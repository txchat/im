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

type LeaveGroupByUIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupByUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupByUIDLogic {
	return &LeaveGroupByUIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LeaveGroupByUID user all client unregister comet group.
func (l *LeaveGroupByUIDLogic) LeaveGroupByUID(in *logic.LeaveGroupByUIDReq) (*logic.Reply, error) {
	reply, err := l.leaveGroupByUID(l.ctx, in.GetAppId(), in.GetUid(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *LeaveGroupByUIDLogic) leaveGroupByUID(c context.Context, appId string, uid []string, gid []string) (reply *cometclient.LeaveGroupsReply, err error) {
	keys, err := l.svcCtx.KeysWithServersByUID(c, appId, uid)
	if err != nil {
		return
	}

	for server, sKeys := range keys {
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
