// Code generated by goctl. DO NOT EDIT.
// Source: logic.proto

package server

import (
	"context"

	xlogic "github.com/txchat/im/api/logic"

	"github.com/txchat/im/app/logic/internal/logic"
	"github.com/txchat/im/app/logic/internal/svc"
)

type LogicServer struct {
	svcCtx *svc.ServiceContext
	xlogic.UnimplementedLogicServer
}

func NewLogicServer(svcCtx *svc.ServiceContext) *LogicServer {
	return &LogicServer{
		svcCtx: svcCtx,
	}
}

func (s *LogicServer) Connect(ctx context.Context, in *xlogic.ConnectReq) (*xlogic.ConnectReply, error) {
	l := logic.NewConnectLogic(ctx, s.svcCtx)
	return l.Connect(in)
}

func (s *LogicServer) Disconnect(ctx context.Context, in *xlogic.DisconnectReq) (*xlogic.Reply, error) {
	l := logic.NewDisconnectLogic(ctx, s.svcCtx)
	return l.Disconnect(in)
}

func (s *LogicServer) Heartbeat(ctx context.Context, in *xlogic.HeartbeatReq) (*xlogic.Reply, error) {
	l := logic.NewHeartbeatLogic(ctx, s.svcCtx)
	return l.Heartbeat(in)
}

func (s *LogicServer) Receive(ctx context.Context, in *xlogic.ReceiveReq) (*xlogic.Reply, error) {
	l := logic.NewReceiveLogic(ctx, s.svcCtx)
	return l.Receive(in)
}

func (s *LogicServer) SendByUID(ctx context.Context, in *xlogic.SendByUIDReq) (*xlogic.Reply, error) {
	l := logic.NewSendByUIDLogic(ctx, s.svcCtx)
	return l.SendByUID(in)
}

func (s *LogicServer) PushByUID(ctx context.Context, in *xlogic.PushByUIDReq) (*xlogic.Reply, error) {
	l := logic.NewPushByUIDLogic(ctx, s.svcCtx)
	return l.PushByUID(in)
}

func (s *LogicServer) PushByKey(ctx context.Context, in *xlogic.PushByKeyReq) (*xlogic.Reply, error) {
	l := logic.NewPushByKeyLogic(ctx, s.svcCtx)
	return l.PushByKey(in)
}

func (s *LogicServer) PushGroup(ctx context.Context, in *xlogic.PushGroupReq) (*xlogic.Reply, error) {
	l := logic.NewPushGroupLogic(ctx, s.svcCtx)
	return l.PushGroup(in)
}

func (s *LogicServer) JoinGroupByKey(ctx context.Context, in *xlogic.JoinGroupByKeyReq) (*xlogic.Reply, error) {
	l := logic.NewJoinGroupByKeyLogic(ctx, s.svcCtx)
	return l.JoinGroupByKey(in)
}

func (s *LogicServer) JoinGroupByUID(ctx context.Context, in *xlogic.JoinGroupByUIDReq) (*xlogic.Reply, error) {
	l := logic.NewJoinGroupByUIDLogic(ctx, s.svcCtx)
	return l.JoinGroupByUID(in)
}

func (s *LogicServer) LeaveGroupByKey(ctx context.Context, in *xlogic.LeaveGroupByKeyReq) (*xlogic.Reply, error) {
	l := logic.NewLeaveGroupByKeyLogic(ctx, s.svcCtx)
	return l.LeaveGroupByKey(in)
}

func (s *LogicServer) LeaveGroupByUID(ctx context.Context, in *xlogic.LeaveGroupByUIDReq) (*xlogic.Reply, error) {
	l := logic.NewLeaveGroupByUIDLogic(ctx, s.svcCtx)
	return l.LeaveGroupByUID(in)
}

func (s *LogicServer) DelGroup(ctx context.Context, in *xlogic.DelGroupReq) (*xlogic.Reply, error) {
	l := logic.NewDelGroupLogic(ctx, s.svcCtx)
	return l.DelGroup(in)
}
