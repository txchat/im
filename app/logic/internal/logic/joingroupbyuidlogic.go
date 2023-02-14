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

type JoinGroupByUIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupByUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupByUIDLogic {
	return &JoinGroupByUIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// JoinGroupByUID user all client register comet group.
func (l *JoinGroupByUIDLogic) JoinGroupByUID(in *logic.JoinGroupByUIDReq) (*logic.Reply, error) {
	reply, err := l.joinGroupByUID(l.ctx, in.GetAppId(), in.GetUid(), in.GetGid())
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true, Msg: msg}, nil
}

func (l *JoinGroupByUIDLogic) joinGroupByUID(c context.Context, appId string, uid []string, gid []string) (reply *cometclient.JoinGroupsReply, err error) {
	keys, err := l.svcCtx.KeysWithServersByUID(c, appId, uid)
	if err != nil {
		return
	}
	for server, sKeys := range keys {
		//todo 优化成批量redis操作
		for _, key := range sKeys {
			_ = l.svcCtx.Repo.IncGroupServer(c, appId, key, server, gid)
		}
		if reply, err = l.svcCtx.CometRPC.JoinGroups(context.WithValue(c, model.CtxKeyTODO, server), &cometclient.JoinGroupsReq{Keys: sKeys, Gid: l.svcCtx.CometGroupsID(appId, gid)}); err != nil {
			return
		}
	}
	return
}
