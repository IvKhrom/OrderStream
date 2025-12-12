package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockery --name Pool --output ../../mocks --outpkg mocks --case underscore

// CommandTag is an alias to pgconn.CommandTag so signatures match pgxpool.
type CommandTag = pgconn.CommandTag

// Row is an alias to pgx.Row to match QueryRow return type.
type Row = pgx.Row

// Pool is an abstraction over *pgxpool.Pool used by Postgres adapter.
type Pool interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

// Ensure that *pgxpool.Pool implements Pool
var _ Pool = (*pgxpool.Pool)(nil)
