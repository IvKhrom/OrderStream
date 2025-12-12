package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
)

// MockOrderRepository — простой ручной мок реализации OrderRepository.
// Используется как быстрая замена автоматически сгенерированных моков в тестах.
type MockOrderRepository struct {
	CreateFunc          func(ctx context.Context, o *domain.Order) error
	GetByIDFunc         func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByExternalIDFunc func(ctx context.Context, externalID string, userID uuid.UUID) (*domain.Order, error)
	UpdateFunc          func(ctx context.Context, o *domain.Order) error
	UpdateStatusFunc    func(ctx context.Context, id string, status string) error
	DeleteOrderFunc     func(ctx context.Context, id string) error
}

func (m *MockOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	if m == nil || m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, o)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if m == nil || m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockOrderRepository) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*domain.Order, error) {
	if m == nil || m.GetByExternalIDFunc == nil {
		return nil, nil
	}
	return m.GetByExternalIDFunc(ctx, externalID, userID)
}

func (m *MockOrderRepository) Update(ctx context.Context, o *domain.Order) error {
	if m == nil || m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, o)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if m == nil || m.UpdateStatusFunc == nil {
		return nil
	}
	return m.UpdateStatusFunc(ctx, id, status)
}

func (m *MockOrderRepository) DeleteOrder(ctx context.Context, id string) error {
	if m == nil || m.DeleteOrderFunc == nil {
		return nil
	}
	return m.DeleteOrderFunc(ctx, id)
}
