package orderseventprocessor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ivkhr/orderstream/internal/models"
)

type fakeOrdersService struct{}

func (s *fakeOrdersService) HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error) {
	_ = ctx
	return &models.OrderAck{OrderID: event.OrderID, Status: "processed"}, nil
}

type capturePublisher struct {
	last []byte
}

func (p *capturePublisher) Publish(ctx context.Context, value []byte) error {
	_ = ctx
	p.last = append([]byte(nil), value...)
	return nil
}

func TestOrdersEventProcessor_Handle_PublishesAck(t *testing.T) {
	pub := &capturePublisher{}
	processor := NewOrdersEventProcessor(&fakeOrdersService{}, pub)

	ev := &models.OrderEvent{OrderID: "id-1"}
	if err := processor.Handle(context.Background(), ev); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	var ack models.OrderAck
	if err := json.Unmarshal(pub.last, &ack); err != nil {
		t.Fatalf("unmarshal ack: %v", err)
	}
	if ack.OrderID != "id-1" {
		t.Fatalf("expected ack for id-1, got %q", ack.OrderID)
	}
}


