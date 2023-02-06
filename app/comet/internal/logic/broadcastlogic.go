package logic

import (
	"context"

	"github.com/txchat/im/internal/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BroadcastLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBroadcastLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BroadcastLogic {
	return &BroadcastLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BroadcastLogic) Broadcast(in *comet.BroadcastReq) (*comet.BroadcastReply, error) {
	if in.GetProto() == nil {
		return nil, errors.ErrBroadCastArg
	}
	// TODO use broadcast queue
	go func() {
		for _, bucket := range l.svcCtx.Buckets() {
			bucket.Broadcast(in.GetProto(), in.GetProtoOp())
		}
	}()
	return &comet.BroadcastReply{}, nil
}
