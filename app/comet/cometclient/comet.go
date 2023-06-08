package cometclient

import (
	"context"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/api/protocol"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ListCastReply    = comet.ListCastReply
	ListCastReq      = comet.ListCastReq
	GroupCastReply   = comet.GroupCastReply
	GroupCastReq     = comet.GroupCastReq
	BroadcastReply   = comet.BroadcastReply
	BroadcastReq     = comet.BroadcastReq
	DelGroupsReply   = comet.DelGroupsReply
	DelGroupsReq     = comet.DelGroupsReq
	JoinGroupsReply  = comet.JoinGroupsReply
	JoinGroupsReq    = comet.JoinGroupsReq
	LeaveGroupsReply = comet.LeaveGroupsReply
	LeaveGroupsReq   = comet.LeaveGroupsReq
	Proto            = protocol.Proto

	Comet interface {
		ListCast(ctx context.Context, in *ListCastReq, opts ...grpc.CallOption) (*ListCastReply, error)
		GroupCast(ctx context.Context, in *GroupCastReq, opts ...grpc.CallOption) (*GroupCastReply, error)
		Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error)
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

func (m *defaultComet) ListCast(ctx context.Context, in *ListCastReq, opts ...grpc.CallOption) (*ListCastReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.ListCast(ctx, in, opts...)
}

func (m *defaultComet) GroupCast(ctx context.Context, in *GroupCastReq, opts ...grpc.CallOption) (*GroupCastReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.GroupCast(ctx, in, opts...)
}

func (m *defaultComet) Broadcast(ctx context.Context, in *BroadcastReq, opts ...grpc.CallOption) (*BroadcastReply, error) {
	client := comet.NewCometClient(m.cli.Conn())
	return client.Broadcast(ctx, in, opts...)
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
