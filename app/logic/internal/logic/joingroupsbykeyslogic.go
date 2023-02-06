package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupsByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupsByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupsByKeysLogic {
	return &JoinGroupsByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupsByKeysLogic) JoinGroupsByKeys(in *logic.GroupsKey) (*logic.Reply, error) {
	reply, err := l.joinGroupsByKeys(l.ctx, in.GetAppId(), in.GetKeys(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// JoinGroupsByKeys Push push a biz message to client.
func (l *JoinGroupsByKeysLogic) joinGroupsByKeys(c context.Context, appId string, keys []string, gids []string) (reply *cometclient.JoinGroupsReply, err error) {
	pushKeys, err := l.svcCtx.KeysWithServers(c, keys)
	if err != nil {
		return
	}
	for server, sKeys := range pushKeys {
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.IncGroupServer(c, appId, key, server, gids)
		}
		if reply, err = l.svcCtx.CometRPC.JoinGroups(context.WithValue(c, "TODO", server), &cometclient.JoinGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gids)}); err != nil {
			return
		}
	}

	return
}
