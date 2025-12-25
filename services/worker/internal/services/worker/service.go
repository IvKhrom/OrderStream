package worker

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

type Storage interface {
	Upsert(ctx context.Context, o *models.Order) (created bool, err error)
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, o *models.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	DeleteOrder(ctx context.Context, id string) error
}
