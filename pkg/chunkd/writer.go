package chunkd

import (
	"io"
)

// Writer represents a concurrent packet writer, and should be created using its NewWriter constructor
type Writer chan struct{}

// NewWriter constructs a new writer queue, with the queue size being the passed number of concurrent writes
func NewWriter(concurrentWrites uint32) Writer {
	// Deadlock occurs if the write channel is unbuffered
	if concurrentWrites < 1 {
		concurrentWrites = 1
	}

	return make(chan struct{}, concurrentWrites)
}

type packet struct {
	dst io.WriterAt
	buf []byte
	off int64
}

// Write writes the passed packet at an offset, to its destination
func (w Writer) Write(p packet) (int, error) {
	// Add write operation to queue
	w <- struct{}{}

	// Defer the removal from queue after write (in case of any panics)
	defer func() { <-w }()

	// Write buf at offset
	return p.dst.WriteAt(p.buf[:], p.off)
}
