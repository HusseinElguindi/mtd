package mtd

import "io"

type packet struct {
	dst    io.WriterAt
	buf    []byte
	offset int64
}

type Writer struct {
}
