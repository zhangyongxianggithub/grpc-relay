package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func init() {
	_, ok := os.LookupEnv("ANSI_OUTPUT_ENABLED")
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			TimeFormat: "2006-01-02T15:04:05.000000",
			AddSource:  true,
			NoColor:    !ok,
		}),
	))
}
