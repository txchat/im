package server

import (
	"time"

	"github.com/txchat/im/pkg/auth"
)

const Name = "mock"

func init() {
	auth.Register(Name, NewAuth)
	auth.RegisterErrorDecoder(Name, nil)
}

func NewAuth(url string, timeout time.Duration) auth.Auth {
	return &mock{}
}
