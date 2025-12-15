package pgstorage

import (
	"context"
	"encoding/json"

	"github.com/ivkhr/orderstream/internal/models"
)

func (p *PGstorage) Create(ctx context.Context, o *models.Order) error {
	payload := json.RawMessage(o.Payload)
	sql := `INSERT INTO orders(order_id, user_id, amount, status, payload, bucket)
		    VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := p.pool.Exec(ctx, sql, o.OrderID, o.UserID, o.Amount, o.Status, payload, o.Bucket)
	return err
}


