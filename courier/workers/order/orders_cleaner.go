package order

import (
	"context"
	"gitlab.com/ptflp/geotask/module/order/service"
	"log"
	"time"
)

const (
	orderCleanInterval = 5 * time.Second
)

// OrderCleaner воркер, который удаляет старые заказы
// используя метод orderService.RemoveOldOrders()
type OrderCleaner struct {
	orderService service.Orderer
}

func NewOrderCleaner(orderService service.Orderer) *OrderCleaner {
	return &OrderCleaner{orderService: orderService}
}

func (o *OrderCleaner) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(orderCleanInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := o.orderService.RemoveOldOrders(ctx)
				if err != nil {
					log.Println("Error removing old orders:", err)
				}
			}
		}
	}()
}
