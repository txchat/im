package comet

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/txchat/im/api/comet"
	"github.com/txchat/im/app/comet/internal/config"
	"github.com/txchat/im/app/comet/internal/http"
	xnet "github.com/txchat/im/app/comet/internal/net"
	"github.com/txchat/im/app/comet/internal/server"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// configFile 配置文件路径
var configFile = flag.String("f", "etc/comet.yaml", "the config file")

func Main() {
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)

	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := xnet.InitServer(ctx, c.Websocket.Bind, runtime.NumCPU(), xnet.WebsocketServer); err != nil {
		panic(err)
	}
	if err := xnet.InitServer(ctx, c.TCP.Bind, runtime.NumCPU(), xnet.TCPServer); err != nil {
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
