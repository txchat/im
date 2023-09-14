package net

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/internal/comet"
)

type Scheme interface {
	Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error)
}

type websocketScheme struct{}

func (ws *websocketScheme) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error) {
	return comet.NewWebsocket(conn, rr, wr)
}

type tcpScheme struct{}

func (tcp *tcpScheme) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (comet.ProtoReaderWriterCloser, error) {
	return comet.NewTCP(conn, rr, wr)
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
			ProtoClientCacheSize: svcCtx.Config.Protocol.CliProto,
			ProtoServerCacheSize: svcCtx.Config.Protocol.SvrProto,
			KeepAlive:            svcCtx.Config.TCP.KeepAlive,
			TCPReceiveBufferSize: svcCtx.Config.TCP.Rcvbuf,
			TCPSendBufferSize:    svcCtx.Config.TCP.Sndbuf,
			HandshakeTimeout:     svcCtx.Config.Protocol.HandshakeTimeout,
			RTO:                  svcCtx.Config.Protocol.Rto,
			MinHeartbeat:         svcCtx.Config.Protocol.MinHeartbeat,
			MaxHeartbeat:         svcCtx.Config.Protocol.MaxHeartbeat,
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
