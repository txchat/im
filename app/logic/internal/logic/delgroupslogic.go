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

type DelGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelGroupsLogic {
	return &DelGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelGroupsLogic) DelGroups(in *logic.DelGroupsReq) (*logic.Reply, error) {
	reply, err := l.delGroups(l.ctx, in.GetAppId(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// DelGroups Push push a biz message to client.
func (l *DelGroupsLogic) delGroups(ctx context.Context, appId string, gids []string) (reply *comet.DelGroupsReply, err error) {
	servers, err := l.svcCtx.Repo.ServersByGids(ctx, appId, gids)
	if err != nil {
		return
	}

	for server, gid := range servers {
		reply, err = l.svcCtx.CometRPC.DelGroups(context.WithValue(ctx, model.CtxKeyTODO, server), &cometclient.DelGroupsReq{
			Gid: l.svcCtx.CometGroupsID(appId, gid),
		})
		if err != nil {
			return
		}
	}
	return
}
