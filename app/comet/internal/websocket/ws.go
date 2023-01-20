package websocket

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Terry-Mao/goim/pkg/bytes"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/dtask"
	"github.com/txchat/im/internal/comet"
)

// InitWebsocket listen all tcp.bind and start accept connections.
func InitWebsocket(svcCtx *svc.ServiceContext, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			err = fmt.Errorf("net.ResolveTCPAddr(tcp, %s) err=%e", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			err = fmt.Errorf("net.ListenTCP(tcp, %s) err=%e", bind, err)
			return
		}
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptWebsocket(svcCtx, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptWebsocket(svcCtx *svc.ServiceContext, lis *net.TCPListener) {
	defer func() {
		buf := make([]byte, 1024*3)
		runtime.Stack(buf, false)
		log.Error().Str("panic", string(buf)).Msg("acceptWebsocket done")
		if r := recover(); r != nil {
			buf := make([]byte, 1024*3)
			runtime.Stack(buf, false)
			log.Error().Interface("recover", r).Str("panic", string(buf)).Msg("Recovered in acceptWebsocket")
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
		log.Info().Str("remoteIP", conn.RemoteAddr().String()).Msg("accept ws conn")
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
		go serveWebsocket(svcCtx, conn, r)
		if r++; r == math.MaxInt32 {
			r = 0
		}
	}
}

func serveWebsocket(svcCtx *svc.ServiceContext, conn net.Conn, r int) {
	var (
		// timer
		tr = svcCtx.Round().Timer(r)
		rp = svcCtx.Round().Reader(r)
		wp = svcCtx.Round().Writer(r)
	)
	NewCometServer(svcCtx).ServeWebsocket(conn, rp, wp, tr)
}

type CometServer struct {
	svcCtx *svc.ServiceContext
}

func NewCometServer(svcCtx *svc.ServiceContext) *CometServer {
	return &CometServer{
		svcCtx: svcCtx,
	}
}

// ServeWebsocket serve a websocket connection.
func (s *CometServer) ServeWebsocket(conn net.Conn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	var (
		err    error
		hb     time.Duration
		p      *protocol.Proto
		b      *comet.Bucket
		trd    *xtime.TimerData
		lastHB = time.Now()
		rb     = rp.Get()
		ch     = comet.NewChannel(s.svcCtx.Config.Protocol.CliProto, s.svcCtx.Config.Protocol.SvrProto)
		rr     = &ch.Reader
		wr     = &ch.Writer
		ws     *websocket.Conn // websocket
		req    *websocket.Request
		tsk    *dtask.Task
	)
	// reader
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// handshake
	step := 0
	trd = tr.Add(s.svcCtx.Config.Protocol.HandshakeTimeout, func() {
		// NOTE: fix close block for tls
		_ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
		_ = conn.Close()
		log.Error().Int("step", step).Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("ws handshake timeout")
	})
	// websocket
	ch.IP, ch.Port, _ = net.SplitHostPort(conn.RemoteAddr().String())
	step = 1
	if req, err = websocket.ReadRequest(rr); err != nil || req.RequestURI != "/sub" {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		if err != io.EOF {
			log.Error().Err(err).Msg("http.ReadRequest(rr)")
		}
		return
	}
	// writer
	wb := wp.Get()
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	step = 2
	if ws, err = websocket.Upgrade(conn, rr, wr, req); err != nil {
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
	step = 3
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, hb, err = s.authWebsocket(ctx, ws, p); err == nil {
			log.Info().Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("authoried")
			b = s.svcCtx.Bucket(ch.Key)
			err = b.Put(ch)
		}
	}
	step = 4
	if err != nil {
		ws.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		if err != io.EOF && err != websocket.ErrMessageClose {
			log.Error().Int("step", step).Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("ws handshake failed")
		}
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, hb)
	// hanshake ok start dispatch goroutine
	step = 5
	tsk = dtask.NewTask()
	go s.dispatchWebsocket(ws, wp, wb, ch, tsk)
	serverHeartbeat := s.svcCtx.RandServerHearbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadWebsocket(ws); err != nil {
			break
		}
		if p.Op == int32(protocol.Op_Heartbeat) {
			tr.Set(trd, hb)
			p.Op = int32(protocol.Op_HeartbeatReply)
			p.Body = nil
			// NOTE: send server heartbeat for a long time
			if now := time.Now(); now.Sub(lastHB) > serverHeartbeat {
				if err1 := s.svcCtx.Heartbeat(ctx, ch.Key); err1 == nil {
					lastHB = now
				}
			}
			step++
		} else {
			if err = s.svcCtx.Operate(ctx, p, ch, tsk); err != nil {
				break
			}
		}
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose && !strings.Contains(err.Error(), "closed") {
		log.Error().Err(err).Str("key", ch.Key).Msg("server ws failed")
	}
	b.Del(ch)
	tr.Del(trd)
	ws.Close()
	ch.Close()
	rp.Put(rb)
	if err = s.svcCtx.Disconnect(ctx, ch.Key); err != nil {
		log.Error().Err(err).Str("key", ch.Key).Msg("operator do disconnect")
	}
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *CometServer) dispatchWebsocket(ws *websocket.Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *comet.Channel, tsk *dtask.Task) {
	var (
		err    error
		finish bool
		online int32
		rto    = s.svcCtx.Config.Protocol.Rto
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
					if err = p.WriteWebsocketHeart(ws, online); err != nil {
						goto failed
					}
				} else if p.Op == int32(protocol.Op_ReceiveMsgReply) {
					//skip
				} else if p.Op == int32(protocol.Op_SendMsg) {
					//skip
				} else {
					if err = p.WriteWebsocket(ws); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			switch p.Op {
			case int32(protocol.Op_SendMsgReply):
				if err = p.WriteWebsocket(ws); err != nil {
					goto failed
				}
			case int32(protocol.Op_RePush):
				p.Op = int32(protocol.Op_ReceiveMsg)
				if err = p.WriteWebsocket(ws); err != nil {
					goto failed
				}
			case int32(protocol.Op_ReceiveMsg):
				if err = p.WriteWebsocket(ws); err != nil {
					goto failed
				}
				p.Op = int32(protocol.Op_RePush)
				seq := strconv.FormatInt(int64(p.Seq), 10)
				if j := tsk.Get(seq); j != nil {
					continue
				}
				//push into task pool
				job, inserted := tsk.AddJobRepeat(time.Duration(rto), 0, func() {
					if _, err = ch.Push(p); err != nil {
						log.Error().Err(err).Msg("task job ch.Push error")
						return
					}
				})
				if !inserted {
					log.Error().Err(err).Msg("tsk.AddJobRepeat error")
					goto failed
				}
				tsk.Add(seq, job)
			default:
				continue
			}
		}
		// only hungry flush response
		if err = ws.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose {
		log.Error().Err(err).Str("key", ch.Key).Msg("dispatch ws error")
	}
	tsk.Stop()
	ws.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't blocking Signal
	for !finish {
		finish = (ch.Ready() == protocol.ProtoFinish)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *CometServer) authWebsocket(ctx context.Context, ws *websocket.Conn, p *protocol.Proto) (key string, hb time.Duration, err error) {
reauth:
	for {
		if err = p.ReadWebsocket(ws); err != nil {
			return
		}
		if p.Op == int32(protocol.Op_Auth) {
			break
		} else {
			log.Error().Int32("operation", p.Op).Msg("ws request operation not auth")
		}
	}
	var errMsg string
	if key, hb, errMsg, err = s.svcCtx.Connect(ctx, p); err != nil {
		if errMsg != "" {
			//error result
			log.Debug().Str("errMsg", errMsg).Msg("Connect reject")
			body, e := base64.StdEncoding.DecodeString(errMsg)
			if e != nil {
				log.Error().Err(e).Str("errMsg", errMsg).Msg("base64 Decode errMsg String failed")
				err = e
				return
			}
			p.Op = int32(protocol.Op_AuthReply)
			p.Body = body
			if e := p.WriteWebsocket(ws); e != nil {
				return
			}
			e = ws.Flush()
			if e != nil {
				err = e
			}
			goto reauth
		}
		log.Error().Err(err).Msg("can not call logic.Connect")
		return
	}
	p.Op = int32(protocol.Op_AuthReply)
	p.Body = nil
	if err = p.WriteWebsocket(ws); err != nil {
		return
	}
	err = ws.Flush()
	return
}
