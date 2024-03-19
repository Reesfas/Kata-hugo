package order

import (
	"context"
	"gitlab.com/ptflp/geotask/module/order/service"
	"log"
	"time"
)

const (
	// order generation interval
	orderGenerationInterval = 10 * time.Millisecond
	maxOrdersCount          = 200
)

// worker generates orders and put them into redis
type OrderGenerator struct {
	orderService service.Orderer
}

func NewOrderGenerator(orderService service.Orderer) *OrderGenerator {
	return &OrderGenerator{orderService: orderService}
}

func (o *OrderGenerator) Run(ctx context.Context) {
	go func() {
		// Создаем таймер, который будет срабатывать с интервалом orderGenerationInterval
		ticker := time.NewTicker(orderGenerationInterval)
		defer ticker.Stop()

		// Главный цикл генерации заказов
		for {
			select {
			case <-ticker.C:
				// Получаем текущее количество заказов
				count, err := o.orderService.GetCount(ctx)
				if err != nil {
					log.Printf("Error getting order count: %v", err)
					continue
				}

				// Если количество заказов меньше maxOrdersCount, генерируем новый заказ
				if count < maxOrdersCount {
					err = o.orderService.GenerateOrder(ctx)
					if err != nil {
						log.Printf("Error generating order: %v", err)
					}
				} else {
					time.Sleep(time.Second * 15)
				}
			}
		}
	}()
}
