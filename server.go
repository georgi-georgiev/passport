package passport

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPServer(lc fx.Lifecycle, conf *Config, log *zap.Logger, app *gin.Engine) *http.Server {
	maxHeaderBytes := 1 << 40

	server := &http.Server{
		Addr:           conf.Server.Host + ":" + conf.Server.Port,
		Handler:        app,
		ReadTimeout:    conf.Server.Timeout.Read * time.Second,
		WriteTimeout:   conf.Server.Timeout.Write * time.Second,
		IdleTimeout:    conf.Server.Timeout.Idle * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server at", zap.String("address", server.Addr))
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return server
}
