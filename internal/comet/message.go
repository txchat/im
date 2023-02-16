package comet

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/txchat/im/api/protocol"
	dtask "github.com/txchat/task"
)

var ErrTaskJobCantInit = errors.New("task job can not init")

type Resend struct {
	sync.RWMutex
	connID    string
	isOnClock int32
	// items key=seq value=proto
	items map[int32]*protocol.Proto

	repeatJob func()
	rto       time.Duration
	taskPool  *dtask.Task
}

func NewResend(connId string, rto time.Duration, taskPool *dtask.Task, repeatJob func()) *Resend {
	return &Resend{
		connID:    connId,
		items:     make(map[int32]*protocol.Proto),
		repeatJob: repeatJob,
		rto:       rto,
		taskPool:  taskPool,
	}
}

func (rp *Resend) Add(p *protocol.Proto) error {
	rp.Lock()
	rp.items[p.GetSeq()] = p
	rp.Unlock()

	if atomic.CompareAndSwapInt32(&rp.isOnClock, 0, 1) {
		//如果未开启定时任务则需要开启
		job, inserted := rp.taskPool.AddJobRepeat(rp.rto, 0, rp.repeatJob)
		if !inserted {
			return ErrTaskJobCantInit
		}
		rp.taskPool.Add(rp.connID, job)
	}
	return nil
}

func (rp *Resend) Del(seq int32) {
	rp.Lock()
	delete(rp.items, seq)
	rp.Unlock()

	if len(rp.items) == 0 && atomic.CompareAndSwapInt32(&rp.isOnClock, 1, 0) {
		//停止任务
		if j := rp.taskPool.Get(rp.connID); j != nil {
			j.Cancel()
		}
	}
}

func (rp *Resend) All() []*protocol.Proto {
	rp.RLock()
	defer rp.RUnlock()
	list := make([]*protocol.Proto, 0, len(rp.items))
	for _, v := range rp.items {
		list = append(list, v)
	}
	return list
}

func (rp *Resend) Stop() {
	//停止任务
	if j := rp.taskPool.Get(rp.connID); j != nil {
		j.Cancel()
	}
}
