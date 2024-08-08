package main

import (
	"log/slog"
	"net/http"

	"github.com/zhangyongxianggithub/grpc-relay/gateway/config"
	_ "github.com/zhangyongxianggithub/grpc-relay/gateway/internal/charge"
	_ "github.com/zhangyongxianggithub/grpc-relay/gateway/internal/eb"
	_ "github.com/zhangyongxianggithub/grpc-relay/gateway/internal/ic"
	"github.com/zhangyongxianggithub/grpc-relay/gateway/internal/mux"
	_ "github.com/zhangyongxianggithub/grpc-relay/log"
	"golang.org/x/xerrors"
)

func main() {
	slog.Info("gateway for ic & data gateway is starting")
	multiplexer := mux.NewChainMux()
	multiplexer.Use(mux.RecoveryInterceptor)
	config.InitializeConfig(multiplexer)
	if err := http.ListenAndServe(config.Server.Server.Listen, multiplexer); err != nil {
		xerr := xerrors.New(err.Error())
		slog.Error("gateway for ic & data gateway started failed",
			"addr", config.Server.Server.Listen, slog.Any("err", xerr))
		panic(err)
	}
}
