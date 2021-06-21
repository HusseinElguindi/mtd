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

func NewWriter(ctx context.Context) Writer {
	return Writer{
		packets: make(chan packet),
		acks:    make(chan ack),

		ctx: ctx,
	}
}

func (w Writer) Write(p packet) (int, error) {
	w.packets <- p
	ack := <-w.acks
	return ack.written, ack.err
}

// Listen - starts to listen for write packets, consuming them as they come in, until context is cancelled
func (w *Writer) Listen() {
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
	for ack.written < len(p.buf[:]) {
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
