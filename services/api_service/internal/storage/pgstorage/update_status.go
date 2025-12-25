package pgstorage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (p *PGstorage) UpdateStatus(ctx context.Context, id string, status string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	res, err := p.pool.Exec(ctx, `UPDATE orders SET status=$1, updated_at=now() WHERE order_id=$2`, status, uid)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}


