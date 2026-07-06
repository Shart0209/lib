package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgPooler interface {
	Master() Connector
	Sync() Connector
	Close()
}

type Connector interface {
	Executor

	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Executor) error) error
}

type Executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
