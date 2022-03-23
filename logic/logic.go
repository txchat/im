package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im-pkg/trace"
	comet "github.com/txchat/im/api/comet/grpc"
	"github.com/txchat/im/common"
	"github.com/txchat/im/logic/auth"
	"github.com/txchat/im/logic/conf"
	"github.com/txchat/im/logic/dao"
	"github.com/txchat/im/naming"
	xkey "github.com/txchat/im/naming/balancer/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var ErrInvalidAuthReq = errors.New("ErrInvalidAuthReq")
var ErrInvalidAppId = errors.New("ErrInvalidAppId")

type Logic struct {
	c           *conf.Config
	dao         *dao.Dao
	cometClient comet.CometClient
	apps        map[string]auth.Auth
}

func New(c *conf.Config) (l *Logic) {
	cometClient, err := newCometClient(c)
	if err != nil {
		panic(err)
	}

	l = &Logic{
		c:           c,
		dao:         dao.New(c),
		cometClient: cometClient,
		apps:        make(map[string]auth.Auth),
	}
	l.loadApps()
	return l
}

func newCometClient(c *conf.Config) (comet.CometClient, error) {
	rb := naming.NewResolver(c.Reg.RegAddrs, c.CometRPCClient.Schema)
	resolver.Register(rb)

	addr := fmt.Sprintf("%s:///%s", c.CometRPCClient.Schema, c.CometRPCClient.SrvName) // "schema://[authority]/service"
	fmt.Println("rpc client call addr:", addr)
	conn, err := common.NewGRPCConnWithKey(addr, time.Duration(c.CometRPCClient.Dial), grpc.WithUnaryInterceptor(trace.OpentracingClientInterceptor))
	if err != nil {
		panic(err)
	}
	return comet.NewCometClient(conn), err
}

func (l *Logic) loadApps() {
	for _, app := range l.c.Apps {
		newAuth, _ := auth.Load(app.AppId)
		if newAuth == nil {
			panic("exec auth not exist:" + app.AppId)
		}
		exec := newAuth(app.AuthUrl, time.Duration(app.Timeout))
		l.apps[app.AppId] = exec
	}
}

func (l *Logic) Close() {
	l.dao.Close()
}

// Connect connected a conn.
func (l *Logic) Connect(c context.Context, server string, p *comet.Proto) (mid string, appId string, key string, hb int64, errMsg string, err error) {
	var (
		authMsg comet.AuthMsg
		bytes   []byte
	)
	log.Info().Str("appId", authMsg.AppId).Bytes("body", p.Body).Str("token", authMsg.Token).Msg("call Connect")
	err = proto.Unmarshal(p.Body, &authMsg)
	if err != nil {
		return
	}

	if authMsg.AppId == "" || authMsg.Token == "" {
		err = ErrInvalidAuthReq
		return
	}
	log.Info().Str("appId", authMsg.AppId).Str("token", authMsg.Token).Msg("call auth")
	appId = authMsg.AppId
	authExec, _ := l.apps[authMsg.AppId]
	if authExec == nil {
		err = ErrInvalidAppId
		return
	}

	mid, errMsg, err = authExec.DoAuth(authMsg.Token, authMsg.Ext)
	if err != nil {
		return
	}

	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	key = uuid.New().String() //连接标识
	if err = l.dao.AddMapping(c, mid, appId, key, server); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", mid, key, server))
		return
	}
	log.Info().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn connected")
	// notify biz user connected
	bytes, err = proto.Marshal(p)
	if err != nil {
		return
	}
	err = l.dao.PublishMsg(c, appId, mid, comet.Op_Auth, key, bytes)
	return
}

// Disconnect disconnect a conn.
func (l *Logic) Disconnect(c context.Context, key string, server string) (has bool, err error) {
	appId, mid, err := l.dao.GetMember(c, key)
	if err != nil {
		return false, err
	}
	if has, err = l.dao.DelMapping(c, mid, appId, key); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.DelMapping(%s,%s) error", mid, appId))
		return
	}
	log.Info().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn disconnected")
	// notify biz user disconnected
	l.dao.PublishMsg(c, appId, mid, comet.Op_Disconnect, key, nil)
	return
}

// Heartbeat heartbeat a conn.
func (l *Logic) Heartbeat(c context.Context, key string, server string) (err error) {
	appId, mid, err := l.dao.GetMember(c, key)
	if err != nil {
		return err
	}

	var has bool
	has, err = l.dao.ExpireMapping(c, mid, appId, key)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("l.dao.ExpireMapping(%s,%s) error", mid, appId))
		return
	}
	if !has {
		if err = l.dao.AddMapping(c, mid, appId, key, server); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("l.dao.AddMapping(%s,%s,%s) error", mid, appId, server))
			return
		}
	}
	log.Debug().Str("key", key).Str("mid", mid).Str("appId", appId).Str("comet", server).Msg("conn heartbeat")
	return
}

