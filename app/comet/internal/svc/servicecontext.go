package svc

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	dtask "github.com/txchat/task"

	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/internal/config"
	"github.com/txchat/im/app/logic/logicclient"
	"github.com/txchat/im/internal/comet"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zhenjl/cityhash"
)

type ServiceContext struct {
	Config   config.Config
	LogicRPC logicclient.Logic

	TaskPool *dtask.Task

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
		TaskPool: dtask.NewTask(),
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

// Close the service resource.
func (s *ServiceContext) Close() {
	s.TaskPool.Stop()
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

// RandServerHeartbeat rand server heartbeat.
func (s *ServiceContext) RandServerHeartbeat() time.Duration {
	return s.Config.Protocol.MinHeartbeat + time.Duration(rand.Int63n(int64(s.Config.Protocol.MaxHeartbeat-s.Config.Protocol.MinHeartbeat)))
}

// Connect connected a connection.
func (s *ServiceContext) Connect(c context.Context, p *protocol.Proto) (key string, hb time.Duration, err error) {
	reply, err := s.LogicRPC.Connect(c, &logicclient.ConnectReq{
		Server: s.serverID,
		Proto:  p,
	})
	if err != nil {
		return
	}

	return reply.Key, time.Duration(reply.Heartbeat), nil
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
