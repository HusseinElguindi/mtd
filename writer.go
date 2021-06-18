package mtd

import (
	"io"
)

type writeObj struct {
	dst    io.WriterAt
	buf    []byte
	offset uint64
	errCh  chan error
}

type Writer struct {
	writeCh chan writeObj
	closeCh chan int
}

// NewWriter initializes a new writer
func NewWriter() Writer {
	w := Writer{
		writeCh: make(chan writeObj),
		closeCh: make(chan int, 1),
	}
	w.listen()
	return w
}

func (w Writer) listen() {
	go func() {
		for {
			select {
			case writeObj := <-w.writeCh:
				writeObj.write()
			case code := <-w.closeCh:
				if code == 0 {
					w.closeCh <- 0
				}

				break
			}
		}
	}()
}

func (wo writeObj) write() {
	// if _, err := wo.dest.Seek(int64(wo.offset), io.SeekStart); err != nil {
	// 	wo.errCh <- err
	// 	return
	// }

	wrote := 0
	for wrote < len(wo.buf) {
		// n, err := wo.dest.Write(wo.buf)
		n, err := wo.dst.WriteAt(wo.buf, int64(wo.offset+uint64(wrote)))
		if err != nil {
			wo.errCh <- err
			return
		}
		wrote += n
	}

	wo.errCh <- nil
	// close(wo.errCh)
}
