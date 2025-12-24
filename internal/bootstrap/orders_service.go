package bootstrap

import (
	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
)

func InitOrdersService(
	storage ordersService.OrdersStorage,
	eventsPublisher ordersService.OrdersEventsPublisher,
	ackRegistry ordersService.AckWaitRegistry,
	cfg *config.Config,
) *ordersService.OrdersService {
	return ordersService.NewOrdersService(storage, eventsPublisher, ackRegistry, cfg.AckWaitTimeout)
}


