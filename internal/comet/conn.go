package comet

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Terry-Mao/goim/pkg/bytes"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/pkg/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProtoReader interface {
	ReadProto(proto *protocol.Proto) error
}

type ProtoWriter interface {
	WriteHeart(proto *protocol.Proto, online int32) error
	WriteProto(proto *protocol.Proto) error
	Flush() error
}

type ProtoCloser interface {
	Close() error
}

type ProtoReaderWriter interface {
	ProtoReader
	ProtoWriter
}

type ProtoReaderWriterCloser interface {
	ProtoReader
	ProtoWriter
	ProtoCloser
}

type ConnectHandler func(ctx context.Context, p *protocol.Proto) (key string, hb time.Duration, err error)
type DisconnectHandler func(ctx context.Context, key string) error
type HeartbeatHandler func(ctx context.Context, key string) error

type Conn struct {
	ctx  context.Context
	conn net.Conn
	rwc  ProtoReaderWriterCloser

	b  *Bucket
	ch *Channel

	rp, wp *bytes.Pool
	rb, wb *bytes.Buffer

	resend *Resend

	connect    ConnectHandler
	disconnect DisconnectHandler
	heartbeat  HeartbeatHandler

	hb     time.Duration
	lastHb time.Time
	hbMask time.Duration
	tp     *xtime.Timer
	td     *xtime.TimerData
}

func (c *Conn) resetHeartbeat() {
	c.tp.Set(c.td, c.hb)
}

func (c *Conn) ReadMessage() error {
	p, err := c.ch.CliProto.Set()
	if err != nil {
		return err
	}
	if err = c.rwc.ReadProto(p); err != nil {
		return err
	}
	switch p.Op {
	case int32(protocol.Op_Heartbeat):
		c.resetHeartbeat()
		p.Op = int32(protocol.Op_HeartbeatReply)
		p.Body = nil
		// NOTE: send server heartbeat for a long time
		if now := time.Now(); now.Sub(c.lastHb) > c.hbMask {
			if err := c.heartbeat(c.ctx, c.ch.Key); err == nil {
				c.lastHb = now
			}
		}
	default:
		//TODO read message callback
	}
	// msg sent from client will be dispatched to client itself
	c.ch.CliProto.SetAdv()
	c.ch.Signal()
	return nil
}

func (c *Conn) Close() error {
	c.b.Del(c.ch)
	c.tp.Del(c.td)
	c.rwc.Close()
	c.ch.Close()
	c.rp.Put(c.rb)
	return c.disconnect(c.ctx, c.ch.Key)
}

func newConn(l *Listener, conn net.Conn, rp, wp *bytes.Pool, tp *xtime.Timer) (*Conn, error) {
	var (
		rb = rp.Get()
		wb = wp.Get()
		ch = NewChannel(l.c.ClientCacheSize, l.c.ServerCacheSize)
		rr = &ch.Reader
		wr = &ch.Writer
	)
	// reset reader buffer
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())

	// set handshake timeout
	step := 0
	td := tp.Add(l.c.HandshakeTimeout, func() {
		// NOTE: fix close block for tls
		_ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
		_ = conn.Close()
	})
	ch.IP, ch.Port, _ = net.SplitHostPort(conn.RemoteAddr().String())

	// upgrade connection
	step = 1
	rwc, err := l.upgrade(conn, rr, wr)
	if err != nil {
		conn.Close()
		tp.Del(td)
		rp.Put(rb)
		wp.Put(wb)
		return nil, fmt.Errorf("comet connect failed step %d error: %v", step, err)
	}

	ctx := context.TODO()
	c := &Conn{
		ctx:        ctx,
		conn:       conn,
		rwc:        rwc,
		ch:         ch,
		rp:         rp,
		wp:         wp,
		rb:         rb,
		wb:         wb,
		connect:    l.connect,
		disconnect: l.disconnect,
		heartbeat:  l.heartheat,
		hbMask:     l.RandServerHeartbeat(),
		tp:         tp,
		td:         td,
	}
	// handshake
	// must not setadv, only used in auth
	step = 2
	if ch.Key, c.hb, err = doAuth(ctx, c); err != nil {
		rwc.Close()
		tp.Del(td)
		rp.Put(rb)
		wp.Put(wb)
		return nil, fmt.Errorf("comet connect failed step %d error: %v", step, err)
	}
	td.Key = ch.Key
	c.resetHeartbeat()

	b := l.Bucket(ch.Key)
	b.Put(ch)
	c.b = b

	step = 3
	c.resend = NewResend(ch.Key, l.c.RTO, l.TaskPool, func() {
		ch.Resend()
	})

	go dispatch(c)
	return c, nil
}

// auth for connect handshake with client.
func doAuth(ctx context.Context, conn *Conn) (key string, hb time.Duration, err error) {
	p, err := conn.ch.CliProto.Set()
	if err != nil {
		return
	}
	for {
		if err = conn.rwc.ReadProto(p); err != nil {
			return
		}
		if p.Op == int32(protocol.Op_Auth) {
			var authReplyBody []byte
			key, hb, err = conn.connect(ctx, p)
			if err != nil {
				// set authReplyBody
				authReplyBody, err = auth.ParseGRPCErr(err)
				if err != nil {
					return
				}
			}
			p.Op = int32(protocol.Op_AuthReply)
			p.Ack = p.Seq
			p.Body = authReplyBody
			if err = conn.rwc.WriteProto(p); err != nil {
				return
			}
			err = conn.rwc.Flush()
			break
		}
	}
	return
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func dispatch(conn *Conn) {
	var (
		err    error
		finish bool
		online int32
	)
	for {
		var p = conn.ch.Ready()
		switch p {
		case protocol.ProtoFinish:
			finish = true
			goto failed
		case protocol.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = conn.ch.CliProto.Get(); err != nil {
					break
				}
				if p.Op == int32(protocol.Op_HeartbeatReply) {
					if err = conn.rwc.WriteHeart(p, online); err != nil {
						goto failed
					}
				} else if p.Op == int32(protocol.Op_MessageReply) {
					//del resend but skip conn write
					conn.resend.Del(p.GetAck())
				} else {
					if err = conn.rwc.WriteProto(p); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				conn.ch.CliProto.GetAdv()
			}
		case protocol.ProtoResend:
			for _, rp := range conn.resend.All() {
				if err = conn.rwc.WriteProto(rp); err != nil {
					goto failed
				}
			}
		default:
			switch p.Op {
			case int32(protocol.Op_Message):
				if err = conn.rwc.WriteProto(p); err != nil {
					goto failed
				}
				if err = conn.resend.Add(p); err != nil {
					logx.Error("append resend job error", "err", err)
					goto failed
				}
			default:
				continue
			}
		}
		// only hungry flush response
		if err = conn.rwc.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		logx.Error("dispatch comet conn error", "err", err, "key", conn.ch.Key)
	}
	conn.resend.Stop()
	conn.rwc.Close()
	conn.wp.Put(conn.wb)
	// must ensure all channel message discard, for reader won't block Signal
	for !finish {
		finish = conn.ch.Ready() == protocol.ProtoFinish
	}
}
