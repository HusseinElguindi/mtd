package chunkd

type Task interface {
	Run(Writer) <-chan error
	Progress() <-chan Progress
	Status() Status
}
