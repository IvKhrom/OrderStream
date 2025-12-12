package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
	"github.com/ivkhr/orderstream/internal/mocks"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateOrderAndPublish(t *testing.T) {
	// Prepare mock pool that accepts Exec
	mp := &mocks.MockPool{}
	called := false
	mp.ExecFunc = func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
		called = true
		var z pgconn.CommandTag
		return z, nil
	}

	p := &Postgres{pool: mp}

	o := &domain.Order{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
		Amount:  10,
		Payload: json.RawMessage(`{"id":"ext123"}`),
		Status:  "new",
		Bucket:  1,
	}

	if err := p.CreateOrder(context.Background(), o); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if !called {
		t.Fatalf("expected Exec to be called")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	mp := &mocks.MockPool{}
	// QueryRow returns row whose Scan returns error
	mp.QueryRowFunc = func(ctx context.Context, sql string, args ...interface{}) Row {
		return &mocks.MockRow{ScanFunc: func(dest ...interface{}) error { return errors.New("no rows") }}
	}
	p := &Postgres{pool: mp}
	_, err := p.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatalf("expected error from GetByID")
	}
}

func TestUpdateNoRows(t *testing.T) {
	mp := &mocks.MockPool{}
	mp.ExecFunc = func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
		var z pgconn.CommandTag
		return z, nil
	}
	p := &Postgres{pool: mp}
	o := &domain.Order{OrderID: uuid.New(), Payload: json.RawMessage(`{"foo":1}`)}
	if err := p.Update(context.Background(), o); err == nil {
		t.Fatalf("expected error when no rows updated")
	}
}

func TestUpdateStatusInvalidUUID(t *testing.T) {
	p := &Postgres{pool: &mocks.MockPool{}}
	if err := p.UpdateStatus(context.Background(), "not-uuid", "x"); err == nil {
		t.Fatalf("expected parse error for invalid uuid")
	}
}

func TestDeleteOrderInvalidUUID(t *testing.T) {
	p := &Postgres{pool: &mocks.MockPool{}}
	if err := p.DeleteOrder(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected parse error for invalid uuid on delete")
	}
}
