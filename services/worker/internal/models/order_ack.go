package models

import "time"

type OrderAck struct {
	OrderID     string    `json:"order_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}


