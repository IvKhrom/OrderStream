package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Ручной мок для Pool/CommandTag/Row Postgres, используемый в тестах.
// Простой вариант мока, который позволяет запускать тесты без внешних генераторов.

type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func (r *MockRow) Scan(dest ...interface{}) error { return r.ScanFunc(dest...) }

type MockPool struct {
	ExecFunc     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	CloseFunc    func()
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, args...)
	}
	var z pgconn.CommandTag
	return z, nil
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return &MockRow{ScanFunc: func(dest ...interface{}) error { return nil }}
}

func (m *MockPool) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}
