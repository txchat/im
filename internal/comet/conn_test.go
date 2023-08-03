package comet

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/stretchr/testify/assert"
	dtask "github.com/txchat/task"
)

type mockUpgrader struct {
}

func (m *mockUpgrader) Upgrade(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (ProtoReaderWriterCloser, error) {
	if mconn, ok := conn.(*MockConn); ok {
		mconn.rr = rr
		mconn.wr = wr
		return mconn, nil
	}
	return nil, errors.New("mock conn crash")
}

func (m *mockUpgrader) UpgradeFailed(conn net.Conn, rr *bufio.Reader, wr *bufio.Writer) (ProtoReaderWriterCloser, error) {
	return nil, errors.New("mock conn crash")
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
	rb, err := authMockFrame("user1-10")
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb), rp, wp, tp)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	// go func() {
	// 	err = conn.ReadMessage()
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// }()

	time.Sleep(10 * time.Second)
}

func Test_HandshakeTimeout(t *testing.T) {
	rb, err := authMockFrame("user1-10")
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame), WithListenerConfig(&ListenerConfig{
		ClientCacheSize:   10,
		ServerCacheSize:   5,
		KeepAlive:         false,
		ReceiveBufferSize: 4096,
		SendBufferSize:    4096,
		HandshakeTimeout:  2 * time.Second,
		RTO:               3 * time.Second,
		MinHeartbeat:      5 * time.Minute,
		MaxHeartbeat:      10 * time.Minute,
	}))
	conn, err := newConn(l, NewMockConn(rb), rp, wp, tp)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	err = conn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
}

func Test_UpgradeConnUpgradeFailed(t *testing.T) {
	rb, err := authMockFrame("user1-10")
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.UpgradeFailed, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb), rp, wp, tp)
	assert.EqualError(t, err, "comet connect failed step 1 error: mock conn crash")
	assert.Nil(t, conn)
}

func Test_UpgradeConnAuthFailed(t *testing.T) {
	rb, err := authMockFrame("")
	if err != nil {
		t.Fatal(err)
	}

	tp := testRound.Timer(1)
	rp := testRound.Reader(1)
	wp := testRound.Writer(1)
	l := NewListener(nil, testRound, testBuckets, testTaskPool, testUpgrader.Upgrade, WithConnectHandle(verifyMockFrame))
	conn, err := newConn(l, NewMockConn(rb), rp, wp, tp)
	assert.EqualError(t, err, "comet connect failed step 2 error: split token failed")
	assert.Nil(t, conn)
}