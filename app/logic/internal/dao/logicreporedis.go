package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/txchat/im/app/logic/internal/config"
)

const (
	UIDFmt             = "%s:%v"       // {appId}:{uid}
	_prefixUIDServer   = "uid_%s:%v"   // uid_{appId}:{uid} -> key:server
	_prefixKeyServer   = "key_%s"      // key_{key} -> server
	_prefixKeyUser     = "usr_%s"      // usr_{key} -> {appId}:{uid}
	_prefixGroupServer = "group_%s:%s" // group_{appId}:{gid} -> {server}:{score}
)

func keyUIDServer(appId string, uid string) string {
	return fmt.Sprintf(_prefixUIDServer, appId, uid)
}

func keyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func keyKeyUser(key string) string {
	return fmt.Sprintf(_prefixKeyUser, key)
}

func keyGroupServer(appId, gid string) string {
	return fmt.Sprintf(_prefixGroupServer, appId, gid)
}

func newRedis(c *config.Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: c.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Network, c.Addr,
				redis.DialConnectTimeout(c.DialTimeout),
				redis.DialReadTimeout(c.ReadTimeout),
				redis.DialWriteTimeout(c.WriteTimeout),
				redis.DialPassword(c.Auth),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

type LogicRepositoryRedis struct {
	redis       *redis.Pool
	redisExpire int32
}

func NewLogicRepositoryRedis(cfg config.Redis) *LogicRepositoryRedis {
	return &LogicRepositoryRedis{
		redis:       newRedis(&cfg),
		redisExpire: int32(cfg.Expire / time.Second),
	}
}

func (repo *LogicRepositoryRedis) pingRedis(ctx context.Context) (err error) {
	conn := repo.redis.Get()
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

func (repo *LogicRepositoryRedis) GetMember(ctx context.Context, key string) (appId string, uid string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	ss, err := redis.String(conn.Do("GET", keyKeyUser(key)))
	if err != nil {
		return "", "", fmt.Errorf("conn.DO(GET) error: %v", err)
	}
	arr := strings.Split(ss, ":")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid key %v", key)
	}
	appId = arr[0]
	uid = arr[1]
	return
}

// KeysByUIDs get a key server by uid.
func (repo *LogicRepositoryRedis) KeysByUIDs(ctx context.Context, appId string, uids []string) (ress map[string]string, onlyUID []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	ress = make(map[string]string)
	for _, uid := range uids {
		if err = conn.Send("HGETALL", keyUIDServer(appId, uid)); err != nil {
			err = fmt.Errorf("conn.Do(HGETALL %s) error: %v", uid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		err = fmt.Errorf("conn.Flush() error %v", err)
		return
	}
	for idx := 0; idx < len(uids); idx++ {
		var (
			res map[string]string
		)
		if res, err = redis.StringMap(conn.Receive()); err != nil {
			return
		}
		if len(res) > 0 {
			onlyUID = append(onlyUID, uids[idx])
		}
		for k, v := range res {
			ress[k] = v
		}
	}
	return
}

func (repo *LogicRepositoryRedis) GetServer(ctx context.Context, key string) (server string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	if server, err = redis.String(conn.Do("GET", keyKeyServer(key))); err != nil {
		err = fmt.Errorf("conn.DO(GET) error %v", err)
	}
	return
}

// ServersByKeys get a server by key.
func (repo *LogicRepositoryRedis) ServersByKeys(ctx context.Context, keys []string) (res []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, keyKeyServer(key))
	}
	if res, err = redis.Strings(conn.Do("MGET", args...)); err != nil {
		err = fmt.Errorf("conn.Do(MGET %v) error %v", args, err)
	}
	return
}

func (repo *LogicRepositoryRedis) AddMapping(ctx context.Context, uid string, appId string, key string, server string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 4
	if uid != "" {
		if err = conn.Send("HSET", keyUIDServer(appId, uid), key, server); err != nil {
			err = fmt.Errorf("conn.Send(HSET %s,%s,%s,%s) error %v", appId, uid, server, key, err)
			return
		}
		if err = conn.Send("EXPIRE", keyUIDServer(appId, uid), repo.redisExpire); err != nil {
			err = fmt.Errorf("conn.Send(EXPIRE %s,%s,%s,%s) error %v", appId, uid, key, server, err)
			return
		}
		n += 2
	}
	if err = conn.Send("SET", keyKeyServer(key), server); err != nil {
		err = fmt.Errorf("conn.Send(HSET %s,%s,%s) error %v", uid, server, key, err)
		return
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), repo.redisExpire); err != nil {
		err = fmt.Errorf("conn.Send(EXPIRE %s,%s,%s) error %v", uid, key, server, err)
		return
	}
	user := fmt.Sprintf(UIDFmt, appId, uid)
	if err = conn.Send("SET", keyKeyUser(key), user); err != nil {
		err = fmt.Errorf("conn.Send(HSET %s,%s,%s) error %v", uid, appId, key, err)
		return
	}
	if err = conn.Send("EXPIRE", keyKeyUser(key), repo.redisExpire); err != nil {
		err = fmt.Errorf("conn.Send(EXPIRE %s,%s,%s) error %v", uid, appId, key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

// ExpireMapping expire a mapping.
func (repo *LogicRepositoryRedis) ExpireMapping(ctx context.Context, uid string, appId string, key string) (has bool, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 2
	if uid != "" {
		if err = conn.Send("EXPIRE", keyUIDServer(appId, uid), repo.redisExpire); err != nil {
			err = fmt.Errorf("conn.Send(EXPIRE %s) error %v", keyUIDServer(appId, uid), err)
			return
		}
		n++
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), repo.redisExpire); err != nil {
		err = fmt.Errorf("conn.Send(EXPIRE %s) error %v", keyKeyServer(key), err)
		return
	}
	if err = conn.Send("EXPIRE", keyKeyUser(key), repo.redisExpire); err != nil {
		err = fmt.Errorf("conn.Send(EXPIRE %s) error %v", keyKeyServer(key), err)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			return
		}
	}
	return
}

func (repo *LogicRepositoryRedis) DelMapping(ctx context.Context, uid string, appId string, key string) (has bool, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 2
	if uid != "" {
		if err = conn.Send("HDEL", keyUIDServer(appId, uid), key); err != nil {
			err = fmt.Errorf("conn.Send(HDEL %s) error %v", keyUIDServer(appId, uid), err)
			return
		}
		n++
	}
	if err = conn.Send("DEL", keyKeyServer(key)); err != nil {
		err = fmt.Errorf("conn.Send(DEL %s) error %v", keyKeyServer(key), err)
		return
	}
	if err = conn.Send("DEL", keyKeyUser(key)); err != nil {
		err = fmt.Errorf(fmt.Sprintf("conn.Send(DEL %s) error %v", keyKeyUser(key)), err)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			return
		}
	}
	return
}

