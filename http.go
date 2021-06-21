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

func (t Task) HTTPInit() error {
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

	wg := &sync.WaitGroup{}

	byteRange := byteRange{}
	errors := make(chan error)
	for i := uint(1); i < t.Chunks; i++ {
		wg.Add(1)
		byteRange.end += chunkSize
		go func() {
			defer wg.Done()

			byteRange := byteRange
			byteRange.end -= 1
			data, err := t.httpGET(&byteRange)
			if err != nil {
				errors <- err
				return
			}
			defer data.Close()

			ack := t.write(data, &byteRange)
			if ack.err != nil {
				errors <- err
				return
			}
		}()
		byteRange.start = byteRange.end
	}
	byteRange.end += meta.contentLength - byteRange.start
	wg.Add(1)
	go func() {
		defer wg.Done()

		byteRange := byteRange
		data, err := t.httpGET(&byteRange)
		if err != nil {
			errors <- err
			return
		}
		defer data.Close()

		ack := t.write(data, &byteRange)
		fmt.Println(ack)
		if ack.err != nil {
			errors <- err
			return
		}
	}()

	go func() {
		for {
			err := <-errors
			fmt.Println(err)
		}
	}()

	wg.Wait()
	fmt.Println("done")
	return nil
}

type byteRange struct {
	start, end int64
}

func (b byteRange) String() string { return fmt.Sprintf("bytes=%d-%d", b.start, b.end) }

func (t Task) httpGET(byteRange *byteRange) (io.ReadCloser, error) {
	fmt.Println(byteRange)
	req, _ := http.NewRequest(http.MethodGet, t.URL, nil)
	// Set task's request headers
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	// Set the byte range header if a range is passed
	if byteRange != nil {
		req.Header.Set("Range", byteRange.String())
	}

	// Perform the GET request with the task's client
	resp, err := t.Client.Do(req)
	return resp.Body, err
}

func (t Task) write(rc io.ReadCloser, byteRange *byteRange) ack {
	buf := make([]byte, t.BufSize)
	var written int64

	packet := packet{
		dst: t.Dst,
		off: 0,
	}
	if byteRange != nil {
		packet.off = byteRange.start
	}

	for {
		select {
		case <-t.Ctx.Done():
			return ack{int(written), nil}
		default:
		}

		// n, err := rc.Read(buf[:])
		n, err := io.ReadFull(rc, buf[:])
		if err != nil && err != io.EOF {
			return ack{int(written), nil}
		}

		if n == 0 {
			break
		}

		fmt.Println(packet.off, n, written)

		packet.buf = buf[:n]
		packet.off += written

		n, err = t.Writer.Write(packet)
		written += int64(n)
		if err != nil {
			return ack{int(written), err}
		}
	}

	return ack{int(written), nil}
}
