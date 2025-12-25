package api_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/services/api_service/internal/models"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *Service) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	return s.storage.GetByExternalID(ctx, externalID, userID)
}

// GetViewByID возвращает заказ, подмешивая обработанный результат из Redis (если он уже готов).
func (s *Service) GetViewByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	o, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if res, ok, _ := s.results.GetOrderResult(ctx, id.String()); ok && res != nil {
		o.Amount = res.Amount
		if res.OrderStatus != "" {
			o.Status = res.OrderStatus
		} else {
			// fallback для старых записей, где status использовался как order status
			o.Status = res.Status
		}
		if len(res.Payload) != 0 {
			o.Payload = res.Payload
		}
		if !res.UpdatedAt.IsZero() {
			o.UpdatedAt = res.UpdatedAt
		}
		if !res.CreatedAt.IsZero() {
			o.CreatedAt = res.CreatedAt
		}
		if res.Bucket != 0 {
			o.Bucket = res.Bucket
		}
	}
	return o, nil
}

// GetViewByExternalID возвращает заказ по external_id, подмешивая обработанный результат из Redis (если он готов).
func (s *Service) GetViewByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	o, err := s.storage.GetByExternalID(ctx, externalID, userID)
	if err != nil {
		return nil, err
	}
	if res, ok, _ := s.results.GetOrderResult(ctx, o.OrderID.String()); ok && res != nil {
		o.Amount = res.Amount
		if res.OrderStatus != "" {
			o.Status = res.OrderStatus
		} else {
			o.Status = res.Status
		}
		if len(res.Payload) != 0 {
			o.Payload = res.Payload
		}
		if !res.UpdatedAt.IsZero() {
			o.UpdatedAt = res.UpdatedAt
		}
		if !res.CreatedAt.IsZero() {
			o.CreatedAt = res.CreatedAt
		}
		if res.Bucket != 0 {
			o.Bucket = res.Bucket
		}
	}
	return o, nil
}


