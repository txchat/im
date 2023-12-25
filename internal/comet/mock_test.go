package comet

import (
	"bufio"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
)

func Test_verifyMockFrame(t *testing.T) {
	reader, err := authMockFrame("user1-10", 0)
	assert.Nil(t, err)

	var p protocol.Proto
	err = p.ReadFrom(bufio.NewReader(reader))
	assert.Nil(t, err)
	key, hb, err := verifyMockFrame(context.Background(), &p)
	assert.Nil(t, err)
	assert.Equal(t, "user1", key)
	assert.Equal(t, 10*time.Second, hb)
}
