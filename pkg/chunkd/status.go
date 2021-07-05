package chunkd

import "sync/atomic"

type Status uint32

const (
	StatusIdle       = Status(0)
	StatusInProgress = iota
	StatusCompleted  = iota
	StatusErrored    = iota
	end              = iota
)

func (s Status) String() string {
	if !s.IsValid() {
		return ""
	}
	return statusLabels[int(s)]
}

func (s Status) IsValid() bool {
	return s < end
}

var statusLabels = [end]string{
	"IDLE",
	"IN PROGRESS",
	"COMPLETED",
	"ERRORED",
}

func (s *Status) StoreAtomic(newStatus Status) {
	if !newStatus.IsValid() {
		return
	}

	atomic.StoreUint32((*uint32)(s), uint32(newStatus))
}

func (s *Status) LoadAtomic() Status {
	return Status(atomic.LoadUint32((*uint32)(s)))
}
