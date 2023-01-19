package logic

import (
	"context"

	"github.com/txchat/im/internel/errors"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupsLogic {
	return &LeaveGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveGroupsLogic) LeaveGroups(in *comet.LeaveGroupsReq) (*comet.LeaveGroupsReply, error) {
	if len(in.GetKeys()) == 0 || len(in.GetGid()) == 0 {
		return nil, errors.ErrPushMsgArg
	}
	for _, key := range in.GetKeys() {
		bucket := l.svcCtx.Bucket(key)
		channel := bucket.Channel(key)
		if channel == nil {
			l.Error("LeaveGroups can not found channel", "key", key)
			continue
		}
		for _, gid := range in.GetGid() {
			if group := bucket.Group(gid); group != nil {
				group.Del(channel)
			}
		}
	}
	return &comet.LeaveGroupsReply{}, nil
}
