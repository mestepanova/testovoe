package app

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	pgConnRetryCount     = 5
	delayBetweenAttempts = time.Second
)

func InitPostgres(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ParseConfig: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool

	for range pgConnRetryCount {
		var attemptPool *pgxpool.Pool
		attemptPool, err = pgxpool.ConnectConfig(ctx, config)
		if err != nil {
			time.Sleep(delayBetweenAttempts)
			continue
		}

		if err = attemptPool.Ping(ctx); err != nil {
			attemptPool.Close()
			time.Sleep(delayBetweenAttempts)
			continue
		}

		pool = attemptPool
		break
	}

	if pool == nil {
		return nil, fmt.Errorf("failed to init pg pool: %w", err)
	}

	return pool, nil
}

func RunMigrations(dsn string) error {
	m, err := migrate.New("file://../migrations", dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}
