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

type LeaveGroupByKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupByKeyLogic {
	return &LeaveGroupByKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LeaveGroupByKey user client unregister comet group.
func (l *LeaveGroupByKeyLogic) LeaveGroupByKey(in *logic.LeaveGroupByKeyReq) (*logic.Reply, error) {
	reply, err := l.leaveGroupByKey(l.ctx, in.GetAppId(), in.GetKey(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *LeaveGroupByKeyLogic) leaveGroupByKey(c context.Context, appId string, keys []string, gid []string) (reply *cometclient.LeaveGroupsReply, err error) {
	pushKeys, err := l.svcCtx.KeysWithServers(c, keys)
	if err != nil {
		return
	}
	for server, sKeys := range pushKeys {
		//todo 优化成批量redis操作
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.DecGroupServer(c, appId, key, server, gid)
		}
		if reply, err = l.svcCtx.CometRPC.LeaveGroups(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.LeaveGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gid)}); err != nil {
			return
		}
	}
	return
}
