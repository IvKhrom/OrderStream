package pgstorage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ivkhr/orderstream/services/worker/internal/models"
)

func (p *PGstorage) Update(ctx context.Context, o *models.Order) error {
	payload := json.RawMessage(o.Payload)
	res, err := p.pool.Exec(ctx, `UPDATE orders SET payload=$1, amount=$2, updated_at=now() WHERE order_id=$3`, payload, o.Amount, o.OrderID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}


