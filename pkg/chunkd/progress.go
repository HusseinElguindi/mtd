package chunkd

type Progress struct {
	Off     uint64
	Written uint64
	Done    bool
}
