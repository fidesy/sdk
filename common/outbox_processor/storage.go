package outbox_processor

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy/sdk/common/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	tableName string
	pool      *pgxpool.Pool
}

func NewStorage(tableName string, pool *pgxpool.Pool) *storage {
	return &storage{
		tableName: fmt.Sprintf("%s_outbox", tableName),
		pool:      pool,
	}
}

func (s *storage) ListOutboxMessages(ctx context.Context, limit uint64) ([]*Message, error) {
	query := postgres.Builder().
		Select("id, message").
		From(s.tableName).
		Limit(limit)

	return postgres.Select[Message](ctx, s.pool, query)
}

func (s *storage) DeleteOutboxMessages(ctx context.Context, ids []int64) error {
	query := postgres.Builder().
		Delete(s.tableName).
		Where(sq.Eq{
			"id": ids,
		})

	_, err := postgres.Exec[Message](ctx, s.pool, query)
	return err
}
