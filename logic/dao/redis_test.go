package dao

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/Shopify/sarama.v1"
	"testing"

	"github.com/txchat/im/logic/conf"
)

func TestDao_IncGroupServer(t *testing.T) {
	type fields struct {
		c           *conf.Config
		kafkaPub    sarama.SyncProducer
		redis       *redis.Pool
		redisExpire int32
	}
	type args struct {
		c      context.Context
		appId  string
		key    string
		server string
		gid    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				c:           testConf,
				kafkaPub:    nil,
				redis:       testRedis,
				redisExpire: 0,
			},
			args: args{
				c:      nil,
				appId:  "dtalk",
				key:    "1",
				server: "grpc://172.0.0.1:8080",
				gid:    []string{"1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				c:           tt.fields.c,
				kafkaPub:    tt.fields.kafkaPub,
				redis:       tt.fields.redis,
				redisExpire: tt.fields.redisExpire,
			}
			if err := d.IncGroupServer(tt.args.c, tt.args.appId, tt.args.key, tt.args.server, tt.args.gid); (err != nil) != tt.wantErr {
				t.Errorf("IncGroupServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDao_DecGroupServer(t *testing.T) {
	type fields struct {
		c           *conf.Config
		kafkaPub    sarama.SyncProducer
		redis       *redis.Pool
		redisExpire int32
	}
	type args struct {
		c      context.Context
		appId  string
		key    string
		server string
		gid    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				c:           testConf,
				kafkaPub:    nil,
				redis:       testRedis,
				redisExpire: 0,
			},
			args: args{
				c:      nil,
				appId:  "dtalk",
				key:    "1",
				server: "grpc://172.0.0.1:8080",
				gid:    []string{"1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				c:           tt.fields.c,
				kafkaPub:    tt.fields.kafkaPub,
				redis:       tt.fields.redis,
				redisExpire: tt.fields.redisExpire,
			}
			if err := d.DecGroupServer(tt.args.c, tt.args.appId, tt.args.key, tt.args.server, tt.args.gid); (err != nil) != tt.wantErr {
				t.Errorf("IncGroupServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDao_ServersByGid(t *testing.T) {
	type fields struct {
		c           *conf.Config
		kafkaPub    sarama.SyncProducer
		redis       *redis.Pool
		redisExpire int32
	}
	type args struct {
		c     context.Context
		appId string
		gid   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				c:           nil,
				kafkaPub:    nil,
				redis:       testRedis,
				redisExpire: 0,
			},
			args: args{
				c:     nil,
				appId: "dtalk",
				gid:   "1",
			},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				c:           tt.fields.c,
				kafkaPub:    tt.fields.kafkaPub,
				redis:       tt.fields.redis,
				redisExpire: tt.fields.redisExpire,
			}
			gotRes, err := d.ServersByGid(tt.args.c, tt.args.appId, tt.args.gid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServersByGid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, re := range gotRes {
				t.Logf("got %v:%v\n", i, re)
			}
		})
	}
}

func TestDao_KeysByMids(t *testing.T) {
	type fields struct {
		c           *conf.Config
		kafkaPub    sarama.SyncProducer
		redis       *redis.Pool
		redisExpire int32
	}
	type args struct {
		c     context.Context
		appId string
		mids  []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantRess   map[string]string
		wantOlMids []string
		wantErr    bool
	}{
		{
			name: "",
			fields: fields{
				c:           testConf,
				kafkaPub:    nil,
				redis:       testRedis,
				redisExpire: 0,
			},
			args: args{
				c:     context.Background(),
				appId: "dtalk",
				mids:  []string{"1ygj6Un2UzL2rev6ub6NukWrGcKjW8LoG", "1LNaxM1BtkkRpWEGty8bDxmvWwRwxsCy1B", "14si8HGSBKN2B4Ps7QQJeRLvqWoXHX2NwB", ""},
			},
			wantRess:   nil,
			wantOlMids: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				c:           tt.fields.c,
				kafkaPub:    tt.fields.kafkaPub,
				redis:       tt.fields.redis,
				redisExpire: tt.fields.redisExpire,
			}
			gotRess, gotOlMids, err := d.KeysByMids(tt.args.c, tt.args.appId, tt.args.mids)
			if err != nil {
				t.Error(err)
			}
			t.Log(gotRess)
			t.Log(gotOlMids)
		})
	}
}
