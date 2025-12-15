package ordersService

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func TestOrdersService_HandleOrderEvent_ReturnsProcessedAck(t *testing.T) {
	st := newMemStorage()
	svc := NewOrdersService(st, nil, nil, 0)

	ev := &models.OrderEvent{
		OrderID:   uuid.New().String(),
		UserID:    uuid.New().String(),
		Payload:   json.RawMessage(`{"id":"ext-1"}`),
		Status:    "new",
		Timestamp: time.Now(),
	}

	ack, err := svc.HandleOrderEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ack.Status != "processed" {
		t.Fatalf("expected processed, got %q", ack.Status)
	}
}


