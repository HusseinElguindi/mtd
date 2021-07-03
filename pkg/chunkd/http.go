package chunkd

import "net/http"

type HTTP struct {
	url     string
	threads uint
	bufsize uint

	Headers http.Header

	progress chan Progress
	status   Status
}

func NewHTTP(url string, threads uint, bufsize uint) HTTP {
	if threads == 0 {
		threads = 1
	}
	if bufsize == 0 {
		bufsize = 64 * 1024
	}

	return HTTP{
		url:     url,
		threads: threads,
		bufsize: bufsize,

		progress: make(chan Progress),
	}
}

func (h *HTTP) Run(w Writer) <-chan error {
	progress := make(chan Progress)
	chunks := make(map[uint64]Progress, h.threads)
	for p := range progress {
		chunks[p.Off] = p
		if p.Done {
			// do dynamic segmentation
		}
	}

	return nil
}

// Change implementation
func (h *HTTP) Progress() <-chan Progress {
	return h.progress
}

func (h *HTTP) Status() Status {
	return h.status.LoadAtomic()
}
