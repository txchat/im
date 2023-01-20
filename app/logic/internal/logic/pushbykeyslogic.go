package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/svc"
	xkey "github.com/txchat/im/naming/balancer/key"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushByKeysLogic {
	return &PushByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushByKeysLogic) PushByKeys(in *logic.KeysMsg) (*logic.Reply, error) {
	reply, err := l.pushByKeys(l.ctx, in.GetAppId(), in.GetToKeys(), in.GetMsg())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// PushByKeys Push push a biz message to client.
func (l *PushByKeysLogic) pushByKeys(c context.Context, appId string, keys []string, msg []byte) (reply *cometclient.PushMsgReply, err error) {
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
		if reply, err = l.svcCtx.CometRPC.PushMsg(context.WithValue(c, xkey.DefaultKey, server), &cometclient.PushMsgReq{Keys: pushKeys[server], Proto: &p}); err != nil {
			return
		}
	}
	return
}
