package main

import (
	"log"
	"os"

	"github.com/husseinelguindi/mtd"
)

func main() {
	URL := "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4"

	f, err := os.OpenFile("./vid.mp4", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	task := mtd.New(URL, 1, 1024*1024*17, f)
	if err := task.DownloadHTTP(); err != nil {
		log.Fatalln(err)
	}
}
