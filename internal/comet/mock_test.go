package comet

import (
	"context"
	"testing"
	"time"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
)

func Test_verifyMockFrame(t *testing.T) {
	reader, err := authMockFrame("user1-10")
	assert.Nil(t, err)

	var p protocol.Proto
	err = p.ReadTCP(bufio.NewReader(reader))
	assert.Nil(t, err)
	key, hb, err := verifyMockFrame(context.Background(), &p)
	assert.Nil(t, err)
	assert.Equal(t, "user1", key)
	assert.Equal(t, 10*time.Second, hb)
}
