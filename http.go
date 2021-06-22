package mtd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	ErrNoContentLength = errors.New("content length is unset or zero")
)

type metadata struct {
	contentLength    int64
	contentType      string
	acceptByteRanges bool
}

// TODO: implement download retries
// TODO: progress bars
// TODO: background program that connects to terminal command (like Docker cli)

func (t Task) httpHEAD() (metadata, error) {
	// Create request
	req, _ := http.NewRequest(http.MethodHead, t.URL, nil)
	// Set task's request headers
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	// Perform the HEAD request with the task's client
	resp, err := t.Client.Do(req)
	if err != nil {
		return metadata{}, err
	}
	defer resp.Body.Close()

	// Parse headers and create metadata
	return metadata{
		contentLength:    resp.ContentLength,
		contentType:      resp.Header.Get("Content-Type"),
		acceptByteRanges: strings.EqualFold(resp.Header.Get("Accept-Ranges"), "Bytes"),
	}, err
}

func (t Task) Download() error {
	// Get metadata from a HEAD request
	meta, err := t.httpHEAD()
	if err != nil {
		return err
	}
	// Ensure the content has a positive, non-zero length
	if meta.contentLength <= 0 {
		return ErrNoContentLength
	}

	// Determine the chunk size taking into consideration the content length, number of chunks,
	// and whether or not the file server allows for byte ranges
	chunkSize := meta.contentLength
	if meta.acceptByteRanges && t.Chunks > 1 {
		chunkSize = meta.contentLength / int64(t.Chunks) // Floored
	}

	// Prepare waitgroup with the expected number of goroutines
	wg := &sync.WaitGroup{}
	wg.Add(int(t.Chunks))

	// Initialize the error channel
	errc := make(chan error)

	// Distibute the byte ranges between t.Chunks number of goroutines
	var start, end int64
	for i := uint(1); i < t.Chunks; i++ {
		end += chunkSize
		go t.httpWorker(byteRange{start, end - 1}, wg, errc)
		start = end
	}
	// Handle remaining bytes (or all bytes if t.Chunks is 1)
	end += meta.contentLength - start
	go t.httpWorker(byteRange{start, end}, wg, errc)

	// Listen for errors
	go func() {
		for {
			err := <-errc
			// TODO: handle errors
			fmt.Println("error!", err)
		}
	}()

	// Wait for worker goroutines to finish
	wg.Wait()
	return nil
}

func (t Task) httpWorker(bRange byteRange, wg *sync.WaitGroup, errc chan<- error) {
	defer wg.Done()

	// Get body
	rc, err := t.httpGET(bRange)
	if err != nil {
		errc <- err
		return
	}
	defer rc.Close()

	// Write the body
	_, err = t.write(rc, bRange)
	if err != nil {
		errc <- err
		return
	}
}

type byteRange struct{ start, end int64 }

func (b byteRange) Header() string { return fmt.Sprintf("bytes=%d-%d", b.start, b.end) }
func (b byteRange) Valid() bool    { return b.end > 0 && b.end > b.start }

func (t Task) httpGET(bRange byteRange) (io.ReadCloser, error) {
	req, _ := http.NewRequest(http.MethodGet, t.URL, nil)
	// Set task's request headers
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	// Set the byte range header if a range is passed
	if bRange.Valid() {
		req.Header.Set("Range", bRange.Header())
	}

	// Perform the GET request with the task's client
	resp, err := t.Client.Do(req)
	return resp.Body, err
}
