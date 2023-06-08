package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	xkafka "github.com/oofpgDLD/kafka-go"
	"github.com/oofpgDLD/kafka-go/trace"
	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/api/protocol"
	"github.com/txchat/im/app/comet/cometclient"
	"github.com/txchat/im/app/logic/internal/config"
	"github.com/txchat/im/app/logic/internal/dao"
	"github.com/txchat/im/pkg/auth"
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
		Producer: xkafka.NewProducer(c.Producer, xkafka.WithProducerInterceptors(trace.ProducerInterceptor)),
		CometRPC: cometclient.NewComet(zrpc.MustNewClient(c.CometRPC, zrpc.WithNonBlock(), zrpc.WithNonBlock())),
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
		exec := newAuth(app.AuthURL, time.Duration(app.Timeout))
		svc.Apps[app.AppId] = exec
	}
}

// PublishReceiveMessage received from A publishing to MQ.
func (s *ServiceContext) PublishReceiveMessage(ctx context.Context, appId, fromId, key string, op protocol.Op, body []byte) error {
	msg := &logic.ReceivedMessage{
		AppId: appId,
		Key:   key,
		From:  fromId,
		Op:    op,
		Body:  body,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("goim-%s-receive", appId)
	_, _, err = s.Producer.Publish(ctx, topic, fromId, data)
	return err
}

// PublishConnection connect and disconnect from comet publishing to MQ.
func (s *ServiceContext) PublishConnection(ctx context.Context, appId string, fromId, key string, op protocol.Op, body []byte) (err error) {
	msg := &logic.ReceivedMessage{
		AppId: appId,
		Key:   key,
		From:  fromId,
		Op:    op,
		Body:  body,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	topic := fmt.Sprintf("goim-%s-connection", appId)
	_, _, err = s.Producer.Publish(ctx, topic, fromId, data)
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

func (s *ServiceContext) KeysWithServersByUID(ctx context.Context, appId string, uid []string) (map[string][]string, error) {
	keyServers, _, err := s.Repo.KeysByUIDs(ctx, appId, uid)
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
