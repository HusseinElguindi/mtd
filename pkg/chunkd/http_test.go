package chunkd

import (
	"log"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestHTTPDownload(t *testing.T) {
	f, err := os.OpenFile("./vid.mp4", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		t.FailNow()
	}
	defer f.Close()

	url := "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4"
	// url = "http://ipv4.download.thinkbroadband.com/10MB.zip"
	httpTask := NewHTTP(
		url,
		HTTPWithBufSize(128*1024),
		HTTPWithDestination(f),
		HTTPWithThreads(uint(runtime.NumCPU())),
		HTTPWithWriter(NewWriter(5)),
	)

	done := make(chan struct{})
	defer func() { done <- struct{}{} }()

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		log.Println(httpTask.Status())
		for {
			select {
			case <-ticker.C:
			case <-done:
				log.Println(httpTask.Status())
				return
			}
			log.Println(httpTask.Status())
		}
	}()

	if err := httpTask.Run(); err != nil {
		t.FailNow()
	}
	log.Println("done")
}
