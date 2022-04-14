package http

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/txchat/im/comet"
	"github.com/txchat/im/comet/conf"
	"net/http"
)

var (
	srv *comet.Comet
)

func Start(addr string, s *comet.Comet) *http.Server {
	srv = s
	gin.ForceConsoleColor()
	switch conf.Conf.Env {
	case conf.DebugMode:
		gin.SetMode(gin.DebugMode)
	case conf.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	}
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
			log.Fatal().Err(err).Msg("http listen failed")
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
