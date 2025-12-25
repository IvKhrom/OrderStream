package bootstrap

import (
	orderseventprocessor "github.com/ivkhr/orderstream/services/worker/internal/services/processors/orders_event_processor"
)

func InitOrdersEventProcessor(
	svc orderseventprocessor.OrdersService,
	results orderseventprocessor.ResultWriter,
	ackPub orderseventprocessor.AckPublisher,
) *orderseventprocessor.Processor {
	return orderseventprocessor.New(svc, results, ackPub)
}


