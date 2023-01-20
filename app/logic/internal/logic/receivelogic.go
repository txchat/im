package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveLogic {
	return &ReceiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveLogic) Receive(in *logic.ReceiveReq) (*logic.Reply, error) {
	if err := l.receive(l.ctx, in.GetKey(), in.GetProto()); err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true}, nil
}

// Receive receive a message from client.
func (l *ReceiveLogic) receive(c context.Context, key string, p *protocol.Proto) (err error) {
	appId, mid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return err
	}

	msg, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	return l.svcCtx.PublishMsg(c, appId, mid, protocol.Op(p.Op), key, msg)
}
