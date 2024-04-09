package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/fidesy/sdk/common/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/opentracing/opentracing-go"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, pgDsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, pgDsn)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func Exec[T any](ctx context.Context, db pgxscan.Querier, sqlizer sq.Sqlizer) (*T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres.Exec")
	defer span.Finish()

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	var model T

	err = pgxscan.Get(ctx, db, &model, query, args...)
	if err != nil {
		return nil, handleError(err)
	}

	return &model, nil
}

func Select[T any](ctx context.Context, db pgxscan.Querier, sqlizer sq.Sqlizer) ([]*T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgres.Select")
	defer span.Finish()

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	var model []*T
	err = pgxscan.Select(ctx, db, &model, query, args...)
	if err != nil {
		return nil, handleError(err)
	}

	return model, nil
}

func WithTransaction(ctx context.Context, pool *pgxpool.Pool, callback func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = callback(tx)
	if err != nil {
		rollBackErr := tx.Rollback(ctx)
		if rollBackErr != nil {
			logger.Errorf("tx.Rollback: %v", err)
		}

		return err
	}

	return tx.Commit(ctx)
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Check for PostgreSQL unique constraint violation error
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}

	return err
}
