package main

import (
	"log/slog"

	_ "bestzyx.com/grpc-relay/log"
	"bestzyx.com/grpc-relay/relay/config"
	_ "bestzyx.com/grpc-relay/relay/internal/gateway"
)

func main() {
	slog.Info("relay for ic & data gateway is starting")
	config.InitializeConfig()
}
