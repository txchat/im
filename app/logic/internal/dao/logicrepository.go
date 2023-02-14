package dao

import (
	"context"
)

type LogicRepository interface {
	GetMember(c context.Context, key string) (appId string, uid string, err error)
	GetServer(c context.Context, key string) (server string, err error)
	AddMapping(c context.Context, uid string, appId string, key string, server string) (err error)
	ExpireMapping(c context.Context, uid string, appId string, key string) (has bool, err error)
	DelMapping(c context.Context, uid string, appId string, key string) (has bool, err error)
	ServersByKeys(c context.Context, keys []string) (res []string, err error)
	KeysByUIDs(c context.Context, appId string, uid []string) (ress map[string]string, onlyUID []string, err error)
	IncGroupServer(c context.Context, appId, key, server string, gid []string) (err error)
	DecGroupServer(c context.Context, appId, key, server string, gid []string) (err error)
	ServersByGid(c context.Context, appId string, gid string) (res []string, err error)
	ServersByGids(c context.Context, appId string, gids []string) (ress map[string][]string, err error)
}
