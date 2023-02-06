package logic

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

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

func (l *HeartbeatLogic) Heartbeat(in *logic.HeartbeatReq) (*logic.Reply, error) {
	// todo: add your logic here and delete this line
	if err := l.heartbeat(l.ctx, in.GetKey(), in.GetServer()); err != nil {
		return nil, err
	}
	return &logic.Reply{IsOk: true}, nil
}

// Heartbeat heartbeat a conn.
func (l *HeartbeatLogic) heartbeat(c context.Context, key string, server string) (err error) {
	appId, mid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return err
	}

	var has bool
	has, err = l.svcCtx.Repo.ExpireMapping(c, mid, appId, key)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.ExpireMapping(%s,%s) error", mid, appId))
		return
	}
	if !has {
		if err = l.svcCtx.Repo.AddMapping(c, mid, appId, key, server); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", mid, appId, server))
			return
		}
	}
	log.Debug().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn heartbeat")
	return
}
