package config

import (
	"time"

	"github.com/txchat/im/internal/comet"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	LogicRPC  zrpc.RpcClientConf
	TCP       TCP
	Websocket Websocket
	Protocol  Protocol
	Bucket    comet.BucketConfig
}

// TCP is tcp config.
type TCP struct {
	Bind         []string `json:",default=[\":3101\"]"`
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
	Bind        []string `json:",default=4096"`
	TLSOpen     bool     `json:",default=4096"`
	TLSBind     []string `json:",default=4096"`
	CertFile    string   `json:",default=4096"`
	PrivateFile string   `json:",default=4096"`
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int
	TimerSize        int
	Task             int
	TaskSize         int
	SvrProto         int
	CliProto         int
	HandshakeTimeout time.Duration
	MinHeartbeat     time.Duration
	MaxHeartbeat     time.Duration
	TaskDuration     time.Duration
	Rto              time.Duration
}
