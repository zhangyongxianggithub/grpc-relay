package main

import (
	"log/slog"
	"os"

	_ "github.com/zhangyongxianggithub/grpc-relay/log"
)

func main() {
	_ = os.Setenv("ANSI_OUTPUT_ENABLED", "")
	slog.Error("colorful text", "aaa", "aaa")
}
