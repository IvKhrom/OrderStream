package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

func (s *Service) HandleOrderEvent(ctx context.Context, event *models.OrderEvent) (*models.OrderResult, error) {
	if event == nil {
		return nil, errors.New("nil event")
	}
	userUUID, err := uuid.Parse(event.UserID)
	if err != nil {
		return nil, err
	}

	extFromPayload := payloadID(event.Payload)
	if event.ExternalID != "" && extFromPayload != "" && event.ExternalID != extFromPayload {
		slog.Warn("worker.handle_order_event.external_id_mismatch",
			"order_id", event.OrderID,
			"event_external_id", event.ExternalID,
			"payload_external_id", extFromPayload,
		)
	}

	slog.Info("worker.handle_order_event.received",
		"order_id", event.OrderID,
		"user_id", event.UserID,
		"status", event.Status,
		"external_id", func() string {
			if event.ExternalID != "" {
				return event.ExternalID
			}
			return extFromPayload
		}(),
	)

	// Обратная совместимость как в @internal: если order_id пустой/"0" — создаём новый.
	if event.OrderID == "" || event.OrderID == "0" {
		oid := uuid.New()
		now := time.Now()
		ord := &models.Order{
			OrderID:   oid,
			UserID:    userUUID,
			Amount:    extractAmount(event.Payload),
			Payload:   event.Payload,
			Status:    "done",
			Bucket:    models.BucketFromUUID(oid, 4),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.storage.Create(ctx, ord); err != nil {
			return nil, err
		}
		return &models.OrderResult{
			OrderID:     ord.OrderID.String(),
			Status:      "created",
			OrderStatus: ord.Status,
			ProcessedAt: time.Now().UTC(),
			UserID:      ord.UserID.String(),
			Amount:      ord.Amount,
			Payload:     ord.Payload,
			CreatedAt:   ord.CreatedAt,
			UpdatedAt:   ord.UpdatedAt,
			Bucket:      ord.Bucket,
		}, nil
	}

	oid, err := uuid.Parse(event.OrderID)
	if err != nil {
		return nil, err
	}

	existing, err := s.storage.GetByID(ctx, oid)
	if err != nil {
		now := time.Now()
		ord := &models.Order{
			OrderID:   oid,
			UserID:    userUUID,
			Amount:    extractAmount(event.Payload),
			Payload:   event.Payload,
			Status:    "done",
			Bucket:    models.BucketFromUUID(oid, 4),
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, upErr := s.storage.Upsert(ctx, ord)
		if upErr != nil {
			return nil, upErr
		}
		return &models.OrderResult{
			OrderID:     event.OrderID,
			Status:      "created",
			OrderStatus: ord.Status,
			ProcessedAt: time.Now().UTC(),
			UserID:      ord.UserID.String(),
			Amount:      ord.Amount,
			Payload:     ord.Payload,
			CreatedAt:   ord.CreatedAt,
			UpdatedAt:   ord.UpdatedAt,
			Bucket:      ord.Bucket,
		}, nil
	}

	// Проверка "можно ли апдейтить": если уже deleted — ничего не меняем.
	if existing.Status == "deleted" && event.Status != "deleted" {
		return &models.OrderResult{
			OrderID:     existing.OrderID.String(),
			Status:      "deleted",
			OrderStatus: "deleted",
			ProcessedAt: time.Now().UTC(),
			UserID:      existing.UserID.String(),
			Amount:      existing.Amount,
			Payload:     existing.Payload,
			CreatedAt:   existing.CreatedAt,
			UpdatedAt:   existing.UpdatedAt,
			Bucket:      existing.Bucket,
		}, nil
	}

	if event.Status == "deleted" {
		if err := s.storage.DeleteOrder(ctx, existing.OrderID.String()); err != nil {
			return nil, err
		}
		return &models.OrderResult{
			OrderID:     existing.OrderID.String(),
			Status:      "deleted",
			OrderStatus: "deleted",
			ProcessedAt: time.Now().UTC(),
			UserID:      existing.UserID.String(),
			Amount:      existing.Amount,
			Payload:     existing.Payload,
			CreatedAt:   existing.CreatedAt,
			UpdatedAt:   time.Now(),
			Bucket:      existing.Bucket,
		}, nil
	}

	if len(event.Payload) != 0 {
		existing.Payload = event.Payload
		existing.Amount = extractAmount(event.Payload)
	}
	existing.UpdatedAt = time.Now()
	if err := s.storage.Update(ctx, existing); err != nil {
		return nil, err
	}
	_ = s.storage.UpdateStatus(ctx, existing.OrderID.String(), "done")

	slog.Info("worker.handle_order_event.done",
		"order_id", event.OrderID,
		"op", "updated",
		"amount", existing.Amount,
	)

	return &models.OrderResult{
		OrderID:     existing.OrderID.String(),
		Status:      "updated",
		OrderStatus: "done",
		ProcessedAt: time.Now().UTC(),
		UserID:      existing.UserID.String(),
		Amount:      existing.Amount,
		Payload:     existing.Payload,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
		Bucket:      existing.Bucket,
	}, nil
}
