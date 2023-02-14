package logic

import (
	"context"
	"fmt"

	"github.com/txchat/im/internal/auth"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
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
	uid, appId, key, hb, err := l.connect(l.ctx, in.GetServer(), in.GetProto())
	if err != nil {
		l.Error("client connect failed", "err", err, "key", key, "uid", uid, "appId", appId, "comet server", in.GetServer())
		return &logic.ConnectReply{}, auth.ToGRPCErr(err)
	}
	l.Info("client connected", "key", key, "uid", uid, "appId", appId, "comet server", in.GetServer())
	return &logic.ConnectReply{
		Key:       key,
		AppId:     appId,
		Uid:       uid,
		Heartbeat: hb,
	}, nil
}

func (l *ConnectLogic) connect(c context.Context, server string, p *protocol.Proto) (uid string, appId string, key string, hb int64, err error) {
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

	uid, err = authExec.DoAuth(authBody.Token, authBody.Ext)
	if err != nil {
		return
	}

	hb = int64(l.svcCtx.Config.Node.Heartbeat) * int64(l.svcCtx.Config.Node.HeartbeatMax)
	key = uuid.New().String() //连接标识
	if err = l.svcCtx.Repo.AddMapping(c, uid, appId, key, server); err != nil {
		l.Error(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", uid, key, server), "err", err)
		return
	}
	// notify biz user connected
	bytes, err = proto.Marshal(p)
	if err != nil {
		return
	}
	err = l.svcCtx.PublishConnection(c, appId, uid, protocol.Op_Auth, key, bytes)
	return
}
