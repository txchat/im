package logic

import (
	"context"

	"github.com/txchat/im/internal/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCastLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCastLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCastLogic {
	return &GroupCastLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCastLogic) GroupCast(in *comet.GroupCastReq) (*comet.GroupCastReply, error) {
	if in.GetProto() == nil || in.GetGid() == "" {
		return nil, errors.ErrBroadCastArg
	}
	for _, bucket := range l.svcCtx.Buckets() {
		bucket.GroupCast(in.GetGid(), in.GetProto())
	}
	return &comet.GroupCastReply{}, nil
}
