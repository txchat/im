package http

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/txchat/im/app/comet/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var svcCtx *svc.ServiceContext

func Start(addr string, svcCtx *svc.ServiceContext) *http.Server {
	svcCtx = svcCtx
	gin.ForceConsoleColor()
	engine := gin.Default()
	SetupEngine(engine)
	pprof.Register(engine)

	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Error("http listen failed", "err", err)
		}
	}()
	return srv
}

func SetupEngine(e *gin.Engine) *gin.Engine {
	statics := e.Group("/statics")
	statics.GET("/num", StaticsNumb)

	//buckets
	buckets := statics.Group("/buckets")
	buckets.GET("", BucketsInfo)

	users := statics.Group("/users")
	users.GET("/:uid", UserInfo)

	group := statics.Group("/group")
	group.GET("", GroupsInfo)
	group.GET("/:gid", GroupInfo)

	return e
}
