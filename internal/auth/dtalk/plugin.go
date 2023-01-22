package acc

import (
	"time"

	"github.com/txchat/im/internal/auth"
)

const Name = "dtalk"

func init() {
	auth.Register(Name, NewAuth)
}

func NewAuth(url string, timeout time.Duration) auth.Auth {
	return &talkClient{url: url, timeout: timeout}
}
