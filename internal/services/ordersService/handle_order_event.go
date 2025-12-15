package ordersService

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func (s *OrdersService) HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error) {
	// Обратная совместимость: если order_id пустой/"0" — создаём заказ и возвращаем ACK с новым id.
	if event.OrderID == "" || event.OrderID == "0" {
		oid := uuid.New()
		var userUUID uuid.UUID
		if parsed, err := uuid.Parse(event.UserID); err == nil {
			userUUID = parsed
		}
		bucket := models.BucketFromUUID(oid, 4)
		ord := &models.Order{
			OrderID:   oid,
			UserID:    userUUID,
			Amount:    0,
			Payload:   event.Payload,
			Status:    "processing",
			Bucket:    bucket,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.storage.Create(ctx, ord); err != nil {
			return nil, err
		}
		_ = s.storage.UpdateStatus(ctx, oid.String(), "done")
		return &models.OrderAck{
			OrderID:     oid.String(),
			Status:      "processed",
			ProcessedAt: time.Now().UTC(),
		}, nil
	}

	// Обычный поток: обновляем существующий заказ по данным события.
	if event.Status == "" {
		event.Status = "processing"
	}
	if err := s.storage.UpdateStatus(ctx, event.OrderID, event.Status); err != nil {
		return nil, err
	}
	// Для событий "new"/"processing" доводим до "done".
	if event.Status == "new" || event.Status == "processing" {
		_ = s.storage.UpdateStatus(ctx, event.OrderID, "done")
		return &models.OrderAck{
			OrderID:     event.OrderID,
			Status:      "processed",
			ProcessedAt: time.Now().UTC(),
		}, nil
	}

	return &models.OrderAck{
		OrderID:     event.OrderID,
		Status:      "updated",
		ProcessedAt: time.Now().UTC(),
	}, nil
}
