package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	xlog "github.com/txchat/im-pkg/log"
	"github.com/txchat/im-pkg/trace"
	"github.com/txchat/im/logic"
	_ "github.com/txchat/im/logic/auth/dtalk"
	_ "github.com/txchat/im/logic/auth/zhaobi"
	"github.com/txchat/im/logic/conf"
	"github.com/txchat/im/logic/grpc"
	"github.com/txchat/im/naming"
)

const (
	srvName = "logic"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
}

var (
	// projectVersion 项目版本
	projectVersion = "3.1.4"
	// goVersion go版本
	goVersion = ""
	// gitCommit git提交commit id
	gitCommit = ""
	// buildTime 编译时间
	buildTime = ""

	isShowVersion = flag.Bool("version", false, "show project version")
)

// showVersion 显示项目版本信息
func showVersion(isShow bool) {
	if isShow {
		fmt.Printf("Project: %s\n", srvName)
		fmt.Printf(" Version: %s\n", projectVersion)
		fmt.Printf(" Go Version: %s\n", goVersion)
		fmt.Printf(" Git Commit: %s\n", gitCommit)
		fmt.Printf(" Build Time: %s\n", buildTime)
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	showVersion(*isShowVersion)

	if err := conf.Init(); err != nil {
		panic(err)
	}

	//log init
	var err error
	log.Logger, err = xlog.Init(conf.Conf.Log)
	if err != nil {
		panic(err)
	}
	log.Logger.With().Str("service", srvName)

	byte, _ := json.Marshal(conf.Conf)
	log.Info().Str("config", string(byte)).Send()

	// trace init
	tracer, tracerCloser := trace.Init(srvName, conf.Conf.Trace, config.Logger(jaeger.NullLogger))
	//不然后续不会有Jaeger实例
	opentracing.SetGlobalTracer(tracer)

	srv := logic.New(conf.Conf)
	rpcSrv := grpc.New(conf.Conf.RPCServer, srv)

	go func() {
		if err := http.ListenAndServe(":8001", nil); err != nil {
			panic(err)
		}
	}()

	// register logic
	_, port, _ := net.SplitHostPort(conf.Conf.RPCServer.Addr)
	addr := fmt.Sprintf("%s:%s", ip.InternalIP(), port)
	if err := naming.Register(conf.Conf.Reg.RegAddrs, conf.Conf.Reg.SrvName, addr, conf.Conf.Reg.Schema, 15); err != nil {
		panic(err)
	}
	fmt.Println("register ok")

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info().Str("signal", s.String()).Send()
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if err := naming.UnRegister(conf.Conf.Reg.SrvName, addr, conf.Conf.Reg.Schema); err != nil {
				log.Error().Err(err).Msg("naming.UnRegister")
			}
			srv.Close()
			rpcSrv.GracefulStop()
			if err := tracerCloser.Close(); err != nil {
				log.Error().Err(err).Msg("tracer close failed")
			}
			//log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
