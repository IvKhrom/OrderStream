package ordersackconsumer

import (
	"context"
	"testing"

	"github.com/ivkhr/orderstream/internal/models"
)

type dummyAckProcessor struct{}

func (p *dummyAckProcessor) Handle(ctx context.Context, ack *models.OrderAck) error { return nil }

func TestNewOrdersAckConsumer(t *testing.T) {
	c := NewOrdersAckConsumer(&dummyAckProcessor{}, []string{"b1"}, "t1", "g1")
	if c.topicName != "t1" || c.groupID != "g1" {
		t.Fatalf("ожидали topic=t1 и group=g1, получили %q/%q", c.topicName, c.groupID)
	}
}


