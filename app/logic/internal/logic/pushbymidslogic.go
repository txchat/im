package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushByMidsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushByMidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushByMidsLogic {
	return &PushByMidsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushByMidsLogic) PushByMids(in *logic.MidsMsg) (*logic.Reply, error) {
	reply, err := l.pushByMids(l.ctx, in.GetAppId(), in.GetToIds(), in.GetMsg())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// PushByMids Push push a biz message to client.
func (l *PushByMidsLogic) pushByMids(c context.Context, appId string, toIds []string, msg []byte) (reply *cometclient.PushMsgReply, err error) {
	var p protocol.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	keys, err := l.svcCtx.KeysWithServersByMid(c, appId, toIds)
	if err != nil {
		return
	}

	for server, sKeys := range keys {
		if reply, err = l.svcCtx.CometRPC.PushMsg(context.WithValue(c, "TODO", server), &cometclient.PushMsgReq{Keys: sKeys, Proto: &p}); err != nil {
			return
		}
	}
	return
}
