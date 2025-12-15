package bootstrap

import (
	"github.com/ivkhr/orderstream/config"
	"github.com/ivkhr/orderstream/internal/services/ordersService"
	"github.com/ivkhr/orderstream/internal/storage/kafkastorage"
	"github.com/ivkhr/orderstream/internal/storage/pgstorage"
)

func InitOrdersService(storage *pgstorage.PGstorage, eventsPublisher *kafkastorage.Publisher, ackRegistry *ordersService.AckRegistry, cfg *config.Config) *ordersService.OrdersService {
	var eventsPub ordersService.OrdersEventsPublisher
	if eventsPublisher != nil {
		eventsPub = eventsPublisher
	}
	var ackWait ordersService.AckWaitRegistry
	if ackRegistry != nil {
		ackWait = ackRegistry
	}
	return ordersService.NewOrdersService(storage, eventsPub, ackWait, cfg.AckWaitTimeout)
}


