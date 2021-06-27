package mtd

import (
	"context"
	"io"
	"net/http"
)

type Task struct {
	URL string
	// Threads uint
	// ChunkSize int64
	Chunks  uint32
	BufSize int64

	Client  http.Client
	Headers map[string]string

	Dst    io.WriterAt
	Writer *Writer

	Ctx context.Context

	status *status
}

func NewTask(URL string, chunks uint32, bufSize int64, dst io.WriterAt, writer *Writer, ctx context.Context) Task {
	return Task{
		URL:     URL,
		Chunks:  chunks,
		BufSize: bufSize,

		Dst:    dst,
		Writer: writer,

		Ctx: ctx,

		status: &status{},
	}
}

func (t Task) Status() Status {
	return t.status.get()
}

func (t Task) write(rc io.ReadCloser, bRange byteRange) (written int, err error) {
	// Determine the smallest useable byte buffer size and allocate the space for it
	bufSize := minInt64(t.BufSize, bRange.end-bRange.start+1)
	buf := make([]byte, bufSize)

	// Prepare backet for reusing
	packet := packet{
		buf: buf,
		dst: t.Dst,
	}

	// Set offset for chunked downloads
	if bRange.Valid() {
		packet.off = bRange.start
	}

	for {
		// Handle cancel without blocking
		select {
		case <-t.Ctx.Done():
			return
		default:
		}

		// Read up to bufSize bytes
		var n int
		n, err = io.ReadFull(rc, buf[:])
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
		n, err = t.Writer.Write(packet)

		// Update offset for next write
		packet.off += int64(n)
		// Update the total number of bytes written
		written += n

		// Handle any write errors
		if err != nil {
			return
		}
	}

	return written, nil
}
