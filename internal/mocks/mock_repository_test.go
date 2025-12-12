package mocks_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
	"github.com/ivkhr/orderstream/internal/mocks"
)

func TestMockUpdateCalled(t *testing.T) {
	m := &mocks.MockOrderRepository{}
	called := false
	m.UpdateFunc = func(ctx context.Context, o *domain.Order) error {
		called = true
		o.Amount = 999
		return nil
	}

	order := &domain.Order{}
	if err := m.Update(context.Background(), order); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected UpdateFunc to be called")
	}
	if order.Amount != 999 {
		t.Fatalf("expected amount to be changed by mock; got %v", order.Amount)
	}

	// Проверяем поведение GetByID по умолчанию (возвращает nil)
	id := uuid.New()
	got, err := m.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error from GetByID: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil default for GetByID, got %+v", got)
	}
}
