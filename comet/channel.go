package comet

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/comet/grpc"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	CliProto Ring
	signal   chan *grpc.Proto
	Writer   bufio.Writer
	Reader   bufio.Reader

	Seq  int32
	Key  string
	IP   string
	Port string

	nodes map[string]*Node
	mutex sync.RWMutex
}

// NewChannel new a channel.
func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *grpc.Proto, svr)
	c.nodes = make(map[string]*Node)
	return c
}

// Push server push message.
func (c *Channel) seqInc() int32 {
	return atomic.AddInt32(&c.Seq, 1)
}

// Push server push message.
func (c *Channel) push(p *grpc.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Push server push message.
func (c *Channel) Push(p *grpc.Proto) (seq int32, err error) {
	p, ok := proto.Clone(p).(*grpc.Proto)
	if ok {
		if p.Op == int32(grpc.Op_ReceiveMsg) {
			p.Seq = c.seqInc()
		}
		seq = p.Seq
		return seq, c.push(p)
	}
	return 0, errors.New("protocol type gRPC proto failed")
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *grpc.Proto {
	return <-c.signal
}

// Signal send signal to the channel, protocol ready.
func (c *Channel) Signal() {
	c.signal <- grpc.ProtoReady
}

// Close close the channel.
func (c *Channel) Close() {
	c.signal <- grpc.ProtoFinish
}

// Close close the channel.
func (c Channel) Groups() map[string]*Node {
	return c.nodes
}

//
func (c *Channel) DelNode(id string) {
	c.mutex.Lock()
	delete(c.nodes, id)
	c.mutex.Unlock()
}

func (c *Channel) SetNode(id string, node *Node) {
	c.mutex.Lock()
	c.nodes[id] = node
	c.mutex.Unlock()
}

func (c *Channel) GetNode(id string) *Node {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.nodes[id]
}

// base info just for debug
func (c *Channel) GetSeq() int32 {
	return c.Seq
}
func (c *Channel) GetKey() string {
	return c.Key
}
func (c *Channel) GetIP() string {
	return c.IP
}
func (c *Channel) GetPort() string {
	return c.Port
}
func (c *Channel) GetGroups() []string {
	groups := make([]string, len(c.nodes))
	i := 0
	for k := range c.nodes {
		groups[i] = k
		i++
	}
	return groups
}
