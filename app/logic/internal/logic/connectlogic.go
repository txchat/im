package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectLogic {
	return &ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConnectLogic) Connect(in *logic.ConnectReq) (*logic.ConnectReply, error) {
	// todo: add your logic here and delete this line

	return &logic.ConnectReply{}, nil
}
