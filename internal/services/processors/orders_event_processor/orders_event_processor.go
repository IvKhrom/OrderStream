package orderseventprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/internal/models"
)

type ordersService interface {
	HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error)
}

// AckPublisher публикует ACK (подтверждение обработки события заказа) в Kafka.
// Вынесено как отдельный интерфейс, чтобы процессор не зависел от конкретной реализации/библиотеки.
type AckPublisher interface {
	PublishOrderAck(ctx context.Context, ack *models.OrderAck) error
}

type OrdersEventProcessor struct {
	ordersService ordersService
	ackPublisher  AckPublisher
}

func NewOrdersEventProcessor(ordersService ordersService, ackPublisher AckPublisher) *OrdersEventProcessor {
	return &OrdersEventProcessor{
		ordersService: ordersService,
		ackPublisher:  ackPublisher,
	}
}


