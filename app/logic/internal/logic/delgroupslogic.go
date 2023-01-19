package logic

import (
	"context"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/app/logic/logic"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelGroupsLogic {
	return &DelGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelGroupsLogic) DelGroups(in *logic.DelGroupsReq) (*logic.Reply, error) {
	// todo: add your logic here and delete this line

	return &logic.Reply{}, nil
}
