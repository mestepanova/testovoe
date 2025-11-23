package storage

import (
	"context"

	"pr-manager-service/internal/usecases"

	"github.com/jackc/pgx/v4"
)

func (s *Storage) UnitOfWork(ctx context.Context, do func(txs usecases.Storage) error) error {
	return s.querier.BeginFunc(ctx, func(tx pgx.Tx) error {
		return do(NewStorage(tx))
	})
}
