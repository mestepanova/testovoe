//go:build integration

package tests

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"pr-manager-service/internal/app"
	"pr-manager-service/internal/generated/api"
	"pr-manager-service/internal/storage"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB      *pgxpool.Pool
	client      *api.ClientWithResponses
	testStorage *storage.Storage
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env.test"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg := app.InitConfig()

	pgContainer := initPostgresContainer(ctx, cfg)
	slog.Debug("postgres container initiallized")

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to NewApp: %s", err)
	}

	testDB = application.PostgresConn
	testStorage = application.Storage

	serverErr := make(chan error, 1)

	go func() {
		if err := application.RunHttpServer(ctx); err != nil {
			log.Printf("failed to RunHttpServer: %s", err.Error())
			serverErr <- err
			return
		}
		serverErr <- nil
	}()


	// NOTE: ожидание запуска http-server
	addr := "http://" + cfg.Addr() + "/"
	deadline := time.Now().Add(5 * time.Second)

	for {
		select {
		case err := <-serverErr:
			if err != nil {
				log.Fatalf("http server failed to start: %v", err)
			}
			log.Fatalf("http server exited before it became ready")
		default:
		}

		if time.Now().After(deadline) {
			log.Fatalf("http server not ready after %s", 5*time.Second)
		}

		resp, err := http.Get(addr)
		if err == nil {
			resp.Body.Close()
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	slog.Debug("http server started")

	c, err := api.NewClientWithResponses("http://" + cfg.Addr())
	if err != nil {
		log.Fatalf("failed to NewClient: %s", err.Error())
	}
	client = c

	slog.Debug("ready to run tests")

	code := m.Run()

	cancel()

	err = <-serverErr
	if err != nil {
		log.Printf("http server stopped with error: %v", err)
		if code == 0 {
			code = 1
		}
	}

	termCtx, termCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := pgContainer.Terminate(termCtx); err != nil {
		log.Printf("failed to terminate postgres container: %v", err)
	}
	termCancel()

	slog.Debug("stopping tests", slog.Int("code", code))
}

func initPostgresContainer(ctx context.Context, cfg *app.Config) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image: "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DBUser,
			"POSTGRES_PASSWORD": cfg.DBPassword,
			"POSTGRES_DB":       cfg.DBName,
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("cannot get host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("cannot get mapped port: %v", err)
	}

	cfg.DBHost = host
	cfg.DBPort = mappedPort.Port()

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", mappedPort.Port())

	if err := app.RunMigrations(cfg.MigrateDSN()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return container
}

func cleanupDB(ctx context.Context, t *testing.T) {
	_, err := testDB.Exec(ctx, `
        truncate table users, teams, pull_requests, users_stats, pull_requests_stats
        restart identity cascade;
    `)
	if err != nil {
		t.Fatalf("truncate failed: %v", err)
	}
}
