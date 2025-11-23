package app

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	log := slog.New(handler)
	slog.SetDefault(log)
}
