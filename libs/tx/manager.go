package tx

import (
	"context"
	"route256/libs/logger"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"route256/libs/tracer"
	"github.com/opentracing/opentracing-go"
)

var txKey = struct{}{}

type Manager struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Manager {
	return &Manager{pool: pool}
}

type DBProvider interface {
	GetDB(ctx context.Context) Querier
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

func (m *Manager) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "libs/tx")
	defer span.Finish()

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}
	defer func() {
		err = tx.Rollback(ctx)

		if err != nil {
			logger.Info("Rollback: ", err)
		}
	}()

	ctxTx := context.WithValue(ctx, txKey, tx)
	if err = fn(ctxTx); err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}

func (m *Manager) GetDB(ctx context.Context) Querier {
	tx, ok := ctx.Value(txKey).(Querier)
	if ok {
		return tx
	}

	return m.pool
}
