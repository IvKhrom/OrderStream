package orders

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/shared/models"
)

type Storage interface {
	Create(ctx context.Context, o *models.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

// HandleOrderEvent — "рабочая" обработка события: обновление статусов в worker БД и формирование ACK.
func (s *Service) HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderAck, error) {
	// Если order_id пустой/"0" — создаём заказ и возвращаем ACK с новым id.
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

	if event.Status == "" {
		event.Status = "processing"
	}
	if err := s.storage.UpdateStatus(ctx, event.OrderID, event.Status); err != nil {
		return nil, err
	}
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


