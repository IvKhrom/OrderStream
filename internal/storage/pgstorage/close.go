package pgstorage

import "context"

func (p *PGstorage) Close(ctx context.Context) {
	if p.pool != nil {
		p.pool.Close()
	}
}


