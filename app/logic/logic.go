package logic

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/txchat/im/api/logic"
	"github.com/txchat/im/app/logic/internal/config"
	"github.com/txchat/im/app/logic/internal/http"
	"github.com/txchat/im/app/logic/internal/server"
	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// configFile 配置文件路径
var configFile = flag.String("f", "etc/logic.yaml", "the config file")

func Main() {
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)

	httpSrv := http.Start(":8001", ctx)
	defer func() {
		ctxTO, _ := context.WithTimeout(context.Background(), 5*time.Second)
		httpSrv.Shutdown(ctxTO)
	}()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		logic.RegisterLogicServer(grpcServer, server.NewLogicServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
