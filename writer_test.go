package mtd

import (
	"bytes"
	"testing"
)

type writeAtBuf struct {
	bytes.Buffer
}

func (w *writeAtBuf) WriteAt(p []byte, off int64) (n int, err error) {
	return w.Write(p[:])
}

func TestWriterWrite(t *testing.T) {
	buf := &writeAtBuf{}

	writer := Writer{}
	p := packet{
		dst: buf,
		buf: []byte("Hello there"),
		off: 0,
	}
	_, err := writer.Write(p)
	if err != nil {
		t.Fail()
	}

	p2 := packet{
		dst: buf,
		buf: []byte("! How are you?"),
		off: 0,
	}
	_, err = writer.Write(p2)
	if err != nil {
		t.Fail()
	}

	res := append(p.buf, p2.buf...)
	if len(buf.Bytes()) != len(res) {
		t.FailNow()
	}
	for i, b := range buf.Bytes() {
		if res[i] != b {
			t.FailNow()
		}
	}
}
