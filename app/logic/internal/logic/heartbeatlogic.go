package logic

import (
	"context"
	"fmt"

	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Heartbeat heartbeat a conn.
func (l *HeartbeatLogic) Heartbeat(in *logic.HeartbeatReq) (*logic.Reply, error) {
	if err := l.heartbeat(l.ctx, in.GetKey(), in.GetServer()); err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true}, nil
}

func (l *HeartbeatLogic) heartbeat(c context.Context, key string, server string) (err error) {
	appId, uid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return err
	}

	var has bool
	has, err = l.svcCtx.Repo.ExpireMapping(c, uid, appId, key)
	if err != nil {
		l.Error(fmt.Sprintf("l.dao.ExpireMapping(%s,%s) error", uid, appId), "err", err)
		return
	}
	if !has {
		if err = l.svcCtx.Repo.AddMapping(c, uid, appId, key, server); err != nil {
			l.Error(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", uid, appId, server), "err", err)
			return
		}
	}
	l.Slow("conn heartbeat", "key", key, "uid", uid, "appId", appId, "comet server", server)
	return
}
