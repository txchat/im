package logic

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/protocol"

	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisconnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisconnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisconnectLogic {
	return &DisconnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DisconnectLogic) Disconnect(in *logic.DisconnectReq) (*logic.Reply, error) {
	// todo: add your logic here and delete this line
	_, err := l.disconnect(l.ctx, in.GetKey(), in.GetServer())
	if err != nil {
		return nil, err
	}
	return &logic.Reply{
		IsOk: true,
	}, nil
}

// Disconnect disconnect a conn.
func (l *DisconnectLogic) disconnect(c context.Context, key string, server string) (has bool, err error) {
	appId, mid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return false, err
	}
	if has, err = l.svcCtx.Repo.DelMapping(c, mid, appId, key); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.DelMapping(%s,%s) error", mid, appId))
		return
	}
	log.Info().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn disconnected")
	// notify biz user disconnected
	l.svcCtx.PublishMsg(c, appId, mid, protocol.Op_Disconnect, key, nil)
	return
}
