package main

import (
	"log/slog"

	_ "github.com/zhangyongxianggithub/grpc-relay/log"
	"github.com/zhangyongxianggithub/grpc-relay/relay/config"
	_ "github.com/zhangyongxianggithub/grpc-relay/relay/internal/gateway"
)

func main() {
	slog.Info("relay for ic & data gateway is starting")
	config.InitializeConfig()
}
