package logic

import (
	"context"

	"github.com/txchat/im/internal/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListCastLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCastLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCastLogic {
	return &ListCastLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCastLogic) ListCast(in *comet.ListCastReq) (*comet.ListCastReply, error) {
	if len(in.GetKeys()) == 0 || in.GetProto() == nil {
		return nil, errors.ErrPushMsgArg
	}
	for _, key := range in.GetKeys() {
		if channel := l.svcCtx.Bucket(key).Channel(key); channel != nil {
			err := channel.Push(in.GetProto())
			if err != nil {
				return nil, err
			}
		}
	}
	return &comet.ListCastReply{}, nil
}
