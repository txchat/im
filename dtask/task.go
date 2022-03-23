package dtask

import (
	"sync"
	"time"
)

type Task struct {
	sync.RWMutex
	clock    *Clock
	jobCache map[string]Job
}

func NewTask() *Task {
	t := &Task{
		clock:    NewClock(),
		jobCache: make(map[string]Job),
	}
	return t
}

func (t *Task) Add(key string, job Job) {
	t.Lock()
	t.jobCache[key] = job
	t.Unlock()
}

func (t *Task) Get(key string) Job {
	t.RLock()
	defer t.RUnlock()
	return t.jobCache[key]
}

func (t *Task) Stop() {
	t.clock.Stop()
}

func (t *Task) AddJobRepeat(interval time.Duration, actionMax uint64, jobFunc func()) (jobScheduled Job, inserted bool) {
	return t.clock.AddJobRepeat(interval, actionMax, jobFunc)
}
