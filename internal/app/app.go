package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"pr-manager-service/internal/generated/api"
	"pr-manager-service/internal/http_server"
	"pr-manager-service/internal/storage"
	"pr-manager-service/internal/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	Cfg          *Config
	PostgresConn *pgxpool.Pool
	Storage      *storage.Storage
	Usecases     *usecases.Usecases
	HttpServer   *http_server.HttpServer
	closeFuncs   []func()
}

func NewApp(ctx context.Context) (*App, error) {
	SetupLogger()

	closeFuncs := make([]func(), 0, 1)

	cfg := InitConfig()
	pgConn, err := InitPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("InitPostgres: %w", err)
	}
	closeFuncs = append(closeFuncs, pgConn.Close)

	storage := storage.NewStorage(pgConn)
	usecases := usecases.NewUsecases(storage)
	httpServer := http_server.NewHttpServer(usecases)

	return &App{
		Storage:      storage,
		Usecases:     usecases,
		HttpServer:   httpServer,
		Cfg:          cfg,
		PostgresConn: pgConn,
		closeFuncs:   closeFuncs,
	}, nil
}

func (a *App) Cleanup() {
	for _, fn := range a.closeFuncs {
		fn()
	}
}

func (a *App) RunHttpServer(ctx context.Context) error {
	router := gin.New()
	api.RegisterHandlers(router, a.HttpServer)

	srv := &http.Server{
		Addr:              a.Cfg.Addr(),
		Handler:           router,
		ReadHeaderTimeout: time.Millisecond * 500,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("http server shutdown error", slog.Any("error", err))
		}
	}()

	slog.Info(fmt.Sprintf("http server starting on %s", a.Cfg.Addr()), slog.Any("cfg", a.Cfg))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server ListenAndServe: %w", err)
	}

	slog.Info("http server stopped gracefully")
	return nil
}
