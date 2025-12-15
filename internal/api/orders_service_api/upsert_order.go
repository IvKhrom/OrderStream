package orders_service_api

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

func (s *OrdersServiceAPI) UpsertOrder(ctx context.Context, req *orders_api.UpsertOrderRequest) (*orders_api.UpsertOrderResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &orders_api.UpsertOrderResponse{}, err
	}
	payload := json.RawMessage(req.PayloadJson)

	orderID, status, err := s.ordersService.UpsertOrder(ctx, req.OrderId, userID, req.Status, payload)
	if err != nil {
		return &orders_api.UpsertOrderResponse{OrderId: orderID, Status: status}, err
	}
	return &orders_api.UpsertOrderResponse{OrderId: orderID, Status: status}, nil
}
