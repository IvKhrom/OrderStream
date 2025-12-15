package bootstrap

import (
	ordersackprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_ack_processor"
	orderseventprocessor "github.com/ivkhr/orderstream/internal/services/processors/orders_event_processor"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
)

func InitOrdersEventProcessor(ordersService *ordersService.OrdersService, ackPublisher *kafkastorage.Publisher) *orderseventprocessor.OrdersEventProcessor {
	return orderseventprocessor.NewOrdersEventProcessor(ordersService, ackPublisher)
}

func InitOrdersAckProcessor(ackRegistry *ordersService.AckRegistry) *ordersackprocessor.OrdersAckProcessor {
	return ordersackprocessor.NewOrdersAckProcessor(ackRegistry)
}


