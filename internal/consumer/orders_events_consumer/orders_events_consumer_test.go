package orderseventsconsumer

import (
	"context"
	"testing"

	"github.com/ivkhr/orderstream/internal/models"
)

type dummyProcessor struct{}

func (p *dummyProcessor) Handle(ctx context.Context, event *models.OrderEvent) error { return nil }

func TestNewOrdersEventsConsumer(t *testing.T) {
	c := NewOrdersEventsConsumer(&dummyProcessor{}, []string{"b1"}, "t1", "g1")
	if c.topicName != "t1" || c.groupID != "g1" {
		t.Fatalf("ожидали topic=t1 и group=g1, получили %q/%q", c.topicName, c.groupID)
	}
}


