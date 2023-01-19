package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupsByMidsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupsByMidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupsByMidsLogic {
	return &JoinGroupsByMidsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupsByMidsLogic) JoinGroupsByMids(in *logic.GroupsMid) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
