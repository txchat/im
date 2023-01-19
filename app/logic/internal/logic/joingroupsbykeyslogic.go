package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupsByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupsByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupsByKeysLogic {
	return &JoinGroupsByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupsByKeysLogic) JoinGroupsByKeys(in *logic.GroupsKey) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
