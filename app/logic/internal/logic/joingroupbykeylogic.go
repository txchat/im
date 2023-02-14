package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/model"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupByKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupByKeyLogic {
	return &JoinGroupByKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// JoinGroupByKey user client register comet group.
func (l *JoinGroupByKeyLogic) JoinGroupByKey(in *logic.JoinGroupByKeyReq) (*logic.Reply, error) {
	reply, err := l.joinGroupByKey(l.ctx, in.GetAppId(), in.GetKey(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *JoinGroupByKeyLogic) joinGroupByKey(c context.Context, appId string, keys []string, gids []string) (reply *cometclient.JoinGroupsReply, err error) {
	pushKeys, err := l.svcCtx.KeysWithServers(c, keys)
	if err != nil {
		return
	}
	for server, sKeys := range pushKeys {
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.IncGroupServer(c, appId, key, server, gids)
		}
		if reply, err = l.svcCtx.CometRPC.JoinGroups(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.JoinGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gids)}); err != nil {
			return
		}
	}

	return
}
