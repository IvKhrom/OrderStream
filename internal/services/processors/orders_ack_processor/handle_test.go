package ordersackprocessor

import (
	"context"
	"testing"

	"github.com/ivkhr/orderstream/internal/models"
)

type captureNotifier struct {
	last string
}

func (n *captureNotifier) Notify(orderID string) {
	n.last = orderID
}

func TestOrdersAckProcessor_Handle_Notifies(t *testing.T) {
	n := &captureNotifier{}
	p := NewOrdersAckProcessor(n)

	if err := p.Handle(context.Background(), &models.OrderAck{OrderID: "id-1"}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if n.last != "id-1" {
		t.Fatalf("expected notify id-1, got %q", n.last)
	}
}


