package comet

import (
	"bufio"
	"sync"

	"github.com/Terry-Mao/goim/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	Timer        int
	TimerSize    int
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
	Task         int
	TaskSize     int
}

// Round used for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	readers []BufferPool
	writers []BufferPool
	timers  []time.Timer
	options RoundOptions
}

// NewRound new a round struct.
func NewRound(opts RoundOptions) (r *Round) {
	var i int
	r = &Round{options: opts}
	// reader
	r.readers = make([]BufferPool, r.options.Reader)
	for i = 0; i < r.options.Reader; i++ {
		r.readers[i] = NewReaderPool(r.options.ReadBufSize)
	}
	// writer
	r.writers = make([]BufferPool, r.options.Writer)
	for i = 0; i < r.options.Writer; i++ {
		r.writers[i] = NewWriterPool(r.options.WriteBufSize)
	}
	// timer
	r.timers = make([]time.Timer, r.options.Timer)
	for i = 0; i < r.options.Timer; i++ {
		r.timers[i].Init(r.options.TimerSize)
	}
	return
}

// Timer get a timer.
func (r *Round) Timer(rn int) *time.Timer {
	return &(r.timers[rn%r.options.Timer])
}

// Reader get a reader memory buffer.
func (r *Round) Reader(rn int) BufferPool {
	return r.readers[rn%r.options.Reader]
}

// Writer get a writer memory buffer pool.
func (r *Round) Writer(rn int) BufferPool {
	return r.writers[rn%r.options.Writer]
}

// BufferPool represents a pool of buffers. The *sync.Pool type satisfies this
// interface.  The type of the value stored in a pool is not specified.
type BufferPool interface {
	// Get gets a value from the pool or returns nil if the pool is empty.
	Get() interface{}
	// Put adds a value to the pool.
	Put(interface{})
}

type readerPool struct {
	pool sync.Pool
}

func NewReaderPool(size int) *readerPool {
	return &readerPool{
		pool: sync.Pool{
			// New optionally specifies a function to generate
			// a value when Get would otherwise return nil.
			New: func() interface{} { return bufio.NewReaderSize(nil, size) },
		},
	}
}

// Get gets a value from the pool or returns nil if the pool is empty.
func (p *readerPool) Get() interface{} {
	return p.pool.Get()
}

// Put adds a value to the pool.
func (p *readerPool) Put(b interface{}) {
	p.pool.Put(b)
}

type writerPool struct {
	pool sync.Pool
}

func NewWriterPool(size int) *writerPool {
	return &writerPool{
		pool: sync.Pool{
			// New optionally specifies a function to generate
			// a value when Get would otherwise return nil.
			New: func() interface{} { return bufio.NewWriterSize(nil, size) },
		},
	}
}

// Get gets a value from the pool or returns nil if the pool is empty.
func (p *writerPool) Get() interface{} {
	return p.pool.Get()
}

// Put adds a value to the pool.
func (p *writerPool) Put(b interface{}) {
	p.pool.Put(b)
}
