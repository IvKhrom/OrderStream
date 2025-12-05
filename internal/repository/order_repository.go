package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	DeleteOrder(ctx context.Context, id string) error
}
