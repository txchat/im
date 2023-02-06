package logic

import (
	"context"

	"github.com/txchat/im/internal/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupsLogic {
	return &JoinGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupsLogic) JoinGroups(in *comet.JoinGroupsReq) (*comet.JoinGroupsReply, error) {
	if len(in.GetKeys()) == 0 || len(in.GetGid()) == 0 {
		return nil, errors.ErrJoinGroupArg
	}
	for _, key := range in.GetKeys() {
		bucket := l.svcCtx.Bucket(key)
		channel := bucket.Channel(key)
		if channel == nil {
			l.Error("JoinGroups can not found channel", "key", key)
			continue
		}
		for _, gid := range in.GetGid() {
			group := bucket.Group(gid)
			if group == nil {
				group, _ = bucket.PutGroup(gid)
			}
			err := group.Put(channel)
			if err != nil {
				l.Error("JoinGroups can not put channel", "key", key,
					"gid", gid, "channel id", channel.Key, "channel ip", channel.IP, "channel port", channel.Port)
				continue
			}
		}
	}
	return &comet.JoinGroupsReply{}, nil
}
