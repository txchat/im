package http

import (
	"net/http"

	"github.com/txchat/im/app/logic/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var serviceContext *svc.ServiceContext

func Start(addr string, svcCtx *svc.ServiceContext) *http.Server {
	serviceContext = svcCtx

	srv := &http.Server{
		Addr:    addr,
		Handler: nil,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Error("http listen failed", "err", err)
		}
	}()
	return srv
}
