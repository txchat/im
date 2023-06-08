package svc

import (
	"github.com/txchat/im/app/logic/logicclient"
	"github.com/txchat/im/examples/server/echo/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	LogicClient logicclient.Logic
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		LogicClient: logicclient.NewLogic(zrpc.MustNewClient(c.LogicRPC, zrpc.WithNonBlock(), zrpc.WithNonBlock())),
	}
}
