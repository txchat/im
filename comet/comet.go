package comet

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"net"
	"time"

	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/txchat/im-pkg/trace"
	comet "github.com/txchat/im/api/comet/grpc"
	logic "github.com/txchat/im/api/logic/grpc"
	"github.com/txchat/im/comet/conf"
	"github.com/txchat/im/common"
	"github.com/txchat/im/naming"
	"github.com/zhenjl/cityhash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

// Comet is comet server.
type Comet struct {
	c         *conf.Config
	round     *Round
	buckets   []*Bucket // subkey bucket
	bucketIdx uint32

	serverID       string
	logicRPCClient logic.LogicClient
}

// NewServer returns a new Server.
func New(c *conf.Config) *Comet {
	s := &Comet{
		c:              c,
		round:          NewRound(c),
		logicRPCClient: newLogicClient(c),
	}
	// init bucket
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(c.Bucket)
	}
	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(c.RPCServer.Addr)
	s.serverID = "grpc://" + addr + ":" + port
	return s
}

func newLogicClient(c *conf.Config) logic.LogicClient {
	rb := naming.NewResolver(c.Reg.RegAddrs, c.LogicRPCClient.Schema)
	resolver.Register(rb)

	addr := fmt.Sprintf("%s:///%s", c.LogicRPCClient.Schema, c.LogicRPCClient.SrvName) // "schema://[authority]/service"
	fmt.Println("rpc client call addr:", addr)

	conn, err := common.NewGRPCConn(addr, time.Duration(c.LogicRPCClient.Dial), grpc.WithUnaryInterceptor(trace.OpentracingClientInterceptor))
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

// Buckets return all buckets.
func (s *Comet) Buckets() []*Bucket {
	return s.buckets
}

// Bucket get the bucket by subkey.
func (s *Comet) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	return s.buckets[idx]
}

// RandServerHearbeat rand server heartbeat.
func (s *Comet) RandServerHearbeat() time.Duration {
	return time.Duration(s.c.Protocol.MinHeartbeat) + time.Duration(rand.Int63n(int64(s.c.Protocol.MaxHeartbeat-s.c.Protocol.MinHeartbeat)))
}

// Close close the server.
func (s *Comet) Close() (err error) {
	return
}

// Connect connected a connection.
func (s *Comet) Connect(c context.Context, p *comet.Proto) (key string, hb time.Duration, errMsg string, err error) {
	var (
		req   logic.ConnectReq
		reply *logic.ConnectReply
	)

	req.Server = s.serverID
	req.Proto = p
	reply, err = s.logicRPCClient.Connect(c, &req)
	if err != nil {
		log.Debug().Err(err).Interface("reply", reply).Msg("logicRPCClient.Connect")
		errStatus, ok := status.FromError(err)
		if !ok {
			return
		}
		if errStatus.Code() == codes.Unauthenticated {
			for _, detail := range errStatus.Details() {
				if d, ok := detail.(*errdetails.ErrorInfo); ok && d.GetReason() == "DEVICE_REJECT" {
					errMsg = d.Metadata["resp_data"]
					log.Debug().Str("errMsg", errMsg).Msg("get errMsg")
				}
			}
		}
		return
	}

	return reply.Key, time.Duration(reply.Heartbeat), errMsg, nil
}

// Disconnect disconnected a connection.
func (s *Comet) Disconnect(c context.Context, key string) (err error) {
	_, err = s.logicRPCClient.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: s.serverID,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
func (s *Comet) Heartbeat(ctx context.Context, key string) (err error) {
	_, err = s.logicRPCClient.Heartbeat(ctx, &logic.HeartbeatReq{
		Server: s.serverID,
		Key:    key,
	})
	return
}

// Receive receive a message.
func (s *Comet) Receive(ctx context.Context, key string, p *comet.Proto) (err error) {
	_, err = s.logicRPCClient.Receive(ctx, &logic.ReceiveReq{Key: key, Proto: p})
	return
}
