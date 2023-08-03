package comet

import (
	"net"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/Terry-Mao/goim/pkg/websocket"
	"github.com/txchat/im/api/protocol"
)

type Websocket struct {
	conn   net.Conn
	wsConn *websocket.Conn
}

func NewWebsocket(conn net.Conn) *Websocket {
	return &Websocket{
		conn: conn,
	}
}

func (ws *Websocket) SchemeName() string {
	return "websocket"
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
