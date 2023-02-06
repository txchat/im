package logicclient

import (
	"context"

	"github.com/txchat/im/api/logic"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BizMsg        = logic.BizMsg
	ConnectReply  = logic.ConnectReply
	ConnectReq    = logic.ConnectReq
	DelGroupsReq  = logic.DelGroupsReq
	DisconnectReq = logic.DisconnectReq
	GroupMsg      = logic.GroupMsg
	GroupsKey     = logic.GroupsKey
	GroupsMid     = logic.GroupsMid
	HeartbeatReq  = logic.HeartbeatReq
	KeysMsg       = logic.KeysMsg
	MidsMsg       = logic.MidsMsg
	ReceiveReq    = logic.ReceiveReq
	Reply         = logic.Reply

	Logic interface {
		Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectReply, error)
		Disconnect(ctx context.Context, in *DisconnectReq, opts ...grpc.CallOption) (*Reply, error)
		Heartbeat(ctx context.Context, in *HeartbeatReq, opts ...grpc.CallOption) (*Reply, error)
		Receive(ctx context.Context, in *ReceiveReq, opts ...grpc.CallOption) (*Reply, error)
		PushByMids(ctx context.Context, in *MidsMsg, opts ...grpc.CallOption) (*Reply, error)
		PushByKeys(ctx context.Context, in *KeysMsg, opts ...grpc.CallOption) (*Reply, error)
		PushGroup(ctx context.Context, in *GroupMsg, opts ...grpc.CallOption) (*Reply, error)
		JoinGroupsByKeys(ctx context.Context, in *GroupsKey, opts ...grpc.CallOption) (*Reply, error)
		JoinGroupsByMids(ctx context.Context, in *GroupsMid, opts ...grpc.CallOption) (*Reply, error)
		LeaveGroupsByKeys(ctx context.Context, in *GroupsKey, opts ...grpc.CallOption) (*Reply, error)
		LeaveGroupsByMids(ctx context.Context, in *GroupsMid, opts ...grpc.CallOption) (*Reply, error)
		DelGroups(ctx context.Context, in *DelGroupsReq, opts ...grpc.CallOption) (*Reply, error)
	}

	defaultLogic struct {
		cli zrpc.Client
	}
)

func NewLogic(cli zrpc.Client) Logic {
	return &defaultLogic{
		cli: cli,
	}
}

func (m *defaultLogic) Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectReply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.Connect(ctx, in, opts...)
}

func (m *defaultLogic) Disconnect(ctx context.Context, in *DisconnectReq, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.Disconnect(ctx, in, opts...)
}

func (m *defaultLogic) Heartbeat(ctx context.Context, in *HeartbeatReq, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.Heartbeat(ctx, in, opts...)
}

func (m *defaultLogic) Receive(ctx context.Context, in *ReceiveReq, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.Receive(ctx, in, opts...)
}

func (m *defaultLogic) PushByMids(ctx context.Context, in *MidsMsg, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.PushByMids(ctx, in, opts...)
}

func (m *defaultLogic) PushByKeys(ctx context.Context, in *KeysMsg, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.PushByKeys(ctx, in, opts...)
}

func (m *defaultLogic) PushGroup(ctx context.Context, in *GroupMsg, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.PushGroup(ctx, in, opts...)
}

func (m *defaultLogic) JoinGroupsByKeys(ctx context.Context, in *GroupsKey, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.JoinGroupsByKeys(ctx, in, opts...)
}

func (m *defaultLogic) JoinGroupsByMids(ctx context.Context, in *GroupsMid, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.JoinGroupsByMids(ctx, in, opts...)
}

func (m *defaultLogic) LeaveGroupsByKeys(ctx context.Context, in *GroupsKey, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.LeaveGroupsByKeys(ctx, in, opts...)
}

func (m *defaultLogic) LeaveGroupsByMids(ctx context.Context, in *GroupsMid, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.LeaveGroupsByMids(ctx, in, opts...)
}

func (m *defaultLogic) DelGroups(ctx context.Context, in *DelGroupsReq, opts ...grpc.CallOption) (*Reply, error) {
	client := logic.NewLogicClient(m.cli.Conn())
	return client.DelGroups(ctx, in, opts...)
}
