package server

import (
	"context"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/logic"
	"github.com/txchat/im/app/comet/internal/svc"
)

type CometServer struct {
	svcCtx *svc.ServiceContext
	comet.UnimplementedCometServer
}

func NewCometServer(svcCtx *svc.ServiceContext) *CometServer {
	return &CometServer{
		svcCtx: svcCtx,
	}
}

func (s *CometServer) ListCast(ctx context.Context, in *comet.ListCastReq) (*comet.ListCastReply, error) {
	l := logic.NewListCastLogic(ctx, s.svcCtx)
	return l.ListCast(in)
}

func (s *CometServer) GroupCast(ctx context.Context, in *comet.GroupCastReq) (*comet.GroupCastReply, error) {
	l := logic.NewGroupCastLogic(ctx, s.svcCtx)
	return l.GroupCast(in)
}

func (s *CometServer) Broadcast(ctx context.Context, in *comet.BroadcastReq) (*comet.BroadcastReply, error) {
	l := logic.NewBroadcastLogic(ctx, s.svcCtx)
	return l.Broadcast(in)
}

func (s *CometServer) JoinGroups(ctx context.Context, in *comet.JoinGroupsReq) (*comet.JoinGroupsReply, error) {
	l := logic.NewJoinGroupsLogic(ctx, s.svcCtx)
	return l.JoinGroups(in)
}

func (s *CometServer) LeaveGroups(ctx context.Context, in *comet.LeaveGroupsReq) (*comet.LeaveGroupsReply, error) {
	l := logic.NewLeaveGroupsLogic(ctx, s.svcCtx)
	return l.LeaveGroups(in)
}

func (s *CometServer) DelGroups(ctx context.Context, in *comet.DelGroupsReq) (*comet.DelGroupsReply, error) {
	l := logic.NewDelGroupsLogic(ctx, s.svcCtx)
	return l.DelGroups(in)
}
