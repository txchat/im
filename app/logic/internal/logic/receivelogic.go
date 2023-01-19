package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveLogic {
	return &ReceiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveLogic) Receive(in *logic.ReceiveReq) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
