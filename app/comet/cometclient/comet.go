package cometclient

import (
	"context"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/api/protocol"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AuthMsg             = comet.AuthMsg
	BroadcastGroupReply = comet.BroadcastGroupReply
	BroadcastGroupReq   = comet.BroadcastGroupReq
	BroadcastReply      = comet.BroadcastReply
	BroadcastReq        = comet.BroadcastReq
	DelGroupsReply      = comet.DelGroupsReply
	DelGroupsReq        = comet.DelGroupsReq
	Empty               = comet.Empty
	JoinGroupsReply     = comet.JoinGroupsReply
	JoinGroupsReq       = comet.JoinGroupsReq
	LeaveGroupsReply    = comet.LeaveGroupsReply
	LeaveGroupsReq      = comet.LeaveGroupsReq
	Proto               = protocol.Proto
	PushMsgReply        = comet.PushMsgReply
	PushMsgReq          = comet.PushMsgReq
	SyncMsg             = comet.SyncMsg

	Comet interface {
		PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
		Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error)
		BroadcastGroup(ctx context.Context, in *BroadcastGroupReq, opts ...grpc.CallOption) (*BroadcastGroupReply, error)
		JoinGroups(ctx context.Context, in *JoinGroupsReq, opts ...grpc.CallOption) (*JoinGroupsReply, error)
		LeaveGroups(ctx context.Context, in *LeaveGroupsReq, opts ...grpc.CallOption) (*LeaveGroupsReply, error)
		DelGroups(ctx context.Context, in *DelGroupsReq, opts ...grpc.CallOption) (*DelGroupsReply, error)
	}

	defaultComet struct {
		cli zrpc.Client
	}
)

func NewComet(cli zrpc.Client) Comet {
	return &defaultComet{
		cli: cli,
	}
}

func (m *defaultComet) PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.PushMsg(ctx, in, opts...)
}

func (m *defaultComet) Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.Broadcast(ctx, in, opts...)
}

func (m *defaultComet) BroadcastGroup(ctx context.Context, in *BroadcastGroupReq, opts ...grpc.CallOption) (*BroadcastGroupReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.BroadcastGroup(ctx, in, opts...)
}

func (m *defaultComet) JoinGroups(ctx context.Context, in *JoinGroupsReq, opts ...grpc.CallOption) (*JoinGroupsReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.JoinGroups(ctx, in, opts...)
}

func (m *defaultComet) LeaveGroups(ctx context.Context, in *LeaveGroupsReq, opts ...grpc.CallOption) (*LeaveGroupsReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.LeaveGroups(ctx, in, opts...)
}

func (m *defaultComet) DelGroups(ctx context.Context, in *DelGroupsReq, opts ...grpc.CallOption) (*DelGroupsReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.DelGroups(ctx, in, opts...)
}
