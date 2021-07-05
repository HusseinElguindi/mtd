package chunkd

type Progress struct {
	// Off     uint64
	Start   uint64
	End     uint64
	Written uint64
	Done    bool
}
