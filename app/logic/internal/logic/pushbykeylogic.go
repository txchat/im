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

type PushByKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushByKeyLogic {
	return &PushByKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushByKey push message from biz level to user client.
func (l *PushByKeyLogic) PushByKey(in *logic.PushByKeyReq) (*logic.Reply, error) {
	reply, err := l.pushByKey(l.ctx, in.GetAppId(), in.GetToKey(), in.GetMsg())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *PushByKeyLogic) pushByKey(c context.Context, appId string, keys []string, msg []byte) (reply *cometclient.ListCastReply, err error) {
	var p protocol.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	pushKeys, err := l.svcCtx.KeysWithServers(c, keys)
	if err != nil {
		return
	}
	for server := range pushKeys {
		if reply, err = l.svcCtx.CometRPC.ListCast(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.ListCastReq{Keys: pushKeys[server], Proto: &p}); err != nil {
			return
		}
	}
	return
}