func (repo *LogicRepositoryRedis) IncGroupServer(ctx context.Context, appId, key, server string, gid []string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 0
	for _, g := range gid {
		if err = conn.Send("ZINCRBY", keyGroupServer(appId, g), "1", server); err != nil {
			err = fmt.Errorf("conn.Send(ZINCRBY %s,%s,%s) error %v", keyGroupServer(appId, g), "1", server, err)
			return
		}
		n++
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

func (repo *LogicRepositoryRedis) DecGroupServer(ctx context.Context, appId, key, server string, gid []string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 0
	for _, g := range gid {
		if err = conn.Send("ZINCRBY", keyGroupServer(appId, g), "-1", server); err != nil {
			err = fmt.Errorf("conn.Send(ZINCRBY %s,%s,%s) error %v", keyGroupServer(appId, g), "-1", server, err)
			return
		}
		n++
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

// ServersByGid get a key server by uid.
func (repo *LogicRepositoryRedis) ServersByGid(ctx context.Context, appId string, gid string) (res []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	res = make([]string, 0)
	ress := make(map[string]string)
	if ress, err = redis.StringMap(conn.Do("ZRANGE", keyGroupServer(appId, gid), "0", "-1", "WITHSCORES")); err != nil {
		err = fmt.Errorf("conn.DO(ZRANGE %s,%s,%s,%s) error %v", keyGroupServer(appId, gid), "0", "-1", "WITHSCORES", err)
	}
	for k := range ress {
		res = append(res, k)
	}
	return
}

// ServersByGids key=server value=gid list
func (repo *LogicRepositoryRedis) ServersByGids(ctx context.Context, appId string, gids []string) (ress map[string][]string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	ress = make(map[string][]string)
	prefix := fmt.Sprintf("group_%s:", appId)
	n := 0
	for _, gid := range gids {
		err = conn.Send("ZRANGE", keyGroupServer(appId, gid), "0", "-1", "WITHSCORES")
		n++
	}
	if err = conn.Flush(); err != nil {
		err = fmt.Errorf("conn.Flush() error %v", err)
		return
	}
	for n > 0 {
		n--
		var res map[string]string
		res, err = redis.StringMap(conn.Receive())
		if err != nil {
			return
		}
		for server := range res {
			if _, ok := ress[server]; !ok {
				ress[server] = make([]string, 0)
			}
			ress[server] = append(ress[server], strings.TrimPrefix(server, prefix))
		}
	}
	return
}
