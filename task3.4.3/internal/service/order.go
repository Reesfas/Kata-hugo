package service

import (
	"context"
	"task3.4.3/internal/repository"
)

type OrderService interface {
	Create(ctx context.Context, order repository.Order) error
	GetByID(ctx context.Context, id string) (repository.Order, error)
	Delete(ctx context.Context, id string) error
	GetInventory(ctx context.Context) (map[string]int, error)
}

type OrderServ struct {
	repo repository.OrderRepository
}

func NewOrderServ(repo repository.OrderRepository) OrderServ {
	return OrderServ{repo: repo}
}

func (o *OrderServ) Create(ctx context.Context, order repository.Order) error {
	return o.repo.Create(ctx, order)
}

func (o *OrderServ) GetByID(ctx context.Context, id string) (repository.Order, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *OrderServ) Delete(ctx context.Context, id string) error {
	return o.repo.Delete(ctx, id)
}
func (o *OrderServ) GetInventory(ctx context.Context) (map[string]int, error) {
	return o.repo.GetInventory(ctx)
}
