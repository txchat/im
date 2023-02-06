package logic

import (
	"context"

	"github.com/txchat/im/app/comet/cometclient"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/model"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupsByKeysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupsByKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupsByKeysLogic {
	return &LeaveGroupsByKeysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveGroupsByKeysLogic) LeaveGroupsByKeys(in *logic.GroupsKey) (*logic.Reply, error) {
	reply, err := l.leaveGroupsByKeys(l.ctx, in.GetAppId(), in.GetKeys(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// LeaveGroupsByKeys Push push a biz message to client.
func (l *LeaveGroupsByKeysLogic) leaveGroupsByKeys(c context.Context, appId string, keys []string, gids []string) (reply *cometclient.LeaveGroupsReply, err error) {
	pushKeys, err := l.svcCtx.KeysWithServers(c, keys)
	if err != nil {
		return
	}
	for server, sKeys := range pushKeys {
		//todo 优化成批量redis操作
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.DecGroupServer(c, appId, key, server, gids)
		}
		if reply, err = l.svcCtx.CometRPC.LeaveGroups(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.LeaveGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gids)}); err != nil {
			return
		}
	}
	return
}
