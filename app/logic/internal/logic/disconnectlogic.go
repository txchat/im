package logic

import (
	"context"
	"fmt"

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

// Disconnect disconnect a conn.
func (l *DisconnectLogic) Disconnect(in *logic.DisconnectReq) (*logic.Reply, error) {
	_, err := l.disconnect(l.ctx, in.GetKey(), in.GetServer())
	if err != nil {
		l.Error("client disconnect failed", "err", err, "key", in.GetKey(), "comet server", in.GetServer())
		return nil, err
	}
	l.Info("client disconnected", "key", in.GetKey(), "comet server", in.GetServer())
	return &logic.Reply{
		IsOk: true,
	}, nil
}

func (l *DisconnectLogic) disconnect(c context.Context, key string, server string) (has bool, err error) {
	appId, uid, err := l.svcCtx.Repo.GetMember(c, key)
	if err != nil {
		return false, err
	}
	if has, err = l.svcCtx.Repo.DelMapping(c, uid, appId, key); err != nil {
		l.Error(fmt.Sprintf("l.dao.DelMapping(%s,%s) error", uid, appId), "err", err)
		return
	}
	// notify biz user disconnected
	err = l.svcCtx.PublishConnection(c, appId, uid, protocol.Op_Disconnect, key, nil)
	return
}
