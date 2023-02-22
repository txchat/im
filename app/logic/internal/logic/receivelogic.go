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

// Receive got message from user client A transfer to biz service.
func (l *ReceiveLogic) Receive(in *logic.ReceiveReq) (*logic.Reply, error) {
	if err := l.receive(l.ctx, in.GetKey(), in.GetProto()); err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true}, nil
}

func (l *ReceiveLogic) receive(c context.Context, key string, p *protocol.Proto) error {
	appId, uid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return err
	}

	msg, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	return l.svcCtx.PublishReceiveMessage(c, appId, uid, key, protocol.Op(p.GetOp()), msg)
}
