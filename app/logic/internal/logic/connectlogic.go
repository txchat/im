package logic

import (
	"context"
	"fmt"

	"github.com/txchat/im/internal/auth"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/internal/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectLogic {
	return &ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Connect connected a conn.
func (l *ConnectLogic) Connect(in *logic.ConnectReq) (*logic.ConnectReply, error) {
	mid, appId, key, hb, err := l.connect(l.ctx, in.GetServer(), in.GetProto())
	if err != nil {
		return &logic.ConnectReply{}, auth.ToGRPCErr(err)
	}
	return &logic.ConnectReply{
		Key:       key,
		AppId:     appId,
		Mid:       mid,
		Heartbeat: hb,
	}, nil
}

func (l *ConnectLogic) connect(c context.Context, server string, p *protocol.Proto) (mid string, appId string, key string, hb int64, err error) {
	var (
		authBody protocol.AuthBody
		bytes    []byte
	)
	err = proto.Unmarshal(p.Body, &authBody)
	if err != nil {
		return
	}

	if authBody.AppId == "" || authBody.Token == "" {
		err = errors.ErrInvalidAuthReq
		return
	}
	appId = authBody.AppId
	authExec := l.svcCtx.Apps[authBody.AppId]
	if authExec == nil {
		err = errors.ErrInvalidAppId
		return
	}

	mid, err = authExec.DoAuth(authBody.Token, authBody.Ext)
	if err != nil {
		return
	}

	hb = int64(l.svcCtx.Config.Node.Heartbeat) * int64(l.svcCtx.Config.Node.HeartbeatMax)
	key = uuid.New().String() //连接标识
	if err = l.svcCtx.Repo.AddMapping(c, mid, appId, key, server); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", mid, key, server))
		return
	}
	log.Info().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn connected")
	// notify biz user connected
	bytes, err = proto.Marshal(p)
	if err != nil {
		return
	}
	err = l.svcCtx.PublishMsg(c, appId, mid, protocol.Op_Auth, key, bytes)
	return
}
