package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupsByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupsByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupsByKeysLogic {
	return &LeaveGroupsByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveGroupsByKeysLogic) LeaveGroupsByKeys(in *logic.GroupsKey) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
