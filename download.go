package mtd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

var (
	ErrNoContentLength = errors.New("no content length was provided")
)

type Task struct {
	url     string
	chunks  uint32
	bufSize uint64

	dst     io.WriterAt
	writer  Writer
	headers http.Header

	closeCh chan uint8
	errCh   chan error
}

func New(URL string, chunks uint32, bufSize uint64, dst io.WriterAt, w Writer) Task {
	return Task{
		url:     URL,
		chunks:  chunks,
		bufSize: bufSize,

		dst:     dst,
		writer:  w,
		headers: nil,

		closeCh: make(chan uint8, chunks),
		errCh:   make(chan error),
	}
}

func (t Task) DownloadHTTP() error {
	contentLength, acceptRanges, _, err := t.head()
	if err != nil {
		return err
	}
	if contentLength <= 0 {
		return ErrNoContentLength
	}

	chunkSize := int64(math.Floor(float64(contentLength) / float64(t.chunks)))
	if !acceptRanges && t.chunks > 1 {
		chunkSize = contentLength
	}

	var start, end uint64
	for i := uint32(1); i < t.chunks; i++ {
		end += uint64(chunkSize)
		go t.downloadChunk(start, end-1)
		start = end
	}
	end += uint64(contentLength) - start
	go t.downloadChunk(start, end)

	// Wait for all goroutines to finish
	for i := uint32(0); i < t.chunks; i++ {
		select {
		case <-t.closeCh: // 1 channel closed
			continue
		case <-t.errCh: // Error occured, close all
			for i := uint32(1); i < t.chunks; i++ {
				t.closeCh <- 1
			}
		}
	}
	return nil
}

// head returns information about a download link
func (t Task) head() (contentLength int64, acceptRanges bool, mime string, err error) {
	req, _ := http.NewRequest(http.MethodHead, t.url, nil)
	setReqHeaders(req, t.headers)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse response headers
	contentLength = resp.ContentLength
	acceptRanges = strings.EqualFold(resp.Header.Get("Accept-Ranges"), "Bytes")
	mime = resp.Header.Get("Content-Type")

	return
}

func (t Task) downloadChunk(start, end uint64) {
	defer func() { t.closeCh <- 1 }()

	req, _ := http.NewRequest(http.MethodGet, t.url, nil)
	setReqHeaders(req, t.headers)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.errCh <- err
		return
	}
	defer resp.Body.Close()

	// TODO: MOVE THIS WHOLE LOOP TO BE A WRITE FUNCTION

	// Buffer
	buf := make([]byte, t.bufSize)
	offset := start
	for {
		select {
		case <-t.closeCh:
			return
		default:
		}

		// Read up to passed read size
		n, err := resp.Body.Read(buf[:])
		if err != nil && err != io.EOF {
			t.errCh <- err
			return
		}

		// No more bytes
		if n == 0 {
			break
		}

		// Write packet
		writeObj := writeObj{
			dst:    t.dst,
			buf:    buf[:n],
			offset: offset,
			errCh:  make(chan error),
		}
		// Send write packet
		t.writer.writeCh <- writeObj
		// Increment offset for next write
		offset += uint64(n)

		// Wait until error chan closes, indicates write was complete
		// err, errored := <-writeObj.errCh
		// if errored && err != nil {
		err = <-writeObj.errCh
		if err != nil {
			t.errCh <- err
		}
		close(writeObj.errCh)
	}

	t.closeCh <- 0
	return
}
