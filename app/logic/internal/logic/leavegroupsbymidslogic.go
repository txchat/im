package logic

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupsByMidsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupsByMidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupsByMidsLogic {
	return &LeaveGroupsByMidsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveGroupsByMidsLogic) LeaveGroupsByMids(in *logic.GroupsMid) (*logic.Reply, error) {
	reply, err := l.leaveGroupsByMids(l.ctx, in.GetAppId(), in.GetMids(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

// LeaveGroupsByMids Push push a biz message to client.
func (l *LeaveGroupsByMidsLogic) leaveGroupsByMids(c context.Context, appId string, mids []string, gids []string) (reply *cometclient.LeaveGroupsReply, err error) {
	keys, err := l.svcCtx.KeysWithServersByMid(c, appId, mids)
	if err != nil {
		return
	}

	for server, sKeys := range keys {
		//todo 优化成批量redis操作
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.DecGroupServer(c, appId, key, server, gids)
		}
		if reply, err = l.svcCtx.CometRPC.LeaveGroups(context.WithValue(c, "TODO", server), &cometclient.LeaveGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gids)}); err != nil {
			return
		}
	}

	return
}
