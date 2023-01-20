package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/app/logic/internal/config"
)

const (
	MidFmt             = "%s:%v"       // {appId}:{uid}
	_prefixMidServer   = "mid_%s:%v"   // mid_{appId}:{uid} -> key:server
	_prefixKeyServer   = "key_%s"      // key_{key} -> server
	_prefixKeyUser     = "usr_%s"      // usr_{key} -> {appId}:{uid}
	_prefixGroupServer = "group_%s:%s" // group_{appId}:{gid} -> {server}:{score}
)

func keyMidServer(appId string, mid string) string {
	return fmt.Sprintf(_prefixMidServer, appId, mid)
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

func (repo *LogicRepositoryRedis) pingRedis(c context.Context) (err error) {
	conn := repo.redis.Get()
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

func (repo *LogicRepositoryRedis) GetMember(c context.Context, key string) (appId string, mid string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	ss, err := redis.String(conn.Do("GET", keyKeyUser(key)))
	if err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.DO(GET)")
		return "", "", err
	}
	arr := strings.Split(ss, ":")
	if len(arr) != 2 {
		return "", "", errors.New("invalid key")
	}
	appId = arr[0]
	mid = arr[1]
	return
}

func (repo *LogicRepositoryRedis) GetServer(c context.Context, key string) (server string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	if server, err = redis.String(conn.Do("GET", keyKeyServer(key))); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.DO(GET)")
	}
	return
}

// ServersByKeys get a server by key.
func (repo *LogicRepositoryRedis) ServersByKeys(c context.Context, keys []string) (res []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, keyKeyServer(key))
	}
	if res, err = redis.Strings(conn.Do("MGET", args...)); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("conn.Do(MGET %v) error", args))
	}
	return
}

func (repo *LogicRepositoryRedis) AddMapping(c context.Context, mid string, appId string, key string, server string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 4
	if mid != "" {
		if err = conn.Send("HSET", keyMidServer(appId, mid), key, server); err != nil {
			log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(HSET %s,%s,%s,%s) error", appId, mid, server, key))
			return
		}
		if err = conn.Send("EXPIRE", keyMidServer(appId, mid), repo.redisExpire); err != nil {
			log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s,%s,%s,%s)", appId, mid, key, server))
			return
		}
		n += 2
	}
	if err = conn.Send("SET", keyKeyServer(key), server); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(HSET %s,%s,%s) error", mid, server, key))
		return
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), repo.redisExpire); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s,%s,%s) error", mid, key, server))
		return
	}
	user := fmt.Sprintf(MidFmt, appId, mid)
	if err = conn.Send("SET", keyKeyUser(key), user); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(HSET %s,%s,%s) error", mid, appId, key))
		return
	}
	if err = conn.Send("EXPIRE", keyKeyUser(key), repo.redisExpire); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s,%s,%s) error", mid, appId, key))
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.Flush() error")
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error().Str("key", key).Err(err).Msg("conn.Receive() error")
			return
		}
	}
	return
}

// ExpireMapping expire a mapping.
func (repo *LogicRepositoryRedis) ExpireMapping(c context.Context, mid string, appId string, key string) (has bool, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 2
	if mid != "" {
		if err = conn.Send("EXPIRE", keyMidServer(appId, mid), repo.redisExpire); err != nil {
			log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s) error", keyMidServer(appId, mid)))
			return
		}
		n++
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), repo.redisExpire); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s) error", keyKeyServer(key)))
		return
	}
	if err = conn.Send("EXPIRE", keyKeyUser(key), repo.redisExpire); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s) error", keyKeyServer(key)))
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.Flush() error")
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			log.Error().Str("key", key).Err(err).Msg("conn.Receive() error")
			return
		}
	}
	return
}

