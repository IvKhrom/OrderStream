package bootstrap

import (
	"github.com/ivkhr/orderstream/internal/services/ordersService"
	ordersackprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_ack_processor"
	orderseventprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_event_processor"
)

func InitOrdersEventProcessor(ordersService *ordersService.OrdersService, ackProducer orderseventprocessor.AckPublisher) *orderseventprocessor.OrdersEventProcessor {
	return orderseventprocessor.NewOrdersEventProcessor(ordersService, ackProducer)
}

func InitOrdersAckProcessor(notifier ordersService.AckNotifier) *ordersackprocessor.OrdersAckProcessor {
	return ordersackprocessor.NewOrdersAckProcessor(notifier)
}
