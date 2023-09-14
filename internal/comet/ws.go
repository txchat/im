package comet

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/txchat/im/api/protocol"
	xhttp "github.com/txchat/im/internal/http"
)

var upgrader = websocket.Upgrader{} // use default options

type Websocket struct {
	conn *wsStream
	rb   *bufio.Reader
	wb   *bufio.Writer
}

func NewWebsocket(conn net.Conn, rb *bufio.Reader, wb *bufio.Writer) (ProtoReaderWriterCloser, error) {
	rb.Reset(conn)
	wb.Reset(conn)

	req, err := http.ReadRequest(rb)
	//req, err := websocket.ReadRequest(rr)
	if err != nil || req.RequestURI != "/sub" {
		return nil, err
	}

	w, err := xhttp.ReadRequest(wb, req)
	if err != nil {
		return nil, err
	}

	var wsConn *websocket.Conn
	if wsConn, err = upgrader.Upgrade(w, req, nil); err != nil {
		return nil, err
	}
	wsStream, err := NewWsStream(wsConn)
	if err != nil {
		return nil, err
	}

	// must reset source Reader Writer
	rb.Reset(wsStream)
	wb.Reset(wsStream)
	return &Websocket{
		conn: wsStream,
		rb:   rb,
		wb:   wb,
	}, nil
}

func (ws *Websocket) SchemeName() string {
	return "websocket"
}

func (ws *Websocket) WriteProto(p *protocol.Proto) error {
	return p.WriteTo(ws.wb)
}

func (ws *Websocket) ReadProto(p *protocol.Proto) error {
	return p.ReadFrom(ws.rb)
}

func (ws *Websocket) Flush() error {
	ws.wb.Flush()
	return ws.conn.Flush()
}

func (ws *Websocket) Close() error {
	return ws.conn.Close()
}

type wsStream struct {
	wsConn *websocket.Conn
	writer io.WriteCloser
	reader io.Reader
}

func NewWsStream(conn *websocket.Conn) (*wsStream, error) {
	writer, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return nil, err
	}
	_, reader, err := conn.NextReader()
	if err != nil {
		return nil, err
	}
	return &wsStream{
		wsConn: conn,
		writer: writer,
		reader: reader,
	}, nil
}

func (c *wsStream) Write(p []byte) (n int, err error) {
	n, err = c.writer.Write(p)
	return
}

func (c *wsStream) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	if err == io.EOF {
		var reader io.Reader
		_, reader, err = c.wsConn.NextReader()
		if err != nil {
			return
		}
		c.reader = reader
	}
	return
}

func (c *wsStream) Flush() (err error) {
	c.writer.Close()
	c.writer, err = c.wsConn.NextWriter(websocket.BinaryMessage)
	return
}

func (c *wsStream) Close() (err error) {
	err = c.wsConn.Close()
	return
}
