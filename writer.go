package mtd

import (
	"context"
	"errors"
	"io"
)

var (
	ErrWritePacket        = errors.New("error writing packet")
	ErrWriterNotListening = errors.New("cannot write to a non-listening writer")
)

type packet struct {
	dst io.WriterAt
	buf []byte
	off int64
}
type ack struct {
	written int
	err     error
}

type Writer struct {
	packets chan packet
	acks    chan ack

	ctx context.Context
}

func (w Writer) Write(p packet) (int, error) {
	if w.ctx == nil {
		return 0, ErrWriterNotListening
	}

	w.packets <- p
	ack := <-w.acks
	return ack.written, ack.err
}

func NewWriter() Writer {
	return Writer{
		packets: make(chan packet),
	}
}

// Listen - starts to listen for write packets, consuming them as they come in, until context is cancelled
func (w *Writer) Listen(ctx context.Context) {
	// Only allow one instance of the writer to listen at once
	if w.ctx != nil {
		return
	}
	w.ctx = ctx

	// Listener loop
	for {
		select {
		// Handle cancellations
		case <-w.ctx.Done():
			return
		// Handle packets one at a time, blocking others trying to send
		case p := <-w.packets:
			w.acks <- p.write(w.ctx)
		}
	}
}

func (p packet) write(ctx context.Context) ack {
	ack := ack{}
	for ack.written < len(p.buf) {
		// Handle cancel without blocking
		select {
		case <-ctx.Done():
			return ack
		default:
		}

		// Write buf at offset
		n, err := p.dst.WriteAt(p.buf[:], p.off+int64(ack.written))
		ack.written += n
		if err != nil {
			ack.err = err
			return ack
		}
	}
	return ack
}
