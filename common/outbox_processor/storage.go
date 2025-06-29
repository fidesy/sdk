package outbox_processor

import (
	"context"
	"fmt"
	"time"

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

func (s *storage) ListOutboxMessages(ctx context.Context, limit int64) ([]*Message, error) {
	query := postgres.Builder().
		Select("id, message, created_at, sent_at").
		From(s.tableName).
		Where(sq.Eq{
			"sent_at": nil,
		}).
		Limit(uint64(limit))

	return postgres.Select[Message](ctx, s.pool, query)
}

func (s *storage) UpdateOutboxMessagesSentAt(ctx context.Context, ids []int64) error {
	query := postgres.Builder().
		Update(s.tableName).
		SetMap(map[string]interface{}{
			"sent_at": time.Now().UTC(),
		}).
		Where(sq.Eq{
			"id": ids,
		}).
		Suffix("RETURNING id")

	_, err := postgres.Exec[Message](ctx, s.pool, query)
	return err
}
