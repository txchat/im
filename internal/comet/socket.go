package comet

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/google/uuid"
	"github.com/txchat/im/api/protocol"
	dtask "github.com/txchat/task"
	"github.com/zhenjl/cityhash"
)

type ListenerConfig struct {
	ClientCacheSize   int
	ServerCacheSize   int
	KeepAlive         bool
	ReceiveBufferSize int
	SendBufferSize    int
	HandshakeTimeout  time.Duration
	RTO               time.Duration
	MinHeartbeat      time.Duration
	MaxHeartbeat      time.Duration
}

var DefaultListenerConfig = &ListenerConfig{
	ClientCacheSize:   10,
	ServerCacheSize:   5,
	KeepAlive:         false,
	ReceiveBufferSize: 4096,
	SendBufferSize:    4096,
	HandshakeTimeout:  5 * time.Second,
	RTO:               3 * time.Second,
	MinHeartbeat:      5 * time.Minute,
	MaxHeartbeat:      10 * time.Minute,
}
var defaultConnectHandle = func(ctx context.Context, p *protocol.Proto) (key string, hb time.Duration, err error) {
	return fmt.Sprintf("mock-%s", uuid.New().String()), 5 * time.Second, nil
}
var defaultDisconnectHandle = func(ctx context.Context, key string) error {
	return nil
}
var defaultHeartbeatHandle = func(ctx context.Context, key string) error {
	return nil
}

type (
	ListenerOptions struct {
		cfg               *ListenerConfig
		connectHandler    ConnectHandler
		disconnectHandler DisconnectHandler
		heartbeatHandler  HeartbeatHandler
	}

	ListenerOption interface {
		apply(*ListenerOptions)
	}
)

type funcListenerOption struct {
	f func(opts *ListenerOptions)
}

func (fn *funcListenerOption) apply(opts *ListenerOptions) {
	fn.f(opts)
}

func newFuncListenerOption(f func(*ListenerOptions)) *funcListenerOption {
	return &funcListenerOption{
		f: f,
	}
}

func WithListenerConfig(cfg *ListenerConfig) ListenerOption {
	if cfg == nil {
		panic("config is nil")
	}
	return newFuncListenerOption(func(o *ListenerOptions) {
		o.cfg = cfg
	})
}

func WithConnectHandle(handle ConnectHandler) ListenerOption {
	if handle == nil {
		panic("handle is nil")
	}
	return newFuncListenerOption(func(o *ListenerOptions) {
		o.connectHandler = handle
	})
}

func WithDisconnectHandle(handle DisconnectHandler) ListenerOption {
	if handle == nil {
		panic("handle is nil")
	}
	return newFuncListenerOption(func(o *ListenerOptions) {
		o.disconnectHandler = handle
	})
}

func WithHeartbeatHandle(handle HeartbeatHandler) ListenerOption {
	if handle == nil {
		panic("handle is nil")
	}
	return newFuncListenerOption(func(o *ListenerOptions) {
		o.heartbeatHandler = handle
	})
}

type UpgradeHandler func(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (ProtoReaderWriterCloser, error)

type Listener struct {
	lis *net.TCPListener
	c   *ListenerConfig

	round    *Round
	buckets  []*Bucket // subkey bucket
	TaskPool *dtask.Task

	upgrade    UpgradeHandler
	connect    ConnectHandler
	disconnect DisconnectHandler
	heartheat  HeartbeatHandler

	r int
}

func NewListener(l *net.TCPListener, round *Round, buckets []*Bucket, tskPool *dtask.Task, upgrader UpgradeHandler, opt ...ListenerOption) *Listener {
	opts := ListenerOptions{}
	for _, o := range opt {
		o.apply(&opts)
	}
	if round == nil {
		panic("round is nil")
	}
	if len(buckets) == 0 {
		panic("buckets is nil")
	}
	if tskPool == nil {
		panic("task pool is nil")
	}
	if upgrader == nil {
		panic("upgrader is nil")
	}
	if opts.cfg == nil {
		opts.cfg = DefaultListenerConfig
	}
	if opts.connectHandler == nil {
		opts.connectHandler = defaultConnectHandle
	}
	if opts.disconnectHandler == nil {
		opts.disconnectHandler = defaultDisconnectHandle
	}
	if opts.heartbeatHandler == nil {
		opts.heartbeatHandler = defaultHeartbeatHandle
	}
	return &Listener{
		lis:        l,
		c:          opts.cfg,
		round:      round,
		buckets:    buckets,
		TaskPool:   tskPool,
		upgrade:    upgrader,
		connect:    opts.connectHandler,
		disconnect: opts.disconnectHandler,
		heartheat:  opts.heartbeatHandler,
	}
}

func (l *Listener) Accept() (*Conn, error) {
	if l.r++; l.r == math.MaxInt32 {
		l.r = 0
	}
	conn, err := l.lis.AcceptTCP()
	if err != nil {
		// if listener close then return
		return nil, err
	}
	if err = conn.SetKeepAlive(l.c.KeepAlive); err != nil {
		return nil, err
	}
	if err = conn.SetReadBuffer(l.c.ReceiveBufferSize); err != nil {
		return nil, err
	}
	if err = conn.SetWriteBuffer(l.c.SendBufferSize); err != nil {
		return nil, err
	}

	tp := l.round.Timer(l.r)
	rp := l.round.Reader(l.r)
	wp := l.round.Writer(l.r)
	return newConn(l, conn, rp, wp, tp)
}

// Bucket get the bucket by subkey.
func (l *Listener) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % uint32(len(l.buckets))
	return l.buckets[idx]
}

// RandServerHeartbeat rand server heartbeat.
func (l *Listener) RandServerHeartbeat() time.Duration {
	return l.c.MinHeartbeat + time.Duration(rand.Int63n(int64(l.c.MaxHeartbeat-l.c.MinHeartbeat)))
}
