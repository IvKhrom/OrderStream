package ordersService

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func (s *OrdersService) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.storage.GetByID(ctx, id)
}


