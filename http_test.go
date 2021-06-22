package mtd

import (
	"context"
	"os"
	"runtime"
	"testing"
)

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
		URL:     "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4",
		Chunks:  uint(runtime.NumCPU()),
		BufSize: 7 * 1024 * 1024, // ~7mb

		Dst:    f,
		Writer: &writer,

		Ctx: ctx,
	}

	if err := task.Download(); err != nil {
		t.FailNow()
	}
}
