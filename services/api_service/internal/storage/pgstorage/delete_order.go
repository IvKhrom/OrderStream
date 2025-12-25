package pgstorage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (p *PGstorage) DeleteOrder(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	res, err := p.pool.Exec(ctx, `UPDATE orders SET status='deleted', updated_at=now() WHERE order_id=$1`, uid)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}


