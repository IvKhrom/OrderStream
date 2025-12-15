package orders_service_api

import (
	"context"

	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

func (s *OrdersServiceAPI) Health(ctx context.Context, _ *orders_api.HealthRequest) (*orders_api.HealthResponse, error) {
	_ = ctx
	return &orders_api.HealthResponse{Status: "ok"}, nil
}


