package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushByKeysLogic {
	return &PushByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushByKeysLogic) PushByKeys(in *logic.KeysMsg) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
