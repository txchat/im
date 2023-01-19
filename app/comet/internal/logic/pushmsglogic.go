package logic

import (
	"context"

	"github.com/txchat/im/internel/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgLogic {
	return &PushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMsgLogic) PushMsg(in *comet.PushMsgReq) (*comet.PushMsgReply, error) {
	if len(in.GetKeys()) == 0 || in.GetProto() == nil {
		return nil, errors.ErrPushMsgArg
	}
	index := make(map[string]int32)
	for _, key := range in.GetKeys() {
		if channel := l.svcCtx.Bucket(key).Channel(key); channel != nil {
			seq, err := channel.Push(in.GetProto())
			if err != nil {
				return nil, err
			}
			index[key] = seq
		}
	}
	return &comet.PushMsgReply{
		Index: index,
	}, nil
}
