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

type PushGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushGroupLogic {
	return &PushGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushGroup push message from biz level to group clients.
func (l *PushGroupLogic) PushGroup(in *logic.PushGroupReq) (*logic.Reply, error) {
	reply, err := l.pushGroup(l.ctx, in.GetAppId(), in.GetGroup(), &protocol.Proto{
		Ver:  model.ProtoVersion,
		Op:   int32(in.GetOp()),
		Body: in.GetBody(),
	})
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *PushGroupLogic) pushGroup(c context.Context, appId string, group string, p *protocol.Proto) (reply *cometclient.GroupCastReply, err error) {
	servers, err := l.svcCtx.Repo.ServersByGid(c, appId, group)
	if err != nil {
		return
	}

	for _, server := range servers {
		if reply, err = l.svcCtx.CometRPC.GroupCast(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.GroupCastReq{Gid: l.svcCtx.CometGid(appId, group), Proto: p}); err != nil {
			return
		}
	}
	return
}