func (repo *LogicRepositoryRedis) DelMapping(c context.Context, mid string, appId string, key string) (has bool, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 2
	if mid != "" {
		if err = conn.Send("HDEL", keyMidServer(appId, mid), key); err != nil {
			log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(HDEL %s) error", keyMidServer(appId, mid)))
			return
		}
		n++
	}
	if err = conn.Send("DEL", keyKeyServer(key)); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(DEL %s) error", keyKeyServer(key)))
		return
	}
	if err = conn.Send("DEL", keyKeyUser(key)); err != nil {
		log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(DEL %s) error", keyKeyUser(key)))
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.Flush() error")
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			log.Error().Str("key", key).Err(err).Msg("conn.Receive() error")
			return
		}
	}
	return
}

// KeysByMids get a key server by mid.
func (repo *LogicRepositoryRedis) KeysByMids(c context.Context, appId string, mids []string) (ress map[string]string, olMids []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	ress = make(map[string]string)
	for _, mid := range mids {
		if err = conn.Send("HGETALL", keyMidServer(appId, mid)); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("conn.Do(HGETALL %s) error", mid))
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error().Err(err).Msg("conn.Flush() error")
		return
	}
	for idx := 0; idx < len(mids); idx++ {
		var (
			res map[string]string
		)
		if res, err = redis.StringMap(conn.Receive()); err != nil {
			log.Error().Err(err).Msg("conn.Receive() error")
			return
		}
		if len(res) > 0 {
			olMids = append(olMids, mids[idx])
		}
		for k, v := range res {
			ress[k] = v
		}
	}
	return
}

//groups
func (repo *LogicRepositoryRedis) IncGroupServer(c context.Context, appId, key, server string, gid []string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 0
	for _, g := range gid {
		if err = conn.Send("ZINCRBY", keyGroupServer(appId, g), "1", server); err != nil {
			log.Error().Str("key", key).Err(err).Msg(
				fmt.Sprintf("conn.Send(ZINCRBY %s,%s,%s) error", keyGroupServer(appId, g), "1", server))
			return
		}
		//if err = conn.Send("EXPIRE", keyGroupServer(appId, g), repo.redisExpire); err != nil {
		//	log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s,%s,%s,%s)", appId, mid, key, server))
		//	return
		//}
		n++
	}
	if err = conn.Flush(); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.Flush() error")
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error().Str("key", key).Err(err).Msg("conn.Receive() error")
			return
		}
	}
	return
}

func (repo *LogicRepositoryRedis) DecGroupServer(c context.Context, appId, key, server string, gid []string) (err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	var n = 0
	for _, g := range gid {
		if err = conn.Send("ZINCRBY", keyGroupServer(appId, g), "-1", server); err != nil {
			log.Error().Str("key", key).Err(err).Msg(
				fmt.Sprintf("conn.Send(ZINCRBY %s,%s,%s) error", keyGroupServer(appId, g), "-1", server))
			return
		}
		//if err = conn.Send("EXPIRE", keyGroupServer(appId, g), repo.redisExpire); err != nil {
		//	log.Error().Str("key", key).Err(err).Msg(fmt.Sprintf("conn.Send(EXPIRE %s,%s,%s,%s)", appId, mid, key, server))
		//	return
		//}
		n++
	}
	if err = conn.Flush(); err != nil {
		log.Error().Str("key", key).Err(err).Msg("conn.Flush() error")
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error().Str("key", key).Err(err).Msg("conn.Receive() error")
			return
		}
	}
	return
}

// KeysByMids get a key server by mid.
func (repo *LogicRepositoryRedis) ServersByGid(c context.Context, appId string, gid string) (res []string, err error) {
	conn := repo.redis.Get()
	defer conn.Close()
	res = make([]string, 0)
	ress := make(map[string]string)
	if ress, err = redis.StringMap(conn.Do("ZRANGE", keyGroupServer(appId, gid), "0", "-1", "WITHSCORES")); err != nil {
		log.Error().Str("appId", appId).Str("gid", gid).Err(err).Msg(
			fmt.Sprintf("conn.DO(ZRANGE %s,%s,%s,%s) error", keyGroupServer(appId, gid), "0", "-1", "WITHSCORES"))
	}
	for k := range ress {
		res = append(res, k)
	}
	return
}
