package grpc

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im-pkg/trace"
	pb "github.com/txchat/im/api/logic/grpc"
	"github.com/txchat/im/logic"
	"github.com/txchat/im/logic/conf"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" //register gzip compressor
	"google.golang.org/grpc/keepalive"
)

func New(c *conf.RPCServer, l *logic.Logic) *grpc.Server {
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(c.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(c.ForceCloseWait),
		Time:                  time.Duration(c.KeepAliveInterval),
		Timeout:               time.Duration(c.KeepAliveTimeout),
		MaxConnectionAge:      time.Duration(c.MaxLifeTime),
	})
	connectionTimeout := grpc.ConnectionTimeout(time.Duration(c.Timeout))
	srv := grpc.NewServer(keepParams, connectionTimeout,
		grpc.ChainUnaryInterceptor(
			trace.OpentracingServerInterceptor,
		))
	pb.RegisterLogicServer(srv, &server{srv: l})
	lis, err := net.Listen(c.Network, c.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	pb.UnimplementedLogicServer
	srv *logic.Logic
}

var _ pb.LogicServer = &server{}

// Connect connect a conn.
func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	mid, appId, key, hb, errMsg, err := s.srv.Connect(ctx, req.Server, req.Proto)
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
				return &pb.ConnectReply{}, err
			}
			log.Debug().Err(s.Err()).Msg("error Unauthenticated")
			return &pb.ConnectReply{}, s.Err()
		}
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{Key: key, AppId: appId, Mid: mid, Heartbeat: hb}, nil
}

// Disconnect disconnect a conn.
func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.Reply, error) {
	_, err := s.srv.Disconnect(ctx, req.Key, req.Server)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true}, nil
}

// Heartbeat beartbeat a conn.
func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.Reply, error) {
	if err := s.srv.Heartbeat(ctx, req.Key, req.Server); err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true}, nil
}

// Receive receive a message from client.
func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (*pb.Reply, error) {
	if err := s.srv.Receive(ctx, req.Key, req.Proto); err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true}, nil
}

// Push push a biz message to client.
func (s *server) PushByMids(ctx context.Context, req *pb.MidsMsg) (*pb.Reply, error) {
	reply, err := s.srv.PushByMids(ctx, req.AppId, req.ToIds, req.Msg)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) PushByKeys(ctx context.Context, req *pb.KeysMsg) (*pb.Reply, error) {
	reply, err := s.srv.PushByKeys(ctx, req.AppId, req.ToKeys, req.Msg)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) PushGroup(ctx context.Context, req *pb.GroupMsg) (*pb.Reply, error) {
	reply, err := s.srv.PushGroup(ctx, req.AppId, req.Group, req.Msg)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) JoinGroupsByKeys(ctx context.Context, req *pb.GroupsKey) (*pb.Reply, error) {
	reply, err := s.srv.JoinGroupsByKeys(ctx, req.AppId, req.Keys, req.Gid)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) JoinGroupsByMids(ctx context.Context, req *pb.GroupsMid) (*pb.Reply, error) {
	reply, err := s.srv.JoinGroupsByMids(ctx, req.AppId, req.Mids, req.Gid)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) LeaveGroupsByKeys(ctx context.Context, req *pb.GroupsKey) (*pb.Reply, error) {
	reply, err := s.srv.LeaveGroupsByKeys(ctx, req.AppId, req.Keys, req.Gid)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) LeaveGroupsByMids(ctx context.Context, req *pb.GroupsMid) (*pb.Reply, error) {
	reply, err := s.srv.LeaveGroupsByMids(ctx, req.AppId, req.Mids, req.Gid)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}

// Push push a biz message to client.
func (s *server) DelGroups(ctx context.Context, req *pb.DelGroupsReq) (*pb.Reply, error) {
	reply, err := s.srv.DelGroups(ctx, req.AppId, req.Gid)
	if err != nil {
		return nil, err
	}
	msg, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}
	return &pb.Reply{IsOk: true, Msg: msg}, nil
}
