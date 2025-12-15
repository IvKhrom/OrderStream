package pgstorage

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *PGstorage) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error) {
	row := p.pool.QueryRow(ctx,
		`SELECT order_id, user_id, amount, status, payload, created_at, updated_at, bucket
		 FROM orders WHERE payload->> 'id' = $1 AND user_id=$2 LIMIT 1`,
		externalID, userID)

	var o models.Order
	var payload []byte
	if err := row.Scan(&o.OrderID, &o.UserID, &o.Amount, &o.Status, &payload, &o.CreatedAt, &o.UpdatedAt, &o.Bucket); err != nil {
		return nil, err
	}
	o.Payload = json.RawMessage(payload)
	return &o, nil
}


