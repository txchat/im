package net

import (
	"errors"
	"fmt"
	"net"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/internal/comet"
)

type Scheme interface {
	Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error)
}

type websocketScheme struct{}

func (ws *websocketScheme) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error) {
	rwc := comet.NewWebsocket(conn)
	return rwc, rwc.Upgrade(rr, wr)
}

type tcpScheme struct{}

func (tcp *tcpScheme) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error) {
	rwc := comet.NewTCP(conn)
	return rwc, rwc.Upgrade(rr, wr)
}

var WebsocketServer = &websocketScheme{}
var TCPServer = &tcpScheme{}

// InitServer listen all tcp.bind and start accept connections.
func InitServer(svcCtx *svc.ServiceContext, address []string, thread int, scheme Scheme) (err error) {
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
		l := comet.NewListener(listener, svcCtx.Round(), svcCtx.Buckets(), svcCtx.TaskPool, scheme.Upgrade, comet.WithListenerConfig(&comet.ListenerConfig{
			ClientCacheSize:   svcCtx.Config.Protocol.CliProto,
			ServerCacheSize:   svcCtx.Config.Protocol.SvrProto,
			KeepAlive:         svcCtx.Config.TCP.KeepAlive,
			ReceiveBufferSize: svcCtx.Config.TCP.Rcvbuf,
			SendBufferSize:    svcCtx.Config.TCP.Sndbuf,
			HandshakeTimeout:  svcCtx.Config.Protocol.HandshakeTimeout,
			RTO:               svcCtx.Config.Protocol.Rto,
			MinHeartbeat:      svcCtx.Config.Protocol.MinHeartbeat,
			MaxHeartbeat:      svcCtx.Config.Protocol.MaxHeartbeat,
		}), comet.WithConnectHandle(svcCtx.Connect), comet.WithDisconnectHandle(svcCtx.Disconnect), comet.WithHeartbeatHandle(svcCtx.Heartbeat))
		// split N core accept
		for i := 0; i < thread; i++ {
			go accept(l)
		}
	}
	return
}

func accept(l *comet.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			// if listener close then return
			//logx.Errorf("listener.Accept(\"%s\") error:", lis.Addr().String(), err)
			continue
			//return
		}
		go func(conn *comet.Conn) {
			defer conn.Close()
			for {
				if err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}(conn)
	}
}
