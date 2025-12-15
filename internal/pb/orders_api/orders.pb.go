package orders_api

import (
	"context"

	proto_models "github.com/ivkhr/orderstream/internal/pb/models"
)

type HealthRequest struct{}
type HealthResponse struct {
	Status string `json:"status,omitempty"`
}

type UpsertOrderRequest struct {
	OrderId     string `json:"order_id,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	Status      string `json:"status,omitempty"`
	PayloadJson string `json:"payload_json,omitempty"`
}
type UpsertOrderResponse struct {
	OrderId string `json:"order_id,omitempty"`
	Status  string `json:"status,omitempty"`
}

type GetOrderByIDRequest struct {
	OrderId string `json:"order_id,omitempty"`
}
type GetOrderByIDResponse struct {
	Order *proto_models.OrderModel `json:"order,omitempty"`
}

type GetOrderByExternalIDRequest struct {
	ExternalId string `json:"external_id,omitempty"`
	UserId     string `json:"user_id,omitempty"`
}
type GetOrderByExternalIDResponse struct {
	Order *proto_models.OrderModel `json:"order,omitempty"`
}

type OrdersServiceServer interface {
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	UpsertOrder(context.Context, *UpsertOrderRequest) (*UpsertOrderResponse, error)
	GetOrderByID(context.Context, *GetOrderByIDRequest) (*GetOrderByIDResponse, error)
	GetOrderByExternalID(context.Context, *GetOrderByExternalIDRequest) (*GetOrderByExternalIDResponse, error)
}

type UnimplementedOrdersServiceServer struct{}
