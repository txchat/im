package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/config"
	"github.com/txchat/im/app/comet/internal/http"
	"github.com/txchat/im/app/comet/internal/server"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/txchat/im/app/comet/internal/tcp"
	"github.com/txchat/im/app/comet/internal/websocket"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// projectName 项目名称
	projectName = "backup"
	// projectVersion 项目版本
	projectVersion = "0.0.1"
	// goVersion go版本
	goVersion = ""
	// gitCommit git提交commit id
	gitCommit = ""
	// buildTime 编译时间
	buildTime = ""
	// osArch 目标主机架构
	osArch = ""
	// isShowVersion 是否显示项目版本信息
	isShowVersion = flag.Bool("v", false, "show project version")
	// configFile 配置文件路径
	configFile = flag.String("f", "etc/comet.yaml", "the config file")
)

func main() {
	flag.Parse()
	showVersion(*isShowVersion)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := websocket.InitWebsocket(ctx, c.Websocket.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	if err := tcp.InitTCP(ctx, c.TCP.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	httpSrv := http.Start(":8000", ctx)
	defer func() {
		ctxTO, _ := context.WithTimeout(context.Background(), 5*time.Second)
		httpSrv.Shutdown(ctxTO)
	}()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		comet.RegisterCometServer(grpcServer, server.NewCometServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// showVersion 显示项目版本信息
func showVersion(isShow bool) {
	if isShow {
		fmt.Printf("Project: %s\n", projectName)
		fmt.Printf(" Version: %s\n", projectVersion)
		fmt.Printf(" Go Version: %s\n", goVersion)
		fmt.Printf(" Git Commit: %s\n", gitCommit)
		fmt.Printf(" Built: %s\n", buildTime)
		fmt.Printf(" OS/Arch: %s\n", osArch)
		os.Exit(0)
	}
}
