package logic

import (
	"context"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelGroupsLogic {
	return &DelGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelGroupsLogic) DelGroups(in *comet.DelGroupsReq) (*comet.DelGroupsReply, error) {
	for _, gid := range in.GetGid() {
		for _, bucket := range l.svcCtx.Buckets() {
			if g := bucket.Group(gid); g != nil {
				bucket.DelGroup(g)
			}
		}
	}
	return &comet.DelGroupsReply{}, nil
}
