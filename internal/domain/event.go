package domain

import (
	"encoding/json"
	"time"
)

type OrderEvent struct {
	EventID    string          `json:"event_id"`    // "0"=create, ">0"=update_id
	OrderID    string          `json:"order_id"`    // Internal order ID
	ExternalID string          `json:"external_id"` // External order ID
	UserID     string          `json:"user_id"`     // User ID for security
	Payload    json.RawMessage `json:"payload,omitempty"`
	Status     string          `json:"status"` // new, processing, done, cancelled
	Timestamp  time.Time       `json:"timestamp"`
}

type OrderAck struct {
	EventID     string    `json:"event_id"`
	OrderID     string    `json:"order_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
