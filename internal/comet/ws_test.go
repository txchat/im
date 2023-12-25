package comet

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
)

func TestNewWsStream(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		io.ReadAll(r)
	}()
	req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	req.Header.Add("Sec-Websocket-Version", "13")
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "upgrade")
	req.Header.Add("Sec-Websocket-Key", "2mzaxLsKR++Hp0c5q2ufwg==")
	//w := httptest.NewRecorder()
	respWriter := newHttptestRecorder(r, w)
	wsConn, err := upgrader.Upgrade(respWriter, req, nil)
	assert.Nil(t, err)
	wsStream, err := NewWsStream(wsConn)
	assert.Nil(t, err)
	assert.NotNil(t, wsStream)
}

func TestWsStreamTransBatchProtoWithBody(t *testing.T) {
	c1r, c2w := io.Pipe()
	c2r, c1w := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(2)

	num := 5
	send := make([]*protocol.Proto, 0)
	recv := make([]*protocol.Proto, 0)
	go func() {
		defer wg.Done()

		req, err := http.ReadRequest(bufio.NewReader(c1r))
		assert.Nil(t, err)
		respWriter := newHttptestRecorder(c1r, c1w)
		wsConn, err := upgrader.Upgrade(respWriter, req, nil)
		assert.Nil(t, err)

		wsStream, err := NewWsStream(wsConn)
		assert.Nil(t, err)
		assert.NotNil(t, wsStream)

		for i := 0; i < num; i++ {
			p := &protocol.Proto{
				Ver:  1,
				Op:   int32(protocol.Op_Auth),
				Seq:  0,
				Ack:  0,
				Body: []byte(randStr(500)),
			}
			wb := bufio.NewWriter(wsStream)
			p.WriteTo(wb)
			wb.Flush()
			wsStream.Flush()
			send = append(send, p)
		}
	}()

	go func() {
		defer wg.Done()

		urlStr := "ws://localhost:8300"
		u, err := url.Parse(urlStr)
		assert.Nil(t, err)

		mockConn := NewMockConn(c2r, c2w)
		wsConn, _, err := websocket.NewClient(mockConn, u, nil, 1024, 1024)
		assert.Nil(t, err)

		wsStream, err := NewWsStream(wsConn)
		assert.Nil(t, err)
		assert.NotNil(t, wsStream)

		for i := 0; i < num; i++ {
			p2 := &protocol.Proto{}
			rb := bufio.NewReader(wsStream)
			p2.ReadFrom(rb)
			recv = append(recv, p2)
		}
	}()
	wg.Wait()
	assert.Equal(t, len(send), len(recv))
	for i := 0; i < len(send); i++ {
		assert.EqualValues(t, send[i], recv[i])
	}
	wg.Wait()
}

func TestWsTransBatchProtoWithBody(t *testing.T) {
	c1r, c2w := io.Pipe()
	c2r, c1w := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(2)

	num := 5
	send := make([]*protocol.Proto, 0)
	recv := make([]*protocol.Proto, 0)
	go func() {
		defer wg.Done()

		mockConn := NewMockConn(c1r, c1w)
		wsConn, err := NewWebsocket(mockConn, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
		assert.Nil(t, err)

		for i := 0; i < num; i++ {
			p := &protocol.Proto{
				Ver:  1,
				Op:   int32(protocol.Op_Auth),
				Seq:  0,
				Ack:  0,
				Body: []byte(randStr(500)),
			}
			wsConn.WriteProto(p)
			wsConn.Flush()
			send = append(send, p)
		}
	}()

	go func() {
		defer wg.Done()

		urlStr := "ws://localhost:8300/sub"
		u, err := url.Parse(urlStr)
		assert.Nil(t, err)

		mockConn := NewMockConn(c2r, c2w)
		wsClient, _, err := websocket.NewClient(mockConn, u, nil, 1024, 1024)
		assert.Nil(t, err)

		wsConn, err := FromWebsocketConn(wsClient, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
		assert.Nil(t, err)

		for i := 0; i < num; i++ {
			p2 := &protocol.Proto{}
			wsConn.ReadProto(p2)
			p := &protocol.Proto{
				Ver:  p2.Ver,
				Op:   p2.Op,
				Seq:  p2.Seq,
				Ack:  p2.Ack,
				Body: make([]byte, len(p2.Body)),
			}
			copy(p.Body, p2.Body)
			recv = append(recv, p)
		}
	}()
	wg.Wait()
	assert.Equal(t, len(send), len(recv))
	for i := 0; i < len(send); i++ {
		assert.EqualValues(t, send[i], recv[i])
	}
}

type httptestRecorder struct {
	*httptest.ResponseRecorder
	rb       *bufio.Reader
	wb       *bufio.Writer
	mockConn *MockConn
}

func newHttptestRecorder(r io.Reader, w io.Writer) *httptestRecorder {
	return &httptestRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		rb:               bufio.NewReader(r),
		wb:               bufio.NewWriter(w),
		mockConn:         NewMockConn(r, w),
	}
}

func (c *httptestRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return c.mockConn, bufio.NewReadWriter(c.rb, c.wb), nil
}
