package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	xkafka "github.com/txchat/dtalk/pkg/mq/kafka"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/config"
	"github.com/txchat/im/app/logic/internal/dao"
	"github.com/txchat/im/internal/auth"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	Repo     dao.LogicRepository
	Apps     map[string]auth.Auth
	Producer *xkafka.Producer
	CometRPC cometclient.Comet
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:   c,
		Repo:     dao.NewLogicRepositoryRedis(c.RedisDB),
		Apps:     make(map[string]auth.Auth),
		Producer: xkafka.NewProducer(c.Kafka),
		CometRPC: cometclient.NewComet(zrpc.MustNewClient(c.CometRPC, zrpc.WithNonBlock())),
	}
	loadApps(c, s)
	return s
}

func loadApps(c config.Config, svc *ServiceContext) {
	for _, app := range c.Apps {
		newAuth, _ := auth.Load(app.AppId)
		if newAuth == nil {
			panic("exec auth not exist:" + app.AppId)
		}
		exec := newAuth(app.AuthUrl, time.Duration(app.Timeout))
		svc.Apps[app.AppId] = exec
	}
}

// PublishMsg push a message to MQ.
func (s *ServiceContext) PublishMsg(ctx context.Context, appId string, fromId string, op protocol.Op, key string, msg []byte) (err error) {
	pushMsg := &logic.BizMsg{
		AppId:  appId,
		FromId: fromId,
		Op:     int32(op),
		Key:    key,
		Msg:    msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}
	topic := fmt.Sprintf("goim-%s-topic", appId)

	_, _, err = s.Producer.Publish(topic, fromId, b)
	return
}

func (s *ServiceContext) KeysWithServers(ctx context.Context, keys []string) (map[string][]string, error) {
	servers, err := s.Repo.ServersByKeys(ctx, keys)
	if err != nil {
		return nil, err
	}
	pushKeys := make(map[string][]string)
	for i, key := range keys {
		server := servers[i]
		if server != "" && key != "" {
			pushKeys[server] = append(pushKeys[server], key)
		}
	}
	return pushKeys, nil
}

func (s *ServiceContext) KeysWithServersByMid(ctx context.Context, appId string, mid []string) (map[string][]string, error) {
	keyServers, _, err := s.Repo.KeysByMids(ctx, appId, mid)
	if err != nil {
		return nil, err
	}

	keys := make(map[string][]string)
	for key, server := range keyServers {
		if key == "" || server == "" {
			continue
		}
		keys[server] = append(keys[server], key)
	}
	return keys, nil
}

func (s *ServiceContext) CometGid(appId string, gid string) string {
	return fmt.Sprintf("%s_%s", appId, gid)
}

func (s *ServiceContext) CometGroupsID(appId string, groups []string) []string {
	res := make([]string, 0, len(groups))
	for _, g := range groups {
		res = append(res, s.CometGid(appId, g))
	}
	return res
}
