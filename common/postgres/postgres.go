package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy/sdk/common/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
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
	var span opentracing.Span

	skipSpan, ok := ctx.Value("skip_span").(bool)
	if !ok || !skipSpan {
		span, ctx = opentracing.StartSpanFromContext(ctx, "postgres.Exec")
		defer span.Finish()
	}

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

func ExecWithOutbox(ctx context.Context, db *pgxpool.Pool, dst Model, sqlizer sq.Sqlizer) error {
	var span opentracing.Span

	skipSpan, ok := ctx.Value("skip_span").(bool)
	if !ok || !skipSpan {
		span, ctx = opentracing.StartSpanFromContext(ctx, "postgres.Exec")
		defer span.Finish()
	}

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	err = WithTransaction(ctx, db, func(tx pgx.Tx) error {
		err = pgxscan.Get(ctx, tx, dst, query, args...)
		if err != nil {
			return handleError(err)
		}

		message, err := json.Marshal(dst)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}

		outboxQuery := Builder().
			Insert(fmt.Sprintf("%s_outbox", dst.TableName())).
			SetMap(map[string]interface{}{
				"message": message,
			})

		outboxSql, outboxArgs, err := outboxQuery.ToSql()
		if err != nil {
			return fmt.Errorf("outboxQuery.ToSql: %w", err)
		}

		_, err = tx.Exec(ctx, outboxSql, outboxArgs...)
		if err != nil {
			return handleError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("WithTransaction: %w", err)
	}

	return nil
}

func ExecWithOutboxTx(ctx context.Context, tx pgx.Tx, dst Model, sqlizer sq.Sqlizer) error {
	var span opentracing.Span

	skipSpan, ok := ctx.Value("skip_span").(bool)
	if !ok || !skipSpan {
		span, ctx = opentracing.StartSpanFromContext(ctx, "postgres.ExecWithOutboxTx")
		defer span.Finish()
	}

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("sqlizer.ToSql: %v", err)
	}

	err = pgxscan.Get(ctx, tx, dst, query, args...)
	if err != nil {
		return handleError(err)
	}

	message, err := json.Marshal(dst)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	outboxQuery := Builder().
		Insert(fmt.Sprintf("%s_outbox", dst.TableName())).
		SetMap(map[string]interface{}{
			"message": message,
		})

	outboxSql, outboxArgs, err := outboxQuery.ToSql()
	if err != nil {
		return fmt.Errorf("outboxQuery.ToSql: %w", err)
	}

	_, err = tx.Exec(ctx, outboxSql, outboxArgs...)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func Select[T any](ctx context.Context, db pgxscan.Querier, sqlizer sq.Sqlizer) ([]*T, error) {
	var span opentracing.Span

	skipSpan, ok := ctx.Value("skip_span").(bool)
	if !ok || !skipSpan {
		span, ctx = opentracing.StartSpanFromContext(ctx, "postgres.Select")
		defer span.Finish()
	}

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
