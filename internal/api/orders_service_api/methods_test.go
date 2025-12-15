package orders_service_api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

func TestOrdersServiceAPI_Health(t *testing.T) {
	api := NewOrdersServiceAPI(&stubOrdersService{})
	resp, err := api.Health(context.Background(), &orders_api.HealthRequest{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("ожидали ok, получили %q", resp.Status)
	}
}

func TestOrdersServiceAPI_GetOrderByID_InvalidUUID(t *testing.T) {
	api := NewOrdersServiceAPI(&stubOrdersService{})
	_, err := api.GetOrderByID(context.Background(), &orders_api.GetOrderByIDRequest{OrderId: "bad"})
	if err == nil {
		t.Fatalf("ожидали ошибку uuid.Parse")
	}
}

func TestOrdersServiceAPI_GetOrderByID_Ok(t *testing.T) {
	api := NewOrdersServiceAPI(&stubOrdersService{})
	id := uuid.New().String()
	resp, err := api.GetOrderByID(context.Background(), &orders_api.GetOrderByIDRequest{OrderId: id})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if resp.Order == nil || resp.Order.OrderId != id {
		t.Fatalf("ожидали order с id=%s", id)
	}
}

func TestOrdersServiceAPI_GetOrderByExternalID_InvalidUserID(t *testing.T) {
	api := NewOrdersServiceAPI(&stubOrdersService{})
	_, err := api.GetOrderByExternalID(context.Background(), &orders_api.GetOrderByExternalIDRequest{ExternalId: "ext", UserId: "bad"})
	if err == nil {
		t.Fatalf("ожидали ошибку uuid.Parse")
	}
}

func TestOrdersServiceAPI_UpsertOrder_InvalidUserID(t *testing.T) {
	api := NewOrdersServiceAPI(&stubOrdersService{})
	_, err := api.UpsertOrder(context.Background(), &orders_api.UpsertOrderRequest{OrderId: "0", UserId: "bad", PayloadJson: "{}"})
	if err == nil {
		t.Fatalf("ожидали ошибку uuid.Parse")
	}
}


