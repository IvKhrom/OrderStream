package orders_service_api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/models"
)

type stubOrdersService struct{}

func (s *stubOrdersService) UpsertOrder(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error) {
	_ = ctx
	_ = payload
	return orderID, status, nil
}
func (s *stubOrdersService) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	_ = ctx
	return &models.Order{
		OrderID:   id,
		UserID:    uuid.New(),
		Amount:    1,
		Status:    "new",
		Payload:   json.RawMessage(`{"id":"ext"}`),
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
		Bucket:    1,
	}, nil
}
func (s *stubOrdersService) GetOrderByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	_ = ctx
	_ = externalID
	return &models.Order{
		OrderID:   uuid.New(),
		UserID:    userID,
		Amount:    2,
		Status:    "new",
		Payload:   json.RawMessage(`{}`),
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
		Bucket:    2,
	}, nil
}

func TestMapOrderToProto(t *testing.T) {
	o := &models.Order{
		OrderID:   uuid.New(),
		UserID:    uuid.New(),
		Amount:    3.5,
		Status:    "new",
		Payload:   json.RawMessage(`{"x":1}`),
		CreatedAt: time.Unix(1, 0),
		UpdatedAt: time.Unix(2, 0),
		Bucket:    3,
	}
	p := mapOrderToProto(o)
	if p == nil || p.OrderId == "" {
		t.Fatalf("ожидали заполненную proto-модель")
	}
}


