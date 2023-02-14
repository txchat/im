package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/model"
	"github.com/txchat/im/app/logic/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushByUIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushByUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushByUIDLogic {
	return &PushByUIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushByUID push message from biz level to user all clients.
func (l *PushByUIDLogic) PushByUID(in *logic.PushByUIDReq) (*logic.Reply, error) {
	reply, err := l.pushByUID(l.ctx, in.GetAppId(), in.GetToId(), in.GetMsg())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *PushByUIDLogic) pushByUID(c context.Context, appId string, toIds []string, msg []byte) (reply *cometclient.ListCastReply, err error) {
	var p protocol.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	keys, err := l.svcCtx.KeysWithServersByUID(c, appId, toIds)
	if err != nil {
		return
	}

	for server, sKeys := range keys {
		if reply, err = l.svcCtx.CometRPC.ListCast(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.ListCastReq{Keys: sKeys, Proto: &p}); err != nil {
			return
		}
	}
	return
}
