package net

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strings"
	"time"

	"github.com/Terry-Mao/goim/pkg/bytes"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/internal/comet"
	"github.com/txchat/im/pkg/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type CometConnCreator func(conn net.Conn) Conn

var WebsocketServer = func(conn net.Conn) Conn {
	return NewWebsocket(conn)
}
var TCPServer = func(conn net.Conn) Conn {
	return NewTCP(conn)
}

// InitServer listen all tcp.bind and start accept connections.
func InitServer(svcCtx *svc.ServiceContext, address []string, thread int, scheme CometConnCreator) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	if scheme == nil {
		return errors.New("nil scheme served")
	}
	for _, bind = range address {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			err = fmt.Errorf("net.ResolveTCPAddr(tcp, %s) err=%e", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			err = fmt.Errorf("net.ListenTCP(tcp, %s) err=%e", bind, err)
			return
		}
		// split N core accept
		for i := 0; i < thread; i++ {
			go accept(svcCtx, listener, scheme)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func accept(svcCtx *svc.ServiceContext, lis *net.TCPListener, scheme CometConnCreator) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			logx.Errorf("listener.Accept(\"%s\") error:", lis.Addr().String(), err)
			continue
			//return
		}
		logx.Info("accept conn", "remoteIP", conn.RemoteAddr().String())
		if err = conn.SetKeepAlive(svcCtx.Config.TCP.KeepAlive); err != nil {
			logx.Error("conn.SetKeepAlive()", "err", err)
			return
		}
		if err = conn.SetReadBuffer(svcCtx.Config.TCP.Rcvbuf); err != nil {
			logx.Error("conn.SetReadBuffer()", "err", err)
			return
		}
		if err = conn.SetWriteBuffer(svcCtx.Config.TCP.Sndbuf); err != nil {
			logx.Error("conn.SetWriteBuffer()", "err", err)
			return
		}
		go serve(svcCtx, conn, r, scheme)
		if r++; r == math.MaxInt32 {
			r = 0
		}
	}
}

func serve(svcCtx *svc.ServiceContext, conn net.Conn, r int, scheme CometConnCreator) {
	var (
		// timer
		tr = svcCtx.Round().Timer(r)
		rp = svcCtx.Round().Reader(r)
		wp = svcCtx.Round().Writer(r)
	)

	NewCometServer(svcCtx).Serve(scheme(conn), conn, rp, wp, tr)
}

type CometServer struct {
	svcCtx *svc.ServiceContext

	resend *comet.Resend
}

func NewCometServer(svcCtx *svc.ServiceContext) *CometServer {
	return &CometServer{
		svcCtx: svcCtx,
	}
}

