package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/comet/cometclient"
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
	// todo: add your logic here and delete this line
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
func (l *DelGroupsLogic) delGroups(c context.Context, appId string, gids []string) (reply *comet.DelGroupsReply, err error) {

	server2gids := make(map[string][]string)

	for _, gid := range gids {
		cometServers, err := l.svcCtx.Repo.ServersByGid(c, appId, gid)
		if err != nil {
			continue
		}

		for _, server := range cometServers {
			server2gids[server] = append(server2gids[server], gid)
		}
	}

	for server, gids := range server2gids {
		log.Debug().Str("appId", appId).Interface("gids", gids).Str("server", server).Msg("DelGroups pushing")
		reply, err = l.svcCtx.CometRPC.DelGroups(l.ctx, &cometclient.DelGroupsReq{
			Gid: appGids(appId, gids),
		})
		if err != nil {
			log.Error().Err(err).Str("appId", appId).Strs("groups", gids).Str("server", server).Msg("DelGroups l.cometClient.DelGroups")
		}
	}

	return
}
