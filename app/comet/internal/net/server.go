package net

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/Terry-Mao/goim/pkg/bytes"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/internal/auth"
	"github.com/txchat/im/internal/comet"
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
	defer func() {
		buf := make([]byte, 1024*3)
		runtime.Stack(buf, false)
		log.Error().Str("panic", string(buf)).Msg("accept done")
		if r := recover(); r != nil {
			buf := make([]byte, 1024*3)
			runtime.Stack(buf, false)
			log.Error().Interface("recover", r).Str("panic", string(buf)).Msg("Recovered in accept")
		}
	}()
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			log.Error().Stack().Err(err).Msg(fmt.Sprintf("listener.Accept(\"%s\")", lis.Addr().String()))
			continue
			//return
		}
		log.Info().Str("remoteIP", conn.RemoteAddr().String()).Msg("accept conn")
		if err = conn.SetKeepAlive(svcCtx.Config.TCP.KeepAlive); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetKeepAlive()")
			return
		}
		if err = conn.SetReadBuffer(svcCtx.Config.TCP.Rcvbuf); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetReadBuffer()")
			return
		}
		if err = conn.SetWriteBuffer(svcCtx.Config.TCP.Sndbuf); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetWriteBuffer()")
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
		log.Error().Int("step", step).Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("handshake timeout")
	})
	ch.IP, ch.Port, _ = net.SplitHostPort(conn.RemoteAddr().String())

	step = 1
	if err = cometConn.Upgrade(rr, wr); err != nil {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		wp.Put(wb)
		if err != io.EOF {
			log.Error().Err(err).Msg("websocket.NewServerConn")
		}
		return
	}

	// must not setadv, only used in auth
	step = 2
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, hb, err = s.auth(ctx, cometConn, p); err == nil {
			log.Info().Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("authoried")
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
			log.Error().Int("step", step).Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("handshake failed")
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
		log.Error().Err(err).Str("key", ch.Key).Msg("server comet failed")
	}
	b.Del(ch)
	tr.Del(trd)
	cometConn.Close()
	ch.Close()
	rp.Put(rb)
	if err = s.svcCtx.Disconnect(ctx, ch.Key); err != nil {
		log.Error().Str("key", ch.Key).Err(err).Msg("operator do disconnect")
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
		log.Error().Int32("operation", p.Op).Msg("comet conn request operation not auth")
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
				} else if p.Op == int32(protocol.Op_ReceiveMsgReply) {
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
			case int32(protocol.Op_ReceiveMsg):
				if err = conn.WriteProto(p); err != nil {
					goto failed
				}
				if err = s.resend.Add(p); err != nil {
					log.Error().Err(err).Msg("tsk.AddJobRepeat error")
					goto failed
				}
			case int32(protocol.Op_Transparent):
				if err = conn.WriteProto(p); err != nil {
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
		log.Error().Err(err).Str("key", ch.Key).Msg("dispatch comet conn error")
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
	switch p.Op {
	case int32(protocol.Op_SendMsg):
		err := s.svcCtx.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
		//标明Ack的消息序列
		p.Ack = p.Seq
		p.Op = int32(protocol.Op_SendMsgReply)
		p.Body = nil
	case int32(protocol.Op_ReceiveMsgReply):
		err := s.svcCtx.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
	default:
		return s.svcCtx.Receive(ctx, ch.Key, p)
	}
	return nil
}
