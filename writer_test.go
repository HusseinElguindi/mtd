package mtd

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

type writeAtBuf bytes.Buffer

func (w writeAtBuf) WriteAt(p []byte, off int64) (n int, err error) {
	b := bytes.Buffer(w)
	return b.Write(p[:])
}

func TestWriterWrite(t *testing.T) {
	b := &writeAtBuf{}

	writer := NewWriter()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go writer.Listen(ctx)

	p := packet{
		dst: b,
		buf: []byte("hello"),
		off: 0,
	}
	fmt.Println(writer.Write(p))
}
