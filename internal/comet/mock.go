package comet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/protocol"
)

func authMockFrame(token string) (io.Reader, error) {
	r, w := io.Pipe()
	b, err := proto.Marshal(&protocol.AuthBody{
		AppId: "mock",
		Token: token,
		Ext:   nil,
	})
	if err != nil {
		return nil, err
	}
	p := &protocol.Proto{
		Ver:  1,
		Op:   int32(protocol.Op_Auth),
		Seq:  0,
		Ack:  0,
		Body: b,
	}
	go func() {
		//write to
		wr := bufio.NewWriter(w)
		p.WriteTo(wr)
		err = wr.Flush()
	}()
	return r, err
}

func verifyMockFrame(ctx context.Context, p *protocol.Proto) (key string, hb time.Duration, err error) {
	var authMsg protocol.AuthBody
	switch p.Op {
	case int32(protocol.Op_Auth):
		if err = proto.Unmarshal(p.GetBody(), &authMsg); err != nil {
			return
		}
	default:
		err = fmt.Errorf("unsupport protocol opration %v", p.Op)
		return
	}
	items := strings.Split(authMsg.Token, "-")
	if len(items) != 2 {
		err = fmt.Errorf("split token failed")
		return
	}
	key = items[0]
	n, err := strconv.Atoi(items[1])
	if err != nil {
		return
	}
	hb = time.Second * time.Duration(n)
	return
}

type MockConn struct {
	rr *bufio.Reader
	wr *bufio.Writer

	rb       io.Reader
	isClosed bool
}

func NewMockConn(rb io.Reader) *MockConn {
	return &MockConn{
		rb: rb,
	}
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *MockConn) Read(b []byte) (n int, err error) {
	if c.isClosed {
		return 0, net.ErrClosed
	}
	return c.rb.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *MockConn) Write(b []byte) (n int, err error) {
	if c.isClosed {
		return 0, net.ErrClosed
	}
	fmt.Printf("write %d bytes\n", len(b))
	return len(b), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *MockConn) Close() error {
	fmt.Print("conn closed")
	c.isClosed = true
	return nil
}

// LocalAddr returns the local network address, if known.
func (c *MockConn) LocalAddr() net.Addr {
	return &net.IPAddr{
		IP:   net.IPv4zero,
		Zone: "",
	}
}

// RemoteAddr returns the remote network address, if known.
func (c *MockConn) RemoteAddr() net.Addr {
	return &net.IPAddr{
		IP:   net.IPv4zero,
		Zone: "",
	}
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *MockConn) SetDeadline(t time.Time) error {
	fmt.Printf("setted deadline as %v\n", t.GoString())
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *MockConn) SetReadDeadline(t time.Time) error {
	fmt.Printf("setted read deadline as %v\n", t.GoString())
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *MockConn) SetWriteDeadline(t time.Time) error {
	fmt.Printf("setted write deadline as %v\n", t.GoString())
	return nil
}

func (c *MockConn) ReadProto(proto *protocol.Proto) error {
	return proto.ReadFrom(c.rr)
}

func (c *MockConn) WriteHeart(proto *protocol.Proto, online int32) error {
	return proto.WriteTo(c.wr)
}

func (c *MockConn) WriteProto(proto *protocol.Proto) error {
	return proto.WriteTo(c.wr)
}

func (c *MockConn) Flush() error {
	return c.wr.Flush()
}
