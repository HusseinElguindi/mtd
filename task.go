package mtd

import (
	"context"
	"io"
	"net/http"
)

type DownloadType int

const (
	DownloadTypeHTTP = DownloadType(0)
)

type Task struct {
	DownloadType DownloadType
	URL          string
	// Threads
	// ChunkSize
	Chunks  uint
	BufSize int64

	Client  http.Client
	Headers map[string]string

	Dst    io.WriterAt
	Writer *Writer

	Ctx context.Context
}

func (t Task) write(rc io.ReadCloser, byteRange byteRange) (written int, err error) {
	buf := make([]byte, t.BufSize)

	packet := packet{
		dst: t.Dst,
		off: 0,
	}
	if byteRange.Valid() {
		packet.off = byteRange.start
	}

	for {
		select {
		case <-t.Ctx.Done():
			return
		default:
		}

		// n, err := rc.Read(buf[:])
		var n int
		n, err = io.ReadFull(rc, buf[:])
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}

		packet.buf = buf[:n]
		packet.off += int64(written)

		n, err = t.Writer.Write(packet)
		written += n
		if err != nil {
			return
		}
	}

	return
}
