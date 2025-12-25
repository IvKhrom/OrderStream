package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGStorge(dsn string) (*PGstorage, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PGstorage{pool: pool}, nil
}


