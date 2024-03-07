package repository

import (
	"context"
	"database/sql"
)

type Order struct {
	id       int
	petId    int
	quantity int
	shipdate string
	status   string
	complete bool
}

type OrderRepository interface {
	Create(ctx context.Context, order Order) error
	GetByID(ctx context.Context, id string) (Order, error)
	Delete(ctx context.Context, id string) error
	GetInventory(ctx context.Context) (map[string]int, error)
}

type OrderRep struct {
	db *sql.DB
}

func NewOrderRep(db *sql.DB) *OrderRep {
	return &OrderRep{db: db}
}

func (o *OrderRep) Create(ctx context.Context, order Order) error {
	query := `INSERT INTO order (petId, quantity, shipdate, status, complete) VALUES ($1, $2, $3, $4, $5)`
	_, err := o.db.ExecContext(ctx, query, order.petId, order.quantity, order.shipdate, order.status, order.complete)
	return err
}

func (o *OrderRep) GetByID(ctx context.Context, id string) (Order, error) {
	order := Order{}
	query := `SELECT * FROM order WHERE id = $1 AND deleted = false`
	err := o.db.QueryRowContext(ctx, string(query), id).Scan(&order.petId, &order.quantity, &order.shipdate, &order.status, &order.complete)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (o *OrderRep) Delete(ctx context.Context, id string) error {
	query := `UPDATE order SET deleted = true WHERE id = $1 AND deleted = false`
	_, err := o.db.ExecContext(ctx, query, id)
	return err
}

func (o *OrderRep) GetInventory(ctx context.Context) (map[string]int, error) {
	rows, err := o.db.QueryContext(ctx, "SELECT status, quantity FROM order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make(map[string]int)

	for rows.Next() {
		var status string
		var quantity int
		if err = rows.Scan(&status, &quantity); err != nil {
			return nil, err
		}
		inventory[status] = quantity
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventory, nil
}