func (s *CometServer) Serve(cometConn Conn, conn net.Conn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	var (
		err    error
		hb     time.Duration
		p      *protocol.Proto
		b      *comet.Bucket
		trd    *xtime.TimerData
		lastHb = time.Now()
		rb     = rp.Get()
		wb     = wp.Get()
		ch     = comet.NewChannel(s.svcCtx.Config.Protocol.CliProto, s.svcCtx.Config.Protocol.SvrProto)
		rr     = &ch.Reader
		wr     = &ch.Writer
	)
	// reset reader buffer
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handshake
	step := 0
	trd = tr.Add(s.svcCtx.Config.Protocol.HandshakeTimeout, func() {
		// NOTE: fix close block for tls
		_ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
		_ = conn.Close()
		logx.Error("handshake timeout", "step", step, "key", ch.Key, "remoteIP", conn.RemoteAddr().String())
	})
	ch.IP, ch.Port, _ = net.SplitHostPort(conn.RemoteAddr().String())

	step = 1
	if err = cometConn.Upgrade(rr, wr); err != nil {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		wp.Put(wb)
		if err != io.EOF {
			logx.Error("cometConn.Upgrade", "err", err)
		}
		return
	}

	// must not setadv, only used in auth
	step = 2
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, hb, err = s.auth(ctx, cometConn, p); err == nil {
			logx.Error("authentication", "key", ch.Key, "remoteIP", conn.RemoteAddr().String())
			b = s.svcCtx.Bucket(ch.Key)
			err = b.Put(ch)
		}
	}
	if err != nil {
		cometConn.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		if err != io.EOF && err != websocket.ErrMessageClose {
			logx.Error("handshake failed", "step", step, "key", ch.Key, "remoteIP", conn.RemoteAddr().String())
		}
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, hb)

	// handshake ok start dispatch goroutine
	step = 3
	s.resend = comet.NewResend(ch.Key, s.svcCtx.Config.Protocol.Rto, s.svcCtx.TaskPool, func() {
		ch.Resend()
	})

	go s.dispatch(cometConn, wp, wb, ch)

	serverHeartbeat := s.svcCtx.RandServerHeartbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = cometConn.ReadProto(p); err != nil {
			break
		}
		if p.Op == int32(protocol.Op_Heartbeat) {
			tr.Set(trd, hb)
			p.Op = int32(protocol.Op_HeartbeatReply)
			p.Body = nil
			// NOTE: send server heartbeat for a long time
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				if err1 := s.svcCtx.Heartbeat(ctx, ch.Key); err1 == nil {
					lastHb = now
				}
			}
			step++
		} else {
			if err = s.operate(ctx, p, ch); err != nil {
				break
			}
		}
		// msg sent from client will be dispatched to client itself
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose && !strings.Contains(err.Error(), "closed") {
		logx.Error("server comet failed", "err", err, "key", ch.Key)
	}
	b.Del(ch)
	tr.Del(trd)
	cometConn.Close()
	ch.Close()
	rp.Put(rb)
	if err = s.svcCtx.Disconnect(ctx, ch.Key); err != nil {
		logx.Error("operator do disconnect", "err", err, "key", ch.Key)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *CometServer) auth(ctx context.Context, conn Conn, p *protocol.Proto) (key string, hb time.Duration, err error) {
auth:
	var authReplyBody []byte
	success := true
	for {
		if err = conn.ReadProto(p); err != nil {
			return
		}
		if p.Op == int32(protocol.Op_Auth) {
			break
		}
		logx.Error("comet conn request operation not auth", "operation", p.Op)
	}
	if key, hb, err = s.svcCtx.Connect(ctx, p); err != nil {
		success = false
		// set authReplyBody
		authReplyBody, err = auth.ParseGRPCErr(err)
		if err != nil {
			return
		}
	}
	p.Op = int32(protocol.Op_AuthReply)
	p.Ack = p.Seq
	p.Body = authReplyBody
	if err = conn.WriteProto(p); err != nil {
		return
	}
	err = conn.Flush()
	if !success {
		goto auth
	}
	return
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *CometServer) dispatch(conn Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *comet.Channel) {
	var (
		err    error
		finish bool
		online int32
	)
	for {
		var p = ch.Ready()
		switch p {
		case protocol.ProtoFinish:
			finish = true
			goto failed
		case protocol.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					break
				}
				if p.Op == int32(protocol.Op_HeartbeatReply) {
					if err = conn.WriteHeart(p, online); err != nil {
						goto failed
					}
				} else if p.Op == int32(protocol.Op_MessageReply) {
					//del resend but skip conn write
					s.resend.Del(p.GetAck())
				} else {
					if err = conn.WriteProto(p); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		case protocol.ProtoResend:
			for _, rp := range s.resend.All() {
				if err = conn.WriteProto(rp); err != nil {
					goto failed
				}
			}
		default:
			switch p.Op {
			case int32(protocol.Op_Message):
				if err = conn.WriteProto(p); err != nil {
					goto failed
				}
				if err = s.resend.Add(p); err != nil {
					logx.Error("append resend job error", "err", err)
					goto failed
				}
			default:
				continue
			}
		}
		// only hungry flush response
		if err = conn.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		logx.Error("dispatch comet conn error", "err", err, "key", ch.Key)
	}
	s.resend.Stop()
	conn.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't block Signal
	for !finish {
		finish = ch.Ready() == protocol.ProtoFinish
	}
}

// Operate operate.
func (s *CometServer) operate(ctx context.Context, p *protocol.Proto, ch *comet.Channel) error {
	//switch p.Op {
	//case int32(protocol.Op_Message):
	//	err := s.svcCtx.Receive(ctx, ch.Key, p)
	//	if err != nil {
	//		//下层业务调用失败，返回error的话会直接断开连接
	//		return err
	//	}
	//	//标明Ack的消息序列
	//	p.Ack = p.Seq
	//	p.Op = int32(protocol.Op_MessageReply)
	//	p.Body = nil
	//case int32(protocol.Op_MessageReply):
	//	err := s.svcCtx.Receive(ctx, ch.Key, p)
	//	if err != nil {
	//		//下层业务调用失败，返回error的话会直接断开连接
	//		return err
	//	}
	//}
	return nil
}
