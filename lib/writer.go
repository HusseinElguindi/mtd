package mtd

import (
	"io"
	"sync"
)

type Writer struct {
	sync.Mutex
}

type packet struct {
	dst io.WriterAt
	buf []byte
	off int64
}

func (w *Writer) Write(p packet) (int, error) {
	// Block until free
	w.Lock()
	defer w.Unlock()

	// Write buf at offset
	return p.dst.WriteAt(p.buf[:], p.off)
}
