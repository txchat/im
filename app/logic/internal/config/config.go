package config

import (
	"time"

	xtime "github.com/Terry-Mao/goim/pkg/time"
	xkafka "github.com/txchat/dtalk/pkg/mq/kafka"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CometRPC zrpc.RpcClientConf
	Apps     []*App
	RedisDB  Redis
	Kafka    xkafka.ProducerConfig
	Node     Node
	Backoff  Backoff
}

type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Expire       time.Duration
}

type App struct {
	AppId   string
	AuthUrl string
	Timeout xtime.Duration
}

// Node node config.
type Node struct {
	HeartbeatMax int
	Heartbeat    xtime.Duration
}

// Backoff message.
type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}
