package mtd

import "sync"

type Status int
type status struct {
	sync.Mutex
	s Status
}

const (
	IDLE        = Status(0)
	IN_PROGRESS = iota
	COMPLETED   = iota
	ERRORED     = iota
)

func (status *status) set(newStatus Status) {
	status.Lock()
	defer status.Unlock()
	status.s = newStatus
}

func (status *status) get() Status {
	status.Lock()
	defer status.Unlock()
	return status.s
}

func (s Status) String() string {
	switch s {
	case IDLE:
		return "IDLE"
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case COMPLETED:
		return "COMPLETED"
	case ERRORED:
		return "ERRORED"
	}
	return "INVALID"
}
