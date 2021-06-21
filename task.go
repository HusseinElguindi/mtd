package mtd

import (
	"context"
	"io"
	"net/http"
)

type DownloadType int

const (
	DownloadTypeHTTP = DownloadType(0)
)

type Task struct {
	DownloadType DownloadType
	URL          string
	Chunks       uint
	BufSize      int64

	Client  http.Client
	Headers map[string]string

	Dst    io.WriterAt
	Writer Writer

	Ctx context.Context
}
