package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/model"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelGroupLogic {
	return &DelGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelGroupLogic) DelGroup(in *logic.DelGroupReq) (*logic.Reply, error) {
	reply, err := l.delGroup(l.ctx, in.GetAppId(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *DelGroupLogic) delGroup(ctx context.Context, appId string, gid []string) (reply *comet.DelGroupsReply, err error) {
	servers, err := l.svcCtx.Repo.ServersByGids(ctx, appId, gid)
	if err != nil {
		return
	}

	for server, g := range servers {
		reply, err = l.svcCtx.CometRPC.DelGroups(context.WithValue(ctx, model.CtxKeyTODO, server), &cometclient.DelGroupsReq{
			Gid: l.svcCtx.CometGroupsID(appId, g),
		})
		if err != nil {
			return
		}
	}
	return
}
