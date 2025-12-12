package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ivkhr/orderstream/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool Pool
}

func New(dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close(ctx context.Context) {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *Postgres) CreateOrder(ctx context.Context, o *domain.Order) error {
	payload := json.RawMessage(o.Payload)
	sql := `INSERT INTO orders(order_id, user_id, amount, status, payload, bucket)
		    VALUES ($1,$2,$3,$4,$5,$6)`

	// external_id не хранится как отдельная колонка — внешние id берутся из JSON payload
	_, err := p.pool.Exec(ctx, sql, o.OrderID, o.UserID, o.Amount, o.Status, payload, o.Bucket)
	return err
}

func (p *Postgres) Create(ctx context.Context, o *domain.Order) error {
	return p.CreateOrder(ctx, o)
}

func (p *Postgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	row := p.pool.QueryRow(ctx,
		`SELECT order_id, user_id, amount, status, payload, created_at, updated_at, bucket
		 FROM orders WHERE order_id=$1`, id)

	var o domain.Order
	var payload []byte
	if err := row.Scan(&o.OrderID, &o.UserID, &o.Amount, &o.Status, &payload, &o.CreatedAt, &o.UpdatedAt, &o.Bucket); err != nil {
		return nil, err
	}
	o.Payload = json.RawMessage(payload)
	return &o, nil
}

func (p *Postgres) GetByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*domain.Order, error) {
	// Находим по id маркетплейса внутри JSON payload (payload->>'id') и по user_id
	row := p.pool.QueryRow(ctx,
		`SELECT order_id, user_id, amount, status, payload, created_at, updated_at, bucket
		 FROM orders WHERE payload->> 'id' = $1 AND user_id=$2 LIMIT 1`,
		externalID, userID)

	var o domain.Order
	var payload []byte
	if err := row.Scan(&o.OrderID, &o.UserID, &o.Amount, &o.Status, &payload, &o.CreatedAt, &o.UpdatedAt, &o.Bucket); err != nil {
		return nil, err
	}
	o.Payload = json.RawMessage(payload)
	return &o, nil
}

func (p *Postgres) Update(ctx context.Context, o *domain.Order) error {
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

func (p *Postgres) UpdateStatus(ctx context.Context, id string, status string) error {
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

func (p *Postgres) DeleteOrder(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Мягкое (soft) удаление
	res, err := p.pool.Exec(ctx, `UPDATE orders SET status='deleted', updated_at=now() WHERE order_id=$1`, uid)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
