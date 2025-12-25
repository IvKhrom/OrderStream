package api_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/services/api_service/internal/models"
)

type Storage interface {
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, o *models.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	DeleteOrder(ctx context.Context, id string) error
}

type EventsPublisher interface {
	PublishOrderEvent(ctx context.Context, event *models.OrderEvent) error
}

type ResultsStore interface {
	GetOrderAck(ctx context.Context, orderID string) (*models.OrderAck, bool, error)
	GetOrderResult(ctx context.Context, orderID string) (*models.OrderResult, bool, error)
}
