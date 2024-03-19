package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.com/ptflp/geotask/module/order/models"
	"time"
)

type OrderStorager interface {
	Save(ctx context.Context, order models.Order, maxAge time.Duration) error                       // сохранить заказ с временем жизни
	GetByID(ctx context.Context, orderID int) (*models.Order, error)                                // получить заказ по id
	GenerateUniqueID(ctx context.Context) (int64, error)                                            // сгенерировать уникальный id
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // получить заказы в радиусе от точки
	GetCount(ctx context.Context) (int, error)                                                      // получить количество заказов
	RemoveOldOrders(ctx context.Context, maxAge time.Duration) error                                // удалить старые заказы по истечению времени maxAge
}

type OrderStorage struct {
	storage *redis.Client
}

func NewOrderStorage(storage *redis.Client) OrderStorager {
	return &OrderStorage{storage: storage}
}

func (o *OrderStorage) Save(ctx context.Context, order models.Order, maxAge time.Duration) error {
	// Преобразуем заказ в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Сохраняем заказ в Redis с временем жизни
	err = o.storage.Set(fmt.Sprintf("order:%d", order.ID), orderJSON, maxAge).Err()
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderStorage) RemoveOldOrders(ctx context.Context, maxAge time.Duration) error {
	// Вычисляем максимальное время для удаления старых заказов
	maxTime := time.Now().Add(-maxAge).Unix()

	// Получаем ID всех старых ордеров с помощью ZRangeByScore
	oldOrderIDs, err := o.storage.ZRangeByScore("orders", redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%d", maxTime),
		Offset: 0,
		Count:  -1,
	}).Result()
	if err != nil {
		return err
	}

	if len(oldOrderIDs) == 0 {
		return nil
	}

	var oldOrderIDsInterface []interface{}
	for _, id := range oldOrderIDs {
		oldOrderIDsInterface = append(oldOrderIDsInterface, id)
	}

	_, err = o.storage.ZRem("orders", oldOrderIDsInterface...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderStorage) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	// Получаем данные о заказе из Redis по ключу order:ID
	data, err := o.storage.Get(fmt.Sprintf("order:%d", orderID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order models.Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *OrderStorage) saveOrderWithGeo(ctx context.Context, order models.Order, maxAge time.Duration) error {
	var err error
	var data []byte

	data, err = json.Marshal(order)
	if err != nil {
		return err
	}

	err = o.storage.Set(fmt.Sprintf("order:%d", order.ID), data, maxAge).Err()
	if err != nil {
		return err
	}

	err = o.storage.GeoAdd("orders", &redis.GeoLocation{
		Name:      fmt.Sprintf("order:%d", order.ID),
		Longitude: order.Lng,
		Latitude:  order.Lat,
	}).Err()
	if err != nil {
		return err
	}

	err = o.storage.ZAdd("orders", redis.Z{
		Score:  float64(time.Now().Unix()), // Время создания заказа
		Member: fmt.Sprintf("order:%d", order.ID),
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderStorage) GetCount(ctx context.Context) (int, error) {
	count, err := o.storage.ZCard("orders").Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (o *OrderStorage) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	var err error
	var orders []models.Order
	var ordersLocation []redis.GeoLocation

	ordersLocation, err = o.getOrdersByRadius(ctx, lng, lat, radius, unit)
	if err != nil {
		return nil, err
	}

	if len(ordersLocation) == 0 {
		return nil, nil
	}

	orders = make([]models.Order, 0, len(ordersLocation))

	for _, loc := range ordersLocation {
		orderData, err := o.storage.Get(loc.Name).Bytes()
		if err != nil {
			return nil, err
		}

		var order models.Order
		if err = json.Unmarshal(orderData, &order); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderStorage) getOrdersByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]redis.GeoLocation, error) {
	// В данном методе мы получаем список ордеров в радиусе от точки.
	// Возвращаем список ордеров с координатами и расстоянием до точки.

	// Создаем запрос для получения ордеров в заданном радиусе.
	query := &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	}

	orders, err := o.storage.GeoRadius("orders", lng, lat, query).Result()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderStorage) GenerateUniqueID(ctx context.Context) (int64, error) {
	id, err := o.storage.Incr("order:id").Result()
	if err != nil {
		return 0, err
	}
	return id, nil
}
