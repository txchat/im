package logic

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/txchat/im/internal/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (l *ConnectLogic) Connect(in *logic.ConnectReq) (*logic.ConnectReply, error) {
	mid, appId, key, hb, errMsg, err := l.connect(l.ctx, in.GetServer(), in.GetProto())
	if err != nil {
		if errMsg != "" {
			s := status.New(codes.Unauthenticated, "")
			s, err = s.WithDetails(&errdetails.ErrorInfo{
				Reason: "DEVICE_REJECT",
				Domain: "unknown",
				Metadata: map[string]string{
					"resp_data": errMsg,
				},
			})
			if err != nil {
				log.Debug().Err(err).Msg("error of grpc custom handle error")
				return &logic.ConnectReply{}, err
			}
			log.Debug().Err(s.Err()).Msg("error Unauthenticated")
			return &logic.ConnectReply{}, s.Err()
		}
		return &logic.ConnectReply{}, err
	}
	return &logic.ConnectReply{
		Key:       key,
		AppId:     appId,
		Mid:       mid,
		Heartbeat: hb,
	}, nil
}

// Connect connected a conn.
func (l *ConnectLogic) connect(c context.Context, server string, p *protocol.Proto) (mid string, appId string, key string, hb int64, errMsg string, err error) {
	var (
		authMsg comet.AuthMsg
		bytes   []byte
	)
	l.Info()
	log.Info().Str("appId", authMsg.AppId).Bytes("body", p.Body).Str("token", authMsg.Token).Msg("call Connect")
	err = proto.Unmarshal(p.Body, &authMsg)
	if err != nil {
		return
	}

	if authMsg.AppId == "" || authMsg.Token == "" {
		err = errors.ErrInvalidAuthReq
		return
	}
	log.Info().Str("appId", authMsg.AppId).Str("token", authMsg.Token).Msg("call auth")
	appId = authMsg.AppId
	authExec := l.svcCtx.Apps[authMsg.AppId]
	if authExec == nil {
		err = errors.ErrInvalidAppId
		return
	}

	mid, errMsg, err = authExec.DoAuth(authMsg.Token, authMsg.Ext)
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
