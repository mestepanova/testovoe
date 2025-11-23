package main

import (
	"context"
	"log/slog"

	"pr-manager-service/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	var err error
	defer func() {
		if p := recover(); p != nil {
			slog.Error("service failed: panic", slog.Any("panic", p))
		}

		if err != nil {
			slog.Error("service failed: error", slog.Any("error", err))
		}
	}()

	if err = godotenv.Load("../.env"); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.NewApp(ctx)
	if err != nil {
		return
	}
	defer application.Cleanup()

	if err = app.RunMigrations(application.Cfg.MigrateDSN()); err != nil {
		return
	}

	if err = application.RunHttpServer(ctx); err != nil {
		return
	}
}
