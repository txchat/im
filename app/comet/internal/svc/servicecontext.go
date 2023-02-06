package svc

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/internal/config"
	"github.com/txchat/im/app/logic/logicclient"
	"github.com/txchat/im/internal/comet"
	dtask "github.com/txchat/task"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zhenjl/cityhash"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceContext struct {
	Config   config.Config
	LogicRPC logicclient.Logic

	serverID  string
	round     *comet.Round
	buckets   []*comet.Bucket // subkey bucket
	bucketIdx uint32
}

func NewServiceContext(c config.Config) *ServiceContext {
	_, port, _ := net.SplitHostPort(c.ListenOn)
	svc := &ServiceContext{
		Config:   c,
		LogicRPC: logicclient.NewLogic(zrpc.MustNewClient(c.LogicRPC, zrpc.WithNonBlock(), zrpc.WithNonBlock())),
		serverID: fmt.Sprintf("grpc://%s:%v", ip.InternalIP(), port),
		round: comet.NewRound(comet.RoundOptions{
			Reader:       c.TCP.Reader,
			ReadBuf:      c.TCP.ReadBuf,
			ReadBufSize:  c.TCP.ReadBufSize,
			Writer:       c.TCP.Writer,
			WriteBuf:     c.TCP.WriteBuf,
			WriteBufSize: c.TCP.WriteBufSize,
			Timer:        c.Protocol.Timer,
			TimerSize:    c.Protocol.TimerSize,
			Task:         c.Protocol.Task,
			TaskSize:     c.Protocol.TaskSize,
		}),
		buckets:   make([]*comet.Bucket, c.Bucket.Size),
		bucketIdx: uint32(c.Bucket.Size),
	}
	for i := 0; i < c.Bucket.Size; i++ {
		svc.buckets[i] = comet.NewBucket(&c.Bucket)
	}
	return svc
}

// Round return all round.
func (s *ServiceContext) Round() *comet.Round {
	return s.round
}

// Buckets return all buckets.
func (s *ServiceContext) Buckets() []*comet.Bucket {
	return s.buckets
}

// Bucket get the bucket by subkey.
func (s *ServiceContext) Bucket(subKey string) *comet.Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	return s.buckets[idx]
}

// RandServerHearbeat rand server heartbeat.
func (s *ServiceContext) RandServerHearbeat() time.Duration {
	return s.Config.Protocol.MinHeartbeat + time.Duration(rand.Int63n(int64(s.Config.Protocol.MaxHeartbeat-s.Config.Protocol.MinHeartbeat)))
}

// Connect connected a connection.
func (s *ServiceContext) Connect(c context.Context, p *protocol.Proto) (key string, hb time.Duration, errMsg string, err error) {
	reply, err := s.LogicRPC.Connect(c, &logicclient.ConnectReq{
		Server: s.serverID,
		Proto:  p,
	})
	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			return
		}
		if errStatus.Code() == codes.Unauthenticated {
			for _, detail := range errStatus.Details() {
				if d, ok := detail.(*errdetails.ErrorInfo); ok && d.GetReason() == "DEVICE_REJECT" {
					errMsg = d.Metadata["resp_data"]
				}
			}
		}
		return
	}

	return reply.Key, time.Duration(reply.Heartbeat), errMsg, nil
}

// Disconnect disconnected a connection.
func (s *ServiceContext) Disconnect(ctx context.Context, key string) error {
	_, err := s.LogicRPC.Disconnect(ctx, &logicclient.DisconnectReq{
		Server: s.serverID,
		Key:    key,
	})
	return err
}

// Heartbeat a connection session.
func (s *ServiceContext) Heartbeat(ctx context.Context, key string) error {
	_, err := s.LogicRPC.Heartbeat(ctx, &logicclient.HeartbeatReq{
		Server: s.serverID,
		Key:    key,
	})
	return err
}

// Receive receive a message.
func (s *ServiceContext) Receive(ctx context.Context, key string, p *protocol.Proto) error {
	_, err := s.LogicRPC.Receive(ctx, &logicclient.ReceiveReq{Key: key, Proto: p})
	return err
}

// Operate operate.
func (s *ServiceContext) Operate(ctx context.Context, p *protocol.Proto, ch *comet.Channel, tsk *dtask.Task) error {
	switch p.Op {
	case int32(protocol.Op_SendMsg):
		//标明Ack的消息序列
		p.Ack = p.Seq
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
		//p.Op = int32(protocol.Op_SendMsgReply)
	case int32(protocol.Op_ReceiveMsgReply):
		//从task中删除某一条
		if j := tsk.Get(strconv.FormatInt(int64(p.Ack), 10)); j != nil {
			j.Cancel()
		}
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
	case int32(protocol.Op_SyncMsgReq):
		err := s.Receive(ctx, ch.Key, p)
		if err != nil {
			//下层业务调用失败，返回error的话会直接断开连接
			return err
		}
		p.Op = int32(protocol.Op_SyncMsgReply)
	default:
		return s.Receive(ctx, ch.Key, p)
	}
	return nil
}
