package config

import (
	"time"

	xtime "github.com/Terry-Mao/goim/pkg/time"
	xkafka "github.com/oofpgDLD/kafka-go"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CometRPC zrpc.RpcClientConf
	RedisDB  Redis
	Producer xkafka.ProducerConfig
	Apps     []*App
	Node     Node
	Backoff  Backoff
}

type Redis struct {
	Network      string `json:",default=tcp"`
	Addr         string
	Auth         string        `json:",optional"`
	Active       int           `json:",default=60000"`
	Idle         int           `json:",default=1024"`
	DialTimeout  time.Duration `json:",default=200ms"`
	ReadTimeout  time.Duration `json:",default=500ms"`
	WriteTimeout time.Duration `json:",default=500ms"`
	IdleTimeout  time.Duration `json:",default=120s"`
	Expire       time.Duration `json:",default=30m"`
}

type App struct {
	AppId   string
	AuthURL string         `json:"AuthURL"`
	Timeout xtime.Duration `json:",default=5s"`
}

// Node node config.
type Node struct {
	HeartbeatMax int            `json:",default=2"`
	Heartbeat    xtime.Duration `json:",default=4m"`
}

// Backoff message.
type Backoff struct {
	MaxDelay  int32   `json:",default=300"`
	BaseDelay int32   `json:",default=3"`
	Factor    float32 `json:",default=1.8"`
	Jitter    float32 `json:",default=1.3"`
}
