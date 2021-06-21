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
	meta, err := t.httpHEAD()
	if err != nil {
		return err
	}

	if meta.contentLength <= 0 {
		return ErrNoContentLength
	}

	chunkSize := meta.contentLength
	if meta.acceptByteRanges && t.Chunks > 1 {
		chunkSize = meta.contentLength / int64(t.Chunks)
	}

	var start, end int64
	wg := &sync.WaitGroup{}
	errc := make(chan error)
	for i := uint(1); i < t.Chunks; i++ {
		end += chunkSize
		wg.Add(1)
		go t.httpWorker(byteRange{start, end - 1}, wg, errc)
		start = end
	}
	end += meta.contentLength - start
	wg.Add(1)
	go t.httpWorker(byteRange{start, end}, wg, errc)

	go func() {
		for {
			err := <-errc
			fmt.Println("error!", err)
		}
	}()

	wg.Wait()
	fmt.Println("done")
	return nil
}

func (t Task) httpWorker(bRange byteRange, wg *sync.WaitGroup, errc chan<- error) {
	defer wg.Done()

	data, err := t.httpGET(bRange)
	if err != nil {
		errc <- err
		return
	}
	defer data.Close()

	_, err = t.write(data, bRange)
	if err != nil {
		errc <- err
		return
	}
}

type byteRange struct{ start, end int64 }

func (b byteRange) Header() string { return fmt.Sprintf("bytes=%d-%d", b.start, b.end) }
func (b byteRange) Valid() bool    { return b.end != 0 && b.end > b.start }

func (t Task) httpGET(byteRange byteRange) (io.ReadCloser, error) {
	req, _ := http.NewRequest(http.MethodGet, t.URL, nil)
	// Set task's request headers
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	// Set the byte range header if a range is passed
	if byteRange.Valid() {
		req.Header.Set("Range", byteRange.Header())
	}

	// Perform the GET request with the task's client
	resp, err := t.Client.Do(req)
	return resp.Body, err
}
