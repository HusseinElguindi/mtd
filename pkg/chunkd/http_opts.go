package chunkd

import (
	"context"
	"io"
	"net/http"
)

func HTTPWithThreads(threads uint) HTTPOption {
	return func(h *HTTP) {
		if threads == 0 {
			threads = 1
		}
		h.threads = threads
	}
}

func HTTPWithBufSize(bufSize uint64) HTTPOption {
	return func(h *HTTP) {
		if bufSize == 0 {
			bufSize = DefaultHTTPBufSize
		}
		h.bufSize = bufSize
	}
}

func HTTPWithWriter(w Writer) HTTPOption {
	return func(h *HTTP) { h.writer = w }
}

func HTTPWithDestination(dst io.WriterAt) HTTPOption {
	return func(h *HTTP) { h.dst = dst }
}

func HTTPWithClient(client http.Client) HTTPOption {
	return func(h *HTTP) { h.client = client }
}
func HTTPWithHeaders(headers http.Header) HTTPOption {
	return func(h *HTTP) { h.Headers = headers }
}

func HTTPWithContext(ctx context.Context) HTTPOption {
	return func(h *HTTP) { h.ctx = ctx }
}
