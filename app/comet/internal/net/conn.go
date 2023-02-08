package net

import (
	"io"
	"net"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/txchat/im/api/protocol"
)

type Conn interface {
	Upgrade(rr *bufio.Reader, wr *bufio.Writer) error
	Close() error
	Flush() error

	WriteProto(proto *protocol.Proto) error
	ReadProto(proto *protocol.Proto) error
	WriteHeart(proto *protocol.Proto, online int32) error
}

type Websocket struct {
	conn   net.Conn
	wsConn *websocket.Conn
}

func NewWebsocket(conn net.Conn) *Websocket {
	return &Websocket{
		conn: conn,
	}
}

func (ws *Websocket) Upgrade(rr *bufio.Reader, wr *bufio.Writer) error {
	req, err := websocket.ReadRequest(rr)
	if err != nil || req.RequestURI != "/sub" {
		return err
	}
	if ws.wsConn, err = websocket.Upgrade(ws.conn, rr, wr, req); err != nil {
		return err
	}
	return nil
}

func (ws *Websocket) WriteProto(p *protocol.Proto) error {
	return p.WriteWebsocket(ws.wsConn)
}

func (ws *Websocket) WriteHeart(p *protocol.Proto, online int32) error {
	return p.WriteWebsocketHeart(ws.wsConn, online)
}

func (ws *Websocket) ReadProto(p *protocol.Proto) error {
	return p.ReadWebsocket(ws.wsConn)
}

func (ws *Websocket) Flush() error {
	return ws.wsConn.Flush()
}

func (ws *Websocket) Close() error {
	return ws.wsConn.Close()
}

type TCP struct {
	rwc io.ReadWriteCloser
	rr  *bufio.Reader
	wr  *bufio.Writer
}

func NewTCP(rwc io.ReadWriteCloser) *TCP {
	return &TCP{
		rwc: rwc,
	}
}

func (tcp *TCP) Upgrade(rr *bufio.Reader, wr *bufio.Writer) error {
	tcp.rr = rr
	tcp.wr = wr
	return nil
}

func (tcp *TCP) WriteProto(p *protocol.Proto) error {
	return p.WriteTCP(tcp.wr)
}

func (tcp *TCP) WriteHeart(p *protocol.Proto, online int32) error {
	return p.WriteTCPHeart(tcp.wr, online)
}

func (tcp *TCP) ReadProto(p *protocol.Proto) error {
	return p.ReadTCP(tcp.rr)
}

func (tcp *TCP) Flush() error {
	return tcp.wr.Flush()
}

func (tcp *TCP) Close() error {
	return tcp.rwc.Close()
}
