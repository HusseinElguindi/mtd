package chunkd

import "sync"

type Task interface {
	Run(Writer) <-chan error
	Progress() *sync.Cond
	// Progress() <-chan Progress
	Status() Status
}