// Receive receive a message from client.
func (l *Logic) Receive(c context.Context, key string, p *comet.Proto) (err error) {
	appId, mid, err := l.dao.GetMember(c, key)
	if err != nil {
		return err
	}

	log.Debug().Str("appId", appId).Str("mid", mid).Interface("proto", p).Msg("receive proto")
	msg, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	return l.dao.PublishMsg(c, appId, mid, comet.Op(p.Op), key, msg)
}

// PushByMids Push push a biz message to client.
func (l *Logic) PushByMids(c context.Context, appId string, toIds []string, msg []byte) (reply *comet.PushMsgReply, err error) {
	log.Debug().Str("appId", appId).Strs("mids", toIds).Msg("PushByMids start")

	var p comet.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	server2keys, err := l.getServerByMids(c, appId, toIds)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		log.Debug().Strs("keys", keys).Str("server", server).Msg("PushByMids pushing")
		reply, err = l.cometClient.PushMsg(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.PushMsgReq{Keys: keys, Proto: &p})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("PushByMids l.cometClient.PushMsg")
		}
	}

	return
}

// PushByKeys Push push a biz message to client.
func (l *Logic) PushByKeys(c context.Context, appId string, keys []string, msg []byte) (reply *comet.PushMsgReply, err error) {
	log.Debug().Str("appId", appId).Strs("keys", keys).Msg("PushByKeys start")

	var p comet.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	server2keys, err := l.getServerByKeys(c, keys)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		log.Debug().Strs("keys", keys).Str("server", server).Msg("PushByMids pushing")
		reply, err = l.cometClient.PushMsg(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.PushMsgReq{Keys: keys, Proto: &p})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("PushByMids l.cometClient.PushMsg")
		}
	}

	return
}

// PushGroup Push push a biz message to client.
func (l *Logic) PushGroup(c context.Context, appId string, group string, msg []byte) (reply *comet.BroadcastGroupReply, err error) {
	log.Debug().Str("appId", appId).Str("group", group).Msg("PushGroup start")
	var p comet.Proto
	err = proto.Unmarshal(msg, &p)
	if err != nil {
		return
	}

	cometServers, err := l.dao.ServersByGid(c, appId, group)
	if err != nil {
		return
	}

	for _, server := range cometServers {
		log.Debug().Str("appId", appId).Str("group", group).Str("server", server).Msg("PushGroup pushing")
		reply, err = l.cometClient.BroadcastGroup(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.BroadcastGroupReq{GroupID: appGid(appId, group), Proto: &p})
		if err != nil {
			log.Error().Err(err).Str("appId", appId).Str("group", group).Str("server", server).Msg("PushGroup l.cometClient.BroadcastGroup")
		}
	}

	return
}

// JoinGroupsByMids Push push a biz message to client.
func (l *Logic) JoinGroupsByMids(c context.Context, appId string, mids []string, gids []string) (reply *comet.JoinGroupsReply, err error) {
	log.Debug().Str("appId", appId).Strs("mids", mids).Strs("group", gids).Msg("JoinGroupsByMids start")

	server2keys, err := l.getServerByMids(c, appId, mids)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		for _, key := range keys {
			_ = l.dao.IncGroupServer(c, appId, key, server, gids)
		}

		log.Debug().Strs("keys", keys).Str("server", server).Msg("JoinGroupsByMids pushing")
		reply, err = l.cometClient.JoinGroups(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.JoinGroupsReq{Keys: keys, Gid: appGids(appId, gids)})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("JoinGroupsByMids l.cometClient.JoinGroups")
		}
	}

	return
}

// JoinGroupsByKeys Push push a biz message to client.
func (l *Logic) JoinGroupsByKeys(c context.Context, appId string, keys []string, gids []string) (reply *comet.JoinGroupsReply, err error) {
	log.Debug().Str("appId", appId).Strs("keys", keys).Strs("group", gids).Msg("JoinGroupsByKeys start")

	server2keys, err := l.getServerByKeys(c, keys)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		for _, key := range keys {
			_ = l.dao.IncGroupServer(c, appId, key, server, gids)
		}

		log.Debug().Strs("keys", keys).Str("server", server).Msg("JoinGroupsByKeys pushing")
		reply, err = l.cometClient.JoinGroups(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.JoinGroupsReq{Keys: keys, Gid: appGids(appId, gids)})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("JoinGroupsByKeys l.cometClient.JoinGroups")
		}
	}

	return
}

