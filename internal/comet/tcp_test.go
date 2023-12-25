package comet

import (
	"bufio"
	"io"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
)

func TestTCPTransProtoNilBody(t *testing.T) {
	r, w := io.Pipe()

	mockConn := NewMockConn(nil, w)
	tcpConn, err := NewTCP(mockConn, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
	assert.Nil(t, err)

	mockConn2 := NewMockConn(r, nil)
	tcpConn2, err := NewTCP(mockConn2, bufio.NewReader(mockConn2), bufio.NewWriter(mockConn2))
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	p := &protocol.Proto{
		Ver:  1,
		Op:   int32(protocol.Op_Auth),
		Seq:  0,
		Ack:  0,
		Body: nil,
	}
	p2 := &protocol.Proto{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		tcpConn.WriteProto(p)
		tcpConn.Flush()
	}()

	go func() {
		defer wg.Done()
		tcpConn2.ReadProto(p2)
	}()
	wg.Wait()
	assert.EqualValues(t, p, p2)
}

func TestTCPTransProtoWithBody(t *testing.T) {
	r, w := io.Pipe()

	mockConn := NewMockConn(nil, w)
	tcpConn, err := NewTCP(mockConn, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
	assert.Nil(t, err)

	mockConn2 := NewMockConn(r, nil)
	tcpConn2, err := NewTCP(mockConn2, bufio.NewReader(mockConn2), bufio.NewWriter(mockConn2))
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	p := &protocol.Proto{
		Ver:  1,
		Op:   int32(protocol.Op_Auth),
		Seq:  0,
		Ack:  0,
		Body: []byte(randStr(500)),
	}
	p2 := &protocol.Proto{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		tcpConn.WriteProto(p)
		tcpConn.Flush()
	}()

	go func() {
		defer wg.Done()
		tcpConn2.ReadProto(p2)
	}()
	wg.Wait()
	assert.EqualValues(t, p, p2)
}

func TestTCPTransBatchProtoWithBody(t *testing.T) {
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
		tcpConn, err := NewTCP(mockConn, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
		assert.Nil(t, err)

		for i := 0; i < num; i++ {
			p := &protocol.Proto{
				Ver:  1,
				Op:   int32(protocol.Op_Auth),
				Seq:  0,
				Ack:  0,
				Body: []byte(randStr(500)),
			}
			tcpConn.WriteProto(p)
			tcpConn.Flush()
			send = append(send, p)
		}
	}()

	go func() {
		defer wg.Done()

		mockConn := NewMockConn(c2r, c2w)
		tcpConn, err := NewTCP(mockConn, bufio.NewReader(mockConn), bufio.NewWriter(mockConn))
		assert.Nil(t, err)

		for i := 0; i < num; i++ {
			p2 := &protocol.Proto{}
			tcpConn.ReadProto(p2)
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

func randStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
