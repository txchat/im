package common

import (
	"context"
	"fmt"
	"time"

	xkey "github.com/txchat/im/naming/balancer/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
	grpcBackoffMaxDelay       = time.Second * 3
)

func NewGRPCConn(addr string, timeout time.Duration, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithInitialWindowSize(grpcInitialWindowSize),
		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
		grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepAliveTime,
			Timeout:             grpcKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	}
	dialOpts = append(dialOpts, opts...)
	return grpc.DialContext(ctx, addr, dialOpts...)
}

func NewGRPCConnWithKey(addr string, timeout time.Duration, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithInitialWindowSize(grpcInitialWindowSize),
		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
		grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepAliveTime,
			Timeout:             grpcKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, xkey.Name)),
	}
	dialOpts = append(dialOpts, opts...)
	return grpc.DialContext(ctx, addr, dialOpts...)
}
