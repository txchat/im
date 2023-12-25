package comet

import (
	"bufio"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
	dtask "github.com/txchat/task"
)

type mockUpgrader struct {
	rr *bufio.Reader
	wr *bufio.Writer
}

func (c *mockUpgrader) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (ProtoReaderWriterCloser, error) {
	if _, ok := conn.(*MockConn); ok {
		c.rr = rr
		c.wr = wr
		return c, nil
	}
	return nil, errors.New("mock conn crash")
}

func (c *mockUpgrader) UpgradeFailed(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (ProtoReaderWriterCloser, error) {
	return nil, errors.New("mock conn crash")
}

func (c *mockUpgrader) ReadProto(proto *protocol.Proto) error {
	return proto.ReadFrom(c.rr)
}

func (c *mockUpgrader) WriteHeart(proto *protocol.Proto, online int32) error {
	return proto.WriteTo(c.wr)
}

func (c *mockUpgrader) WriteProto(proto *protocol.Proto) error {
	return proto.WriteTo(c.wr)
}

func (c *mockUpgrader) Flush() error {
	return c.wr.Flush()
}

func (c *mockUpgrader) Close() error {
	return nil
}

var (
	testRound = NewRound(RoundOptions{
		Reader:       32,
		ReadBuf:      1024,
		ReadBufSize:  8192,
		Writer:       32,
		WriteBuf:     1024,
		WriteBufSize: 8192,
		Timer:        32,
		TimerSize:    2048,
		Task:         32,
		TaskSize:     2048,
	})
	testBuckets  = make([]*Bucket, 32)
	testTaskPool = dtask.NewTask()
	testUpgrader = &mockUpgrader{}
)

func TestMain(m *testing.M) {
	for i := 0; i < len(testBuckets); i++ {
		testBuckets[i] = NewBucket(&BucketConfig{
			Size:          32,
			Channel:       1024,
			Groups:        1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		})
	}
	os.Exit(m.Run())
}

func Test_UpgradeConn(t *testing.T) {
	rb, err := authMockFrame("user1-10", 0)
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb, &MockWriterPrinter{}), rp, wp, tp)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func Test_HandshakeTimeout(t *testing.T) {
	rb, err := authMockFrame("user1-10", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame), WithListenerConfig(&ListenerConfig{
		ProtoClientCacheSize: 10,
		ProtoServerCacheSize: 5,
		KeepAlive:            false,
		TCPReceiveBufferSize: 4096,
		TCPSendBufferSize:    4096,
		HandshakeTimeout:     2 * time.Second,
		RTO:                  3 * time.Second,
		MinHeartbeat:         5 * time.Minute,
		MaxHeartbeat:         10 * time.Minute,
	}))
	conn, err := newConn(l, NewMockConn(rb, &MockWriterPrinter{}), rp, wp, tp)
	assert.EqualError(t, err, "comet connect failed step 2 error: split token failed")
	assert.Nil(t, conn)
}

func Test_UpgradeConnUpgradeFailed(t *testing.T) {
	rb, err := authMockFrame("user1-10", 0)
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.UpgradeFailed, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb, &MockWriterPrinter{}), rp, wp, tp)
	assert.EqualError(t, err, "comet connect failed step 1 error: mock conn crash")
	assert.Nil(t, conn)
}

func Test_UpgradeConnAuthFailed(t *testing.T) {
	rb, err := authMockFrame("", 0)
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb, &MockWriterPrinter{}), rp, wp, tp)
	assert.EqualError(t, err, "comet connect failed step 2 error: split token failed")
	assert.Nil(t, conn)
}
