package pgstorage

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockRow struct {
	scan func(dest ...any) error
}

func (r *mockRow) Scan(dest ...any) error { return r.scan(dest...) }

type mockPool struct {
	exec    func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	queryRow func(ctx context.Context, sql string, args ...any) pgx.Row
	closed  bool
}

func (p *mockPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if p.exec == nil {
		return pgconn.CommandTag{}, nil
	}
	return p.exec(ctx, sql, args...)
}
func (p *mockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if p.queryRow == nil {
		return &mockRow{scan: func(dest ...any) error { return errors.New("no row mock") }}
	}
	return p.queryRow(ctx, sql, args...)
}
func (p *mockPool) Close() { p.closed = true }

func TestPGstorage_Close(t *testing.T) {
	mp := &mockPool{}
	st := &PGstorage{pool: mp}
	st.Close(context.Background())
	if !mp.closed {
		t.Fatalf("ожидали Close() пула")
	}
}

func TestPGstorage_Create(t *testing.T) {
	mp := &mockPool{
		exec: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			if sql == "" || len(args) == 0 {
				t.Fatalf("ожидали sql и args")
			}
			return pgconn.CommandTag{}, nil
		},
	}
	st := &PGstorage{pool: mp}
	o := &models.Order{
		OrderID:   uuid.New(),
		UserID:    uuid.New(),
		Amount:    10,
		Status:    "new",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		Bucket:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := st.Create(context.Background(), o); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}

func TestPGstorage_GetByID(t *testing.T) {
	wantID := uuid.New()
	wantUser := uuid.New()
	payload := []byte(`{"id":"x"}`)

	mp := &mockPool{
		queryRow: func(ctx context.Context, sql string, args ...any) pgx.Row {
			return &mockRow{scan: func(dest ...any) error {
				*(dest[0].(*uuid.UUID)) = wantID
				*(dest[1].(*uuid.UUID)) = wantUser
				*dest[2].(*float64) = 7
				*dest[3].(*string) = "new"
				*dest[4].(*[]byte) = append([]byte(nil), payload...)
				*dest[5].(*time.Time) = time.Unix(1, 0)
				*dest[6].(*time.Time) = time.Unix(2, 0)
				*dest[7].(*int) = 2
				return nil
			}}
		},
	}
	st := &PGstorage{pool: mp}

	got, err := st.GetByID(context.Background(), wantID)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got.OrderID != wantID {
		t.Fatalf("ожидали OrderID %v, получили %v", wantID, got.OrderID)
	}
	if string(got.Payload) != string(payload) {
		t.Fatalf("ожидали payload %q, получили %q", string(payload), string(got.Payload))
	}
}

func TestPGstorage_Update_NoRows(t *testing.T) {
	mp := &mockPool{
		exec: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			// CommandTag{} => RowsAffected() == 0
			return pgconn.CommandTag{}, nil
		},
	}
	st := &PGstorage{pool: mp}
	o := &models.Order{OrderID: uuid.New(), Payload: json.RawMessage(`{"x":1}`)}
	if err := st.Update(context.Background(), o); err == nil {
		t.Fatalf("ожидали ошибку, когда RowsAffected == 0")
	}
}

func TestPGstorage_Update_Ok(t *testing.T) {
	mp := &mockPool{
		exec: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag("UPDATE 1"), nil
		},
	}
	st := &PGstorage{pool: mp}
	o := &models.Order{OrderID: uuid.New(), Payload: json.RawMessage(`{"x":1}`)}
	if err := st.Update(context.Background(), o); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}

func TestPGstorage_GetByExternalID(t *testing.T) {
	wantID := uuid.New()
	wantUser := uuid.New()
	payload := []byte(`{"id":"x"}`)

	mp := &mockPool{
		queryRow: func(ctx context.Context, sql string, args ...any) pgx.Row {
			return &mockRow{scan: func(dest ...any) error {
				*(dest[0].(*uuid.UUID)) = wantID
				*(dest[1].(*uuid.UUID)) = wantUser
				*dest[2].(*float64) = 7
				*dest[3].(*string) = "new"
				*dest[4].(*[]byte) = append([]byte(nil), payload...)
				*dest[5].(*time.Time) = time.Unix(1, 0)
				*dest[6].(*time.Time) = time.Unix(2, 0)
				*dest[7].(*int) = 2
				return nil
			}}
		},
	}
	st := &PGstorage{pool: mp}

	got, err := st.GetByExternalID(context.Background(), "ext", wantUser)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got.OrderID != wantID {
		t.Fatalf("ожидали OrderID %v, получили %v", wantID, got.OrderID)
	}
}