// LeaveGroupsByMids Push push a biz message to client.
func (l *Logic) LeaveGroupsByMids(c context.Context, appId string, mids []string, gids []string) (reply *comet.LeaveGroupsReply, err error) {
	log.Debug().Str("appId", appId).Strs("mids", mids).Strs("group", gids).Msg("LeaveGroupsByMids start")

	server2keys, err := l.getServerByMids(c, appId, mids)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		for _, key := range keys {
			_ = l.dao.DecGroupServer(c, appId, key, server, gids)
		}

		log.Debug().Strs("keys", keys).Str("server", server).Msg("LeaveGroupsByMids pushing")
		reply, err = l.cometClient.LeaveGroups(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.LeaveGroupsReq{Keys: keys, Gid: appGids(appId, gids)})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("LeaveGroupsByMids l.cometClient.LeaveGroups")
		}
	}

	return
}

// LeaveGroupsByKeys Push push a biz message to client.
func (l *Logic) LeaveGroupsByKeys(c context.Context, appId string, keys []string, gids []string) (reply *comet.LeaveGroupsReply, err error) {
	log.Debug().Str("appId", appId).Strs("keys", keys).Strs("group", gids).Msg("LeaveGroupsByKeys start")

	server2keys, err := l.getServerByKeys(c, keys)
	if err != nil {
		return
	}

	for server, keys := range server2keys {
		for _, key := range keys {
			_ = l.dao.DecGroupServer(c, appId, key, server, gids)
		}

		log.Debug().Strs("keys", keys).Str("server", server).Msg("LeaveGroupsByKeys pushing")
		reply, err = l.cometClient.LeaveGroups(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.LeaveGroupsReq{Keys: keys, Gid: appGids(appId, gids)})
		if err != nil {
			log.Error().Err(err).Strs("keys", keys).Str("server", server).Msg("LeaveGroupsByKeys l.cometClient.LeaveGroups")
		}
	}

	return
}

// DelGroups Push push a biz message to client.
func (l *Logic) DelGroups(c context.Context, appId string, gids []string) (reply *comet.DelGroupsReply, err error) {
	log.Debug().Str("appId", appId).Strs("groups", gids).Msg("DelGroups start")

	server2gids := make(map[string][]string)

	for _, gid := range gids {
		cometServers, err := l.dao.ServersByGid(c, appId, gid)
		if err != nil {
			continue
		}

		for _, server := range cometServers {
			server2gids[server] = append(server2gids[server], gid)
		}
	}

	for server, gids := range server2gids {
		log.Debug().Str("appId", appId).Interface("gids", gids).Str("server", server).Msg("DelGroups pushing")
		reply, err = l.cometClient.DelGroups(context.WithValue(c, xkey.DefaultKey, convertAddress(server)), &comet.DelGroupsReq{Gid: appGids(appId, gids)})
		if err != nil {
			log.Error().Err(err).Str("appId", appId).Strs("groups", gids).Str("server", server).Msg("DelGroups l.cometClient.DelGroups")
		}
	}

	return
}

// getServerByMids 通过 appid 和 mids 找到所在的 comet servers
func (l *Logic) getServerByMids(c context.Context, appId string, mids []string) (map[string][]string, error) {
	ress, _, err := l.dao.KeysByMids(c, appId, mids)
	if err != nil {
		return nil, err
	}

	server2keys := make(map[string][]string)
	for k, s := range ress {
		server2keys[s] = append(server2keys[s], k)
	}

	return server2keys, nil
}

// getServerByKeys 通过 keys 找到所在的 comet servers
func (l *Logic) getServerByKeys(c context.Context, keys []string) (map[string][]string, error) {
	server2keys := make(map[string][]string)
	for _, key := range keys {
		server, _ := l.dao.GetServer(c, key)
		server2keys[server] = append(server2keys[server], key)
	}
	return server2keys, nil
}

func appGid(appId, gid string) string {
	return fmt.Sprintf("%s_%s", appId, gid)
}

func appGids(appId string, gids []string) []string {
	res := make([]string, 0, len(gids))
	for _, g := range gids {
		res = append(res, appGid(appId, g))
	}
	return res
}

func convertAddress(addr string) string {
	// todo
	if len(addr) == 0 {
		return addr
	}
	switch addr[0] {
	case 'g':
		return addr[7:]
	default:
		return addr
	}
}
