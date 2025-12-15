package pgstorage

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *PGstorage) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	row := p.pool.QueryRow(ctx,
		`SELECT order_id, user_id, amount, status, payload, created_at, updated_at, bucket
		 FROM orders WHERE order_id=$1`, id)

	var o models.Order
	var payload []byte
	if err := row.Scan(&o.OrderID, &o.UserID, &o.Amount, &o.Status, &payload, &o.CreatedAt, &o.UpdatedAt, &o.Bucket); err != nil {
		return nil, err
	}
	o.Payload = json.RawMessage(payload)
	return &o, nil
}


