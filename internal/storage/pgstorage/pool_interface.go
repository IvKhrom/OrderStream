package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Pool — минимальный интерфейс для pgxpool.Pool, чтобы pgstorage можно было тестировать без реальной БД.
// Интерфейс используется только внутри storage слоя.
type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}
