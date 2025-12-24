package ordersService

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

type OrdersStorage interface {
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
	Update(ctx context.Context, o *models.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	DeleteOrder(ctx context.Context, id string) error
}

type OrdersEventsPublisher interface {
	PublishOrderEvent(ctx context.Context, event *models.OrderEvent) error
}

type AckWaitRegistry interface {
	Register(orderID string) (ch <-chan struct{}, cleanup func())
}

// AckRegistryContract объединяет ожидание ACK (Register) и нотификацию (Notify)
type AckRegistryContract interface {
	AckWaitRegistry
	AckNotifier
}

type OrdersService struct {
	storage        OrdersStorage
	eventsPub      OrdersEventsPublisher
	ackCoordinator AckCoordinator
}

func NewOrdersService(storage OrdersStorage, eventsPub OrdersEventsPublisher, ackRegistry AckWaitRegistry, ackWaitTimeout time.Duration) *OrdersService {
	return &OrdersService{
		storage:        storage,
		eventsPub:      eventsPub,
		ackCoordinator: NewAckCoordinator(ackRegistry, ackWaitTimeout),
	}
}
