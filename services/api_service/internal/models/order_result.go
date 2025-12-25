package models

import (
	"encoding/json"
	"time"
)

// OrderResult — то, что worker кладёт в Redis по ключу order_result:<order_id>.
// Структура намеренно совместима с worker/internal/models/OrderResult.
type OrderResult struct {
	OrderID     string          `json:"order_id"`
	// Status — результат операции (created|updated|deleted), приходит из worker через Redis/ACK.
	Status      string          `json:"status"`
	ProcessedAt time.Time       `json:"processed_at"`

	// OrderStatus — фактический статус заказа (new|processing|done|deleted).
	OrderStatus string          `json:"order_status,omitempty"`
	UserID    string          `json:"user_id,omitempty"`
	Amount    float64         `json:"amount,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	Bucket    int             `json:"bucket,omitempty"`
}


