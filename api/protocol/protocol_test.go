package protocol

import (
	"bufio"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProto_ReadWrite(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		p := &Proto{
			Ver:  1,
			Op:   2,
			Seq:  3,
			Ack:  4,
			Body: []byte("hello"),
		}
		wb := bufio.NewWriter(w)
		err := p.WriteTo(wb)
		assert.Nil(t, err)
		err = wb.Flush()
		assert.Nil(t, err)
	}()
	rb := bufio.NewReader(r)
	p2 := &Proto{}
	err := p2.ReadFrom(rb)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), p2.Ver)
	assert.Equal(t, int32(2), p2.Op)
	assert.Equal(t, int32(3), p2.Seq)
	assert.Equal(t, int32(4), p2.Ack)
	assert.Equal(t, "hello", string(p2.GetBody()))
}
