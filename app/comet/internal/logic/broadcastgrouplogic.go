package logic

import (
	"context"

	"github.com/txchat/im/internel/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BroadcastGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBroadcastGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BroadcastGroupLogic {
	return &BroadcastGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BroadcastGroupLogic) BroadcastGroup(in *comet.BroadcastGroupReq) (*comet.BroadcastGroupReply, error) {
	if in.GetProto() == nil || in.GetGroupID() == "" {
		return nil, errors.ErrBroadCastArg
	}
	for _, bucket := range l.svcCtx.Buckets() {
		bucket.BroadcastGroup(in.GetGroupID(), in.GetProto())
	}
	return &comet.BroadcastGroupReply{}, nil
}
