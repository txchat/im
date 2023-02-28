package dao

import (
	"context"
)

type LogicRepository interface {
	GetMember(ctx context.Context, key string) (appId string, uid string, err error)
	GetServer(ctx context.Context, key string) (server string, err error)
	AddMapping(ctx context.Context, uid string, appId string, key string, server string) (err error)
	ExpireMapping(ctx context.Context, uid string, appId string, key string) (has bool, err error)
	DelMapping(ctx context.Context, uid string, appId string, key string) (has bool, err error)
	ServersByKeys(ctx context.Context, keys []string) (res []string, err error)
	KeysByUIDs(ctx context.Context, appId string, uid []string) (ress map[string]string, onlyUID []string, err error)
	IncGroupServer(ctx context.Context, appId, key, server string, gid []string) (err error)
	DecGroupServer(ctx context.Context, appId, key, server string, gid []string) (err error)
	ServersByGid(ctx context.Context, appId string, gid string) (res []string, err error)
	ServersByGids(ctx context.Context, appId string, gids []string) (ress map[string][]string, err error)
}
