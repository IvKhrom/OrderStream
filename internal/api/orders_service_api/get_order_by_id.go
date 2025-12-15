package orders_service_api

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

func (s *OrdersServiceAPI) GetOrderByID(ctx context.Context, req *orders_api.GetOrderByIDRequest) (*orders_api.GetOrderByIDResponse, error) {
	oid, err := uuid.Parse(req.OrderId)
	if err != nil {
		return &orders_api.GetOrderByIDResponse{}, err
	}
	o, err := s.ordersService.GetOrderByID(ctx, oid)
	if err != nil {
		return &orders_api.GetOrderByIDResponse{}, err
	}
	return &orders_api.GetOrderByIDResponse{Order: mapOrderToProto(o)}, nil
}


