package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	xlog "github.com/txchat/im-pkg/log"
	"github.com/txchat/im-pkg/trace"
	"github.com/txchat/im/comet"
	"github.com/txchat/im/comet/conf"
	"github.com/txchat/im/comet/grpc"
	"github.com/txchat/im/comet/http"
	"github.com/txchat/im/naming"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	srvName = "comet"
)

var (
	// projectVersion 项目版本
	projectVersion = "3.0.3"
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
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	srv := comet.New(conf.Conf)
	rpcSrv := grpc.New(conf.Conf.RPCServer, srv)
	httpSrv := http.Start(":8000", srv)

	if err := comet.InitWebsocket(srv, conf.Conf.Websocket.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	if err := comet.InitTCP(srv, conf.Conf.TCP.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}

	// register comet
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
			rpcSrv.GracefulStop()
			srv.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("http server shutdown")
			}
			if err := tracerCloser.Close(); err != nil {
				log.Error().Err(err).Msg("tracer close failed")
			}
			log.Info().Msg("comet exit")
			xlog.Close()
			//log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
