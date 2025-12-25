package models

import (
	"encoding/json"
	"time"
)

type OrderEvent struct {
	OrderID    string          `json:"order_id"`    // internal UUID as string
	ExternalID string          `json:"external_id"` // marketplace id extracted from payload
	UserID     string          `json:"user_id"`     // user UUID as string
	Payload    json.RawMessage `json:"payload,omitempty"`
	Status     string          `json:"status"` // new|processing|done|cancelled|deleted
	Timestamp  time.Time       `json:"timestamp"`
}


