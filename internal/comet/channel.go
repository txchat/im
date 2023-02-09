package comet

import (
	"sync"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/txchat/im/api/protocol"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	CliProto Ring
	signal   chan *protocol.Proto
	Writer   bufio.Writer
	Reader   bufio.Reader

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
	c.signal = make(chan *protocol.Proto, svr)
	c.nodes = make(map[string]*Node)
	return c
}

// Push server push message.
func (c *Channel) Push(p *protocol.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *protocol.Proto {
	return <-c.signal
}

// Signal send signal to the channel, protocol ready.
func (c *Channel) Signal() {
	c.signal <- protocol.ProtoReady
}

// Resend send signal to the channel, protocol resend.
func (c *Channel) Resend() {
	c.signal <- protocol.ProtoResend
}

// Close notice goroutine finish.
func (c *Channel) Close() {
	c.signal <- protocol.ProtoFinish
}

// Groups get joined groups.
func (c *Channel) Groups() map[string]*Node {
	return c.nodes
}

// DelNode delete the group node by id
func (c *Channel) DelNode(id string) {
	c.mutex.Lock()
	delete(c.nodes, id)
	c.mutex.Unlock()
}

// SetNode set a group node by id
func (c *Channel) SetNode(id string, node *Node) {
	c.mutex.Lock()
	c.nodes[id] = node
	c.mutex.Unlock()
}

// GetNode get a group node by id
func (c *Channel) GetNode(id string) *Node {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.nodes[id]
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
