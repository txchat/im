package config

import (
	"time"

	xlog "github.com/txchat/im-pkg/log"

	"github.com/txchat/im/internal/comet"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Node int64 `json:",default=1"`
	zrpc.RpcServerConf
	LogicRPC  zrpc.RpcClientConf
	TCP       TCP `json:"TCP"`
	Websocket Websocket
	Protocol  Protocol
	Bucket    comet.BucketConfig
	Zlog      xlog.Config
}

// TCP is tcp config.
type TCP struct {
	Bind         []string `json:",default=[':3101']"`
	Sndbuf       int      `json:",default=4096"`
	Rcvbuf       int      `json:",default=4096"`
	KeepAlive    bool     `json:",default=false"`
	Reader       int      `json:",default=32"`
	ReadBuf      int      `json:",default=1024"`
	ReadBufSize  int      `json:",default=8192"`
	Writer       int      `json:",default=32"`
	WriteBuf     int      `json:",default=1024"`
	WriteBufSize int      `json:",default=8192"`
}

// Websocket is websocket config.
type Websocket struct {
	Bind        []string `json:",default=['3102']"`
	TLSOpen     bool     `json:",optional"`
	TLSBind     []string `json:",optional"`
	CertFile    string   `json:",optional"`
	PrivateFile string   `json:",optional"`
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int           `json:",default=32"`
	TimerSize        int           `json:",default=2048"`
	Task             int           `json:",default=32"`
	TaskSize         int           `json:",default=2048"`
	SvrProto         int           `json:",default=5"`
	CliProto         int           `json:",default=10"`
	HandshakeTimeout time.Duration `json:",default=5s"`
	MinHeartbeat     time.Duration `json:",default=5s"`
	MaxHeartbeat     time.Duration `json:",default=10m"`
	TaskDuration     time.Duration `json:",default=30m"`
	Rto              time.Duration `json:",default=3s"`
	LRUSize          int           `json:",default=86400"`
}
