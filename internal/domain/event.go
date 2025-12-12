package domain

import (
	"encoding/json"
	"time"
)

type OrderEvent struct {
	// Поле EventID удалено: тип события определяется по `status` или по значению `order_id` в запросе.
	OrderID    string          `json:"order_id"`    // Внутренний идентификатор заказа (UUID)
	ExternalID string          `json:"external_id"` // Внешний идентификатор заказа (от маркетплейса)
	UserID     string          `json:"user_id"`     // Идентификатор пользователя (для проверки прав)
	Payload    json.RawMessage `json:"payload,omitempty"`
	// Возможные значения статуса:
	// - "new" — создан, требует обработки (worker должен выполнить создание);
	// - "processing" — в процессе обработки;
	// - "done" — обработан успешно;
	// - "cancelled" — заказ отменён (бизнес-статус, может требовать дополнительной очистки);
	// - "deleted" — мягкое удаление записи (soft-delete), административный статус.
	Status    string    `json:"status"` // new, processing, done, cancelled, deleted
	Timestamp time.Time `json:"timestamp"`
}

type OrderAck struct {
	// EventID удалён; ACK содержит внутренний `order_id` для сопоставления
	OrderID     string    `json:"order_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
