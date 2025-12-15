package ordersService

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func (s *OrdersService) GetOrderByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	return s.storage.GetByExternalID(ctx, externalID, userID)
}


