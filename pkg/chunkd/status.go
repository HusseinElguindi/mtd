package chunkd

import "sync/atomic"

type Status uint32

const (
	StatusInProgress = Status(0)
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
	"IN PROGRESS",
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
