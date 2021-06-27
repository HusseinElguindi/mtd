package scheduler

import (
	"context"
	"sync"

	mtd "github.com/husseinelguindi/mtd/lib"
)

type Handle struct {
	sync.Mutex
	task *mtd.Task

	cancel context.CancelFunc
}

type Scheduler struct {
	sync.RWMutex
	m map[string]*Handle
}

func NewScheduler() Scheduler {
	return Scheduler{
		m: make(map[string]*Handle),
	}
}

func (s *Scheduler) RunTask(t mtd.Task) (id string) {
	handle := &Handle{
		task: &t,
	}
	t.Ctx, handle.cancel = context.WithCancel(context.Background())
	go t.Download()

	s.RLock()
	for {
		id = GenerateID()
		if _, ok := s.m[id]; !ok {
			break
		}
	}
	s.RUnlock()
	s.Set(id, handle)

	return
}

func (s *Scheduler) CancelTask(id string) {
	handle := s.Get(id)
	if handle == nil {
		return
	}

	handle.cancel()
	s.Del(id)
}

func (s *Scheduler) TaskStatus(id string) (mtd.Status, bool) {
	handle := s.Get(id)
	if handle == nil {
		return -1, false
	}

	handle.Lock()
	defer handle.Unlock()

	return handle.task.Status(), true
}

func (s *Scheduler) Get(id string) *Handle {
	s.RLock()
	defer s.RUnlock()
	return s.m[id]
}

func (s *Scheduler) Set(id string, handle *Handle) {
	s.Lock()
	defer s.Unlock()
	s.m[id] = handle
}

func (s *Scheduler) Del(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, id)
}
