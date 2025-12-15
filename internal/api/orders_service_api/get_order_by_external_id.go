package orders_service_api

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

func (s *OrdersServiceAPI) GetOrderByExternalID(ctx context.Context, req *orders_api.GetOrderByExternalIDRequest) (*orders_api.GetOrderByExternalIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &orders_api.GetOrderByExternalIDResponse{}, err
	}
	o, err := s.ordersService.GetOrderByExternalID(ctx, req.ExternalId, userID)
	if err != nil {
		return &orders_api.GetOrderByExternalIDResponse{}, err
	}
	return &orders_api.GetOrderByExternalIDResponse{Order: mapOrderToProto(o)}, nil
}
