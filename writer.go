package mtd

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrWritePacket = errors.New("error writing packet")
	// ErrWriterNotListening = errors.New("cannot write to a non-listening writer")
)

type packet struct {
	dst io.WriterAt
	buf []byte
	off int64
}

// type ack struct {
// 	written int
// 	err     error
// }

type Writer struct {
	sync.Mutex
}

// func NewWriter() Writer {
// 	return Writer{
// 		sync.Mutex{},
// 	}
// }

func (w *Writer) Write(p packet) (written int, err error) {
	// Block until free
	w.Lock()
	defer w.Unlock()

	var n int
	for written < len(p.buf[:]) {
		// Write buf at offset
		n, err = p.dst.WriteAt(p.buf[:], p.off+int64(written))
		written += n
		if err != nil {
			return
		}
	}
	return
}

// func (p packet) write(ctx context.Context) ack {
// 	ack := ack{}
// 	for ack.written < len(p.buf[:]) {
// 		// Handle cancel without blocking
// 		select {
// 		case <-ctx.Done():
// 			return ack
// 		default:
// 		}

// 		// Write buf at offset
// 		n, err := p.dst.WriteAt(p.buf[:], p.off+int64(ack.written))
// 		ack.written += n
// 		if err != nil {
// 			ack.err = err
// 			return ack
// 		}
// 	}
// 	return ack
// }
