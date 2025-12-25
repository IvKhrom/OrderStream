package orderseventprocessor

import (
	"context"

	"github.com/ivkhr/orderstream/shared/models"
)

type OrdersService interface {
	HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error)
}

type AckPublisher interface {
	PublishOrderAck(ctx context.Context, ack *models.OrderAck) error
}

type ResultWriter interface {
	SetOrderAck(ctx context.Context, ack *models.OrderAck) error
}

type Processor struct {
	svc      OrdersService
	results  ResultWriter
	ackPub   AckPublisher
}

func New(svc OrdersService, results ResultWriter, ackPub AckPublisher) *Processor {
	return &Processor{svc: svc, results: results, ackPub: ackPub}
}

func (p *Processor) Handle(ctx context.Context, event *models.OrderEvent) error {
	ack, err := p.svc.HandleOrderEvent(ctx, event)
	if err != nil {
		return err
	}
	if err := p.results.SetOrderAck(ctx, ack); err != nil {
		return err
	}
	return p.ackPub.PublishOrderAck(ctx, ack)
}


