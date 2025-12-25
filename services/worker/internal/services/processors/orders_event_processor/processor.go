package orderseventprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

type OrdersService interface {
	HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderResult, error)
}

type AckPublisher interface {
	PublishOrderAck(ctx context.Context, ack *models.OrderAck) error
}

type ResultWriter interface {
	SetOrderResult(ctx context.Context, res *models.OrderResult) error
}

type Processor struct {
	svc     OrdersService
	results ResultWriter
	ackPub  AckPublisher
}

func New(svc OrdersService, results ResultWriter, ackPub AckPublisher) *Processor {
	return &Processor{svc: svc, results: results, ackPub: ackPub}
}
