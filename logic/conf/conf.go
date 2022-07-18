package conf

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	xtime "github.com/Terry-Mao/goim/pkg/time"
	xlog "github.com/txchat/im-pkg/log"
	"github.com/uber/jaeger-client-go"
	traceConfig "github.com/uber/jaeger-client-go/config"
)

var (
	confPath string
	regAddrs string

	// Conf config
	Conf *Config
)

func init() {
	var (
		defAddrs = os.Getenv("REGADDRS")
	)
	flag.StringVar(&confPath, "conf", "logic.toml", "default config path")
	flag.StringVar(&regAddrs, "reg", defAddrs, "etcd register addrs. eg:127.0.0.1:2379")
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
			SrvName:  "logic",
			RegAddrs: regAddrs,
		},
		CometRPCClient: &RPCClient{
			Schema:  "im",
			SrvName: "comet",
			Dial:    xtime.Duration(time.Second),
			Timeout: xtime.Duration(time.Second)},
		RPCServer: &RPCServer{
			Network:           "tcp",
			Addr:              "3119",
			Timeout:           xtime.Duration(time.Second),
			IdleTimeout:       xtime.Duration(time.Second * 60),
			MaxLifeTime:       xtime.Duration(time.Hour * 2),
			ForceCloseWait:    xtime.Duration(time.Second * 20),
			KeepAliveInterval: xtime.Duration(time.Second * 60),
			KeepAliveTimeout:  xtime.Duration(time.Second * 20),
		},
		Backoff: &Backoff{MaxDelay: 300, BaseDelay: 3, Factor: 1.8, Jitter: 1.3},
	}
}

// Config config.
type Config struct {
	Env            string
	Log            xlog.Config
	Trace          traceConfig.Configuration
	Reg            *Reg
	CometRPCClient *RPCClient
	RPCServer      *RPCServer
	Kafka          *Kafka
	Redis          *Redis
	Node           *Node
	Backoff        *Backoff
	Apps           []*App
}

// Reg is service register/discovery config
type Reg struct {
	Schema   string
	SrvName  string // call
	RegAddrs string // etcd address, be separated from ','
}

type App struct {
	AppId   string
	AuthURL string
	Timeout xtime.Duration
}

// Node node config.
type Node struct {
	HeartbeatMax int
	Heartbeat    xtime.Duration
}

// Backoff backoff.
type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}

// Redis .
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Expire       xtime.Duration
}

// Kafka .
type Kafka struct {
	Topic   string
	Brokers []string
}

// RPCClient is RPC client config.
type RPCClient struct {
	Schema  string
	SrvName string
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
