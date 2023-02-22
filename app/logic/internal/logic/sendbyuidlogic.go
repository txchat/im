package logic

import (
	"context"

	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendByUIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendByUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendByUIDLogic {
	return &SendByUIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendByUID push message from biz level to user all clients.
func (l *SendByUIDLogic) SendByUID(in *logic.SendByUIDReq) (*logic.Reply, error) {
	err := l.svcCtx.PublishReceiveMessage(l.ctx, in.GetAppId(), in.GetUid(), "", in.GetOp(), in.GetBody())
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true}, nil
}
