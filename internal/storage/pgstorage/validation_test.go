package pgstorage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestUpdateStatus_InvalidUUID(t *testing.T) {
	p := &PGstorage{pool: &mockPool{}}
	if err := p.UpdateStatus(context.Background(), "not-uuid", "x"); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestUpdateStatus_Ok(t *testing.T) {
	mp := &mockPool{
		exec: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag("UPDATE 1"), nil
		},
	}
	p := &PGstorage{pool: mp}
	if err := p.UpdateStatus(context.Background(), uuid.New().String(), "done"); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}

func TestDeleteOrder_InvalidUUID(t *testing.T) {
	p := &PGstorage{pool: &mockPool{}}
	if err := p.DeleteOrder(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestDeleteOrder_Ok(t *testing.T) {
	mp := &mockPool{
		exec: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag("UPDATE 1"), nil
		},
	}
	p := &PGstorage{pool: mp}
	if err := p.DeleteOrder(context.Background(), uuid.New().String()); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}


