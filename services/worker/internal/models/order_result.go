package models

import (
	"encoding/json"
	"time"
)

type OrderResult struct {
	OrderID string `json:"order_id"`
	// Status — результат операции для API (created|updated|deleted).
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`

	// OrderStatus — фактический статус заказа (new|processing|done|deleted).
	OrderStatus string          `json:"order_status,omitempty"`
	UserID      string          `json:"user_id,omitempty"`
	Amount      float64         `json:"amount,omitempty"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
	Bucket      int             `json:"bucket,omitempty"`
}
