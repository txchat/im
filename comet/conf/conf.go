package conf

import (
	"flag"
	"os"
	"time"

	"github.com/uber/jaeger-client-go"

	"github.com/BurntSushi/toml"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	xlog "github.com/txchat/im-pkg/log"
	traceConfig "github.com/uber/jaeger-client-go/config"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

var (
	confPath   string
	regAddress string

	// Conf config
	Conf *Config
)

func init() {
	var (
		defAddress = os.Getenv("REGADDRS")
	)
	flag.StringVar(&confPath, "conf", "comet.toml", "default config path.")
	flag.StringVar(&regAddress, "reg", defAddress, "etcd register addrs. eg:127.0.0.1:2379")
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

// Default new a config with specified default value.
func Default() *Config {
	return &Config{
		Env: "",
		Log: xlog.Config{
			Level:   "debug",
			Mode:    "console",
			Path:    "",
			Display: "console",
		},
		Trace: traceConfig.Configuration{
			ServiceName: "answer",
			Gen128Bit:   true,
			Sampler: &traceConfig.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &traceConfig.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: "127.0.0.1:6831",
			},
		},
		Reg: &Reg{
			Schema:   "im",
			SrvName:  "comet",
			RegAddrs: regAddress,
		},
		LogicRPCClient: &RPCClient{
			Schema:  "im",
			SrvName: "logic",
			Dial:    xtime.Duration(time.Second),
			Timeout: xtime.Duration(time.Second),
		},
		RPCServer: &RPCServer{
			Network:           "tcp",
			Addr:              ":3109",
			Timeout:           xtime.Duration(time.Second),
			IdleTimeout:       xtime.Duration(time.Second * 60),
			MaxLifeTime:       xtime.Duration(time.Hour * 2),
			ForceCloseWait:    xtime.Duration(time.Second * 20),
			KeepAliveInterval: xtime.Duration(time.Second * 60),
			KeepAliveTimeout:  xtime.Duration(time.Second * 20),
		},
		TCP: &TCP{
			Bind:         []string{":3101"},
			Sndbuf:       4096,
			Rcvbuf:       4096,
			KeepAlive:    false,
			Reader:       32,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       32,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		},
		Websocket: &Websocket{
			Bind: []string{":3102"},
		},
		Protocol: &Protocol{
			Timer:            32,
			TimerSize:        2048,
			Task:             32,
			TaskSize:         2048,
			CliProto:         5,
			SvrProto:         10,
			HandshakeTimeout: xtime.Duration(time.Second * 5),
			TaskDuration:     xtime.Duration(time.Second * 5),
			MinHeartbeat:     xtime.Duration(time.Minute * 10),
			MaxHeartbeat:     xtime.Duration(time.Minute * 30),
			Rto:              xtime.Duration(time.Second * 3),
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			Groups:        1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
	}
}

// Config is comet config.
type Config struct {
	Env            string
	Log            xlog.Config
	Trace          traceConfig.Configuration
	Reg            *Reg
	TCP            *TCP
	Websocket      *Websocket
	Protocol       *Protocol
	Bucket         *Bucket
	LogicRPCClient *RPCClient
	RPCServer      *RPCServer
}

// Reg is service register/discovery config
type Reg struct {
	Schema   string
	SrvName  string // call
	RegAddrs string // etcd address, be separated from ','
}

// RPCClient is RPC client config.
type RPCClient struct {
	Schema  string
	SrvName string // call
	Dial    xtime.Duration
	Timeout xtime.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           xtime.Duration
	IdleTimeout       xtime.Duration
	MaxLifeTime       xtime.Duration
	ForceCloseWait    xtime.Duration
	KeepAliveInterval xtime.Duration
	KeepAliveTimeout  xtime.Duration
}

// TCP is tcp config.
type TCP struct {
	Bind         []string
	Sndbuf       int
	Rcvbuf       int
	KeepAlive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Websocket is websocket config.
type Websocket struct {
	Bind        []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int
	TimerSize        int
	Task             int
	TaskSize         int
	SvrProto         int
	CliProto         int
	HandshakeTimeout xtime.Duration
	MinHeartbeat     xtime.Duration
	MaxHeartbeat     xtime.Duration
	TaskDuration     xtime.Duration
	Rto              xtime.Duration
}

// Bucket is bucket config.
type Bucket struct {
	Size          int
	Channel       int
	Groups        int
	RoutineAmount uint64
	RoutineSize   int
}
