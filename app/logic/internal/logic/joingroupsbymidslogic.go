package logic

import (
	"context"

	"github.com/txchat/im/app/comet/cometclient"
	xkey "github.com/txchat/im/naming/balancer/key"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupsByMidsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupsByMidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupsByMidsLogic {
	return &JoinGroupsByMidsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupsByMidsLogic) JoinGroupsByMids(in *logic.GroupsMid) (*logic.Reply, error) {
	reply, err := l.joinGroupsByMids(l.ctx, in.GetAppId(), in.GetMids(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// JoinGroupsByMids Push push a biz message to client.
func (l *JoinGroupsByMidsLogic) joinGroupsByMids(c context.Context, appId string, mids []string, gids []string) (reply *cometclient.JoinGroupsReply, err error) {
	keys, err := l.svcCtx.KeysWithServersByMid(c, appId, mids)
	if err != nil {
		return
	}
	for server, sKeys := range keys {
		//todo 优化成批量redis操作
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.IncGroupServer(c, appId, key, server, gids)
		}
		if reply, err = l.svcCtx.CometRPC.JoinGroups(context.WithValue(c, xkey.DefaultKey, server), &cometclient.JoinGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gids)}); err != nil {
			return
		}
	}
	return
}
