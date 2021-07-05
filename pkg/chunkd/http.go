package chunkd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

const DefaultHTTPBufSize = 64 * 1024

var ErrNoContentLength = errors.New("content length is unset or zero")

type HTTP struct {
	url     string
	threads uint
	bufSize uint64

	client  http.Client
	Headers http.Header

	dst    io.WriterAt
	writer Writer

	progress chan Progress
	status   Status

	ctx context.Context
}

type HTTPOption func(*HTTP)

func NewHTTP(url string, opts ...HTTPOption) HTTP {
	h := HTTP{
		url:     url,
		threads: uint(runtime.NumCPU()),
		bufSize: DefaultHTTPBufSize,
	}

	for _, opt := range opts {
		opt(&h)
	}

	if h.progress == nil {
		h.progress = make(chan Progress, h.threads)
	}
	if h.writer == nil {
		h.writer = NewWriter(1)
	}
	if h.ctx == nil {
		h.ctx = context.Background()
	}

	// TODO
	// if h.dst == nil {
	// implement function to make file from passed filename or filename from url/headers
	// }

	return h
}

func (h *HTTP) Run() error {
	h.status.StoreAtomic(StatusInProgress)
	if err := h.run(); err != nil {
		h.status.StoreAtomic(StatusErrored)
		return err
	}
	h.status.StoreAtomic(StatusCompleted)
	return nil
}

func (h HTTP) run() error {
	metadata, err := h.head()
	if err != nil {
		return err
	}
	// Ensure the content has a positive, non-zero length
	if metadata.contentLength <= 0 {
		return ErrNoContentLength
	}

	// Determine the chunk size taking into consideration the content length, number of chunks,
	// and whether or not the file server allows for byte ranges
	chunkSize := uint64(metadata.contentLength)
	if metadata.acceptByteRanges && h.threads > 1 {
		chunkSize = uint64(metadata.contentLength) / uint64(h.threads) // Floored
	}

	// Initialize the error channel
	errc := make(chan error)

	// Distibute the byte ranges between t.Chunks number of goroutines
	var start, end uint64
	for i := uint(1); i < h.threads; i++ {
		end += uint64(chunkSize)
		go h.worker(start, end-1, errc)
		start = end
	}
	// Handle remaining bytes (or all bytes if t.Chunks is 1)
	end += uint64(metadata.contentLength) - start
	go h.worker(start, end, errc)

	// Listen for errors
	go func() {
		for {
			err := <-errc
			// TODO: handle errors
			fmt.Println("error!", err)
		}
	}()

	chunks := make(map[uint64]Progress, h.threads)
	done := uint(0)
	for p := range h.progress {
		chunks[p.Start] = p
		fmt.Printf("progress: %+v\n", p)

		if !p.Done {
			continue
		}

		// Break if all threads finished
		done++
		if done == h.threads {
			break
		}

		// If a thread finished while others are still working
		// if h.threads > 1 {
		// 	fmt.Println("dymanic segmentation")
		// 	// TODO: dynamic segmentation
		// 	min := p
		// 	for _, c := range chunks {
		// 		if c.Written < min.Written {
		// 			min = c
		// 		}
		// 	}

		// 	fmt.Println("biggest gap:", min.End-(min.Start+min.Written-1))

		// 	// https://www.internetdownloadmanager.com/support/segmentation.html
		// 	// min is the progress of the goroutine that wrote the least data; min.off is its "id"
		// 	// do smth with min (maybe use sync.cond?)
		// 	fmt.Printf("min: %+v, done: %+v\n", min, p)
		// }
	}

	return nil
}

type metadata struct {
	contentLength    int64
	contentType      string
	filename         string // from content-disposition header or url
	acceptByteRanges bool
}

func (h *HTTP) head() (metadata, error) {
	// Create request
	req, _ := http.NewRequest(http.MethodHead, h.url, nil)
	// Set task's request headers
	for k := range h.Headers {
		req.Header.Set(k, h.Headers.Get(k))
	}

	// Perform the HEAD request with the task's client
	resp, err := h.client.Do(req)
	if err != nil {
		return metadata{}, err
	}
	defer resp.Body.Close()

	// Parse headers and create metadata
	return metadata{
		contentLength: resp.ContentLength,
		contentType:   resp.Header.Get("Content-Type"),
		// filename:         resp.Header.Get("Content-Disposition"), // needs more work, see https://stackoverflow.com/a/28845255/9443397
		acceptByteRanges: strings.EqualFold(resp.Header.Get("Accept-Ranges"), "Bytes"),
	}, err
}

func (h HTTP) worker(start, end uint64, errc chan<- error) {
	byteRange := byteRange{start, end}
	body, err := h.get(byteRange)
	if err != nil {
		errc <- err
		return // or retry
	}
	defer body.Close()

	_, err = h.write(body, byteRange)
	if err != nil {
		errc <- err
		return
	}
}

type byteRange struct{ start, end uint64 }

func (b byteRange) Header() string { return fmt.Sprintf("bytes=%d-%d", b.start, b.end) }
func (b byteRange) Valid() bool    { return b.end > 0 && b.end > b.start }

func (h HTTP) get(byteRange byteRange) (io.ReadCloser, error) {
	req, _ := http.NewRequest(http.MethodGet, h.url, nil)
	// Set task's request headers
	for k := range h.Headers {
		req.Header.Set(k, h.Headers.Get(k))
	}

	// Set the byte range header if a range is passed
	if byteRange.Valid() {
		req.Header.Set("Range", byteRange.Header())
	}

	// Perform the GET request with the task's client
	resp, err := h.client.Do(req)
	return resp.Body, err
}

func (h HTTP) write(body io.ReadCloser, byteRange byteRange) (written int, err error) {
	// Determine the smallest useable byte buffer size and allocate the space for it
	bufSize := minUint64(h.bufSize, byteRange.end-byteRange.start+1)
	buf := make([]byte, bufSize)

	// Prepare packet for reusing
	packet := packet{
		dst: h.dst,
	}

	// Prepare progress
	progress := Progress{}

	// Set offset for chunked downloads
	if byteRange.Valid() {
		packet.off = int64(byteRange.start)

		progress.Start = uint64(byteRange.start)
		progress.End = uint64(byteRange.end)
	}

	for {
		// Handle cancel without blocking
		select {
		case <-h.ctx.Done():
			return
		default:
		}

		// Read up to bufSize bytes
		var n int
		n, err = io.ReadFull(body, buf[:])
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return
		}

		// Break if no bytes were read (no bytes left to write)
		if n == 0 {
			break
		}

		// Set the packet's buffer as a slice of the newly read bytes
		packet.buf = buf[:n]

		// Write the packet
		n, err = h.writer.Write(packet)

		// Update offset for next write
		packet.off += int64(n)
		// Update the total number of bytes written
		written += n

		// Send progress
		progress.Written = uint64(written)
		if (progress.Start + progress.Written - 1) == progress.End {
			break
		}
		h.progress <- progress

		// Handle any write errors
		if err != nil {
			return
		}
	}

	progress.Done = true
	h.progress <- progress

	return written, nil
}

// Change implementation
// func (h *HTTP) Progress() <-chan Progress {
// 	return h.progress
// }

func (h *HTTP) Status() Status {
	return h.status.LoadAtomic()
}
