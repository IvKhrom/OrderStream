package pgstorage

import (
	"context"
	"encoding/json"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

func (p *PGstorage) Upsert(ctx context.Context, ord *models.Order) (bool, error) {
	payload := json.RawMessage(ord.Payload)

	sql := `INSERT INTO orders(order_id, user_id, amount, status, payload, bucket, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			ON CONFLICT (order_id) DO UPDATE SET
				user_id=EXCLUDED.user_id,
				amount=EXCLUDED.amount,
				status=EXCLUDED.status,
				payload=EXCLUDED.payload,
				updated_at=now()
			RETURNING (orders.created_at = $7) AS inserted`

	var inserted bool
	if err := p.pool.QueryRow(ctx, sql,
		ord.OrderID, ord.UserID, ord.Amount, ord.Status, payload, ord.Bucket,
		ord.CreatedAt, ord.UpdatedAt,
	).Scan(&inserted); err != nil {
		return false, err
	}
	return inserted, nil
}
