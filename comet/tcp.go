package comet

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/Terry-Mao/goim/pkg/bytes"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/dtask"
)

// InitTCP listen all tcp.bind and start accept connections.
func InitTCP(server *Comet, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			log.Error().Stack().Err(err).Msg(fmt.Sprintf("net.ResolveTCPAddr(tcp, %s)", bind))
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			log.Error().Stack().Err(err).Msg(fmt.Sprintf("net.ListenTCP(tcp, %s)", bind))
			return
		}
		log.Info().Str("bind", bind).Msg("start tcp listen")
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptTCP(server, listener)
		}
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func acceptTCP(server *Comet, lis *net.TCPListener) {
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
		log.Info().Str("remoteIP", conn.RemoteAddr().String()).Msg("accept tcp conn")
		if err = conn.SetKeepAlive(server.c.TCP.KeepAlive); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetKeepAlive()")
			return
		}
		if err = conn.SetReadBuffer(server.c.TCP.Rcvbuf); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetReadBuffer()")
			return
		}
		if err = conn.SetWriteBuffer(server.c.TCP.Sndbuf); err != nil {
			log.Error().Stack().Err(err).Msg("conn.SetWriteBuffer()")
			return
		}
		go serveTCP(server, conn, r)
		if r++; r == maxInt {
			r = 0
		}
	}
}

func serveTCP(s *Comet, conn *net.TCPConn, r int) {
	var (
		// timer
		tr = s.round.Timer(r)
		rp = s.round.Reader(r)
		wp = s.round.Writer(r)
	)
	s.ServeTCP(conn, rp, wp, tr)
}

// ServeTCP serve a tcp connection.
func (s *Comet) ServeTCP(conn *net.TCPConn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	var (
		err    error
		hb     time.Duration
		p      *grpc.Proto
		b      *Bucket
		trd    *xtime.TimerData
		lastHb = time.Now()
		rb     = rp.Get()
		wb     = wp.Get()
		ch     = NewChannel(s.c.Protocol.CliProto, s.c.Protocol.SvrProto)
		rr     = &ch.Reader
		wr     = &ch.Writer
		tsk    *dtask.Task
	)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// handshake
	step := 0
	trd = tr.Add(time.Duration(s.c.Protocol.HandshakeTimeout), func() {
		conn.Close()
		log.Error().Int("step", step).Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("tcp handshake timeout")
	})
	ch.IP, ch.Port, _ = net.SplitHostPort(conn.RemoteAddr().String())
	// must not setadv, only used in auth
	step = 1
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Key, hb, err = s.authTCP(ctx, rr, wr, p); err == nil {
			log.Info().Str("key", ch.Key).Str("remoteIP", conn.RemoteAddr().String()).Msg("authoried")
			b = s.Bucket(ch.Key)
			err = b.Put(ch)
		}
	}
	step = 2
	if err != nil {
		conn.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		log.Error().Str("key", ch.Key).Err(err).Msg("handshake failed")
		return
	}
	trd.Key = ch.Key
	tr.Set(trd, hb)
	step = 3
	// hanshake ok start dispatch goroutine
	tsk = dtask.NewTask()
	go s.dispatchTCP(conn, wr, wp, wb, ch, tsk)
	serverHeartbeat := s.RandServerHearbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadTCP(rr); err != nil {
			log.Info().Err(err).Msg("ReadTCP failed")
			break
		}
		if p.Op == int32(grpc.Op_Heartbeat) {
			tr.Set(trd, hb)
			p.Op = int32(grpc.Op_HeartbeatReply)
			p.Body = nil
			// NOTE: send server heartbeat for a long time
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				if err1 := s.Heartbeat(ctx, ch.Key); err1 == nil {
					lastHb = now
				}
			}
			step++
		} else {
			if err = s.Operate(ctx, p, ch, tsk); err != nil {
				break
			}
		}
		// msg sent from client will be dispatched to client itself
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		log.Error().Str("key", ch.Key).Err(err).Msg("server tcp failed")
	}
	b.Del(ch)
	tr.Del(trd)
	rp.Put(rb)
	conn.Close()
	ch.Close()
	if err = s.Disconnect(ctx, ch.Key); err != nil {
		log.Error().Str("key", ch.Key).Err(err).Msg("operator do disconnect")
	}
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (s *Comet) dispatchTCP(conn *net.TCPConn, wr *bufio.Writer, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel, tsk *dtask.Task) {
	var (
		err    error
		finish bool
		online int32
		rto    = s.c.Protocol.Rto
	)
	for {
		var p = ch.Ready()
		switch p {
		case grpc.ProtoFinish:
			finish = true
			goto failed
		case grpc.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					break
				}
				if p.Op == int32(grpc.Op_HeartbeatReply) {
					if err = p.WriteTCPHeart(wr, online); err != nil {
						goto failed
					}
				} else if p.Op == int32(grpc.Op_ReceiveMsgReply) {
					//skip
				} else if p.Op == int32(grpc.Op_SendMsg) {
					//skip
				} else {
					if err = p.WriteTCP(wr); err != nil {
						goto failed
					}
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			switch p.Op {
			case int32(grpc.Op_SendMsgReply):
				if err = p.WriteTCP(wr); err != nil {
					goto failed
				}
			case int32(grpc.Op_RePush):
				pro := proto.Clone(p)
				if p, ok := pro.(*grpc.Proto); ok {
					p.Op = int32(grpc.Op_ReceiveMsg)
					if err = p.WriteTCP(wr); err != nil {
						goto failed
					}
				} else {
					log.Error().Msg("proto can`t clone Proto")
				}
			case int32(grpc.Op_ReceiveMsg):
				if err = p.WriteTCP(wr); err != nil {
					goto failed
				}
				p.Op = int32(grpc.Op_RePush)
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
		if err = wr.Flush(); err != nil {
			break
		}
	}
failed:
	if err != nil {
		log.Error().Str("key", ch.Key).Err(err).Msg("dispatch tcp failed")
	}
	tsk.Stop()
	conn.Close()
	wp.Put(wb)
	// must ensure all channel message discard, for reader won't blocking Signal
	for !finish {
		finish = (ch.Ready() == grpc.ProtoFinish)
	}
}

// auth for goim handshake with client, use rsa & aes.
func (s *Comet) authTCP(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *grpc.Proto) (key string, hb time.Duration, err error) {
reauth:
	for {
		if err = p.ReadTCP(rr); err != nil {
			return
		}
		if p.Op == int32(grpc.Op_Auth) {
			break
		} else {
			log.Error().Int32("option", p.Op).Msg("tcp request option not auth")
		}
	}
	var errMsg string
	if key, hb, errMsg, err = s.Connect(ctx, p); err != nil {
		if errMsg != "" {
			//error result
			log.Debug().Str("errMsg", errMsg).Msg("Connect reject")
			body, e := base64.StdEncoding.DecodeString(errMsg)
			if e != nil {
				log.Error().Err(e).Str("errMsg", errMsg).Msg("base64 Decode errMsg String failed")
				err = e
				return
			}
			p.Op = int32(grpc.Op_AuthReply)
			p.Body = body
			if e := p.WriteTCP(wr); e != nil {
				return
			}
			e = wr.Flush()
			if e != nil {
				err = e
			}
			goto reauth
		}
		log.Error().Err(err).Msg("can not call logic.Connect")
		return
	}
	p.Op = int32(grpc.Op_AuthReply)
	p.Body = nil
	if err = p.WriteTCP(wr); err != nil {
		return
	}
	err = wr.Flush()
	return
}
