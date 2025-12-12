package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID   uuid.UUID       `json:"order_id"`
	UserID    uuid.UUID       `json:"user_id"`
	Amount    float64         `json:"amount,omitempty"`
	Status    string          `json:"status"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Bucket    int             `json:"bucket"`
}

func BucketFromUUID(id uuid.UUID, buckets int) int {
	return int(id[0] % byte(buckets))
}
