package mtd

import (
	"bytes"
	"context"
	"fmt"
	"os"
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

func TestHTTPDownload(t *testing.T) {
	f, err := os.OpenFile("./vid.mp4", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		t.FailNow()
	}
	defer f.Close()

	writer := Writer{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := Task{
		DownloadType: DownloadTypeHTTP,
		URL:          "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4",
		Chunks:       1,
		BufSize:      7 * 1024 * 1024, // 7mb

		Dst:    f,
		Writer: &writer,

		Ctx: ctx,
	}

	err = task.Download()
	fmt.Println(err)
}
