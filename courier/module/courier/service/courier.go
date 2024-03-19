package service

import (
	"context"
	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/module/courier/models"
	"gitlab.com/ptflp/geotask/module/courier/storage"
	"math"
)

// Направления движения курьера
const (
	DirectionUp    = 0
	DirectionDown  = 1
	DirectionLeft  = 2
	DirectionRight = 3
)

const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)

type Courierer interface {
	GetCourier(ctx context.Context) (*models.Courier, error)
	MoveCourier(courier models.Courier, direction, zoom int) error
}

type CourierService struct {
	courierStorage storage.CourierStorager
	allowedZone    geo.PolygonChecker
	disabledZones  []geo.PolygonChecker
}

func NewCourierService(courierStorage storage.CourierStorager, allowedZone geo.PolygonChecker, disbledZones []geo.PolygonChecker) Courierer {
	return &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZones: disbledZones}
}

func (c *CourierService) GetCourier(ctx context.Context) (*models.Courier, error) {
	// Получаем курьера из хранилища
	courier, err := c.courierStorage.GetOne(ctx)
	if err != nil {
		return nil, err
	}

	if !c.allowedZone.Contains(geo.Point(courier.Location)) {
		newLocation := c.allowedZone.RandomPoint()
		courier.Location = models.Point(newLocation)

		err = c.courierStorage.Save(ctx, *courier)
		if err != nil {
			return nil, err
		}
	}

	return courier, nil
}

// MoveCourier : direction - направление движения курьера, zoom - зум карты
func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
	// Вычисляем точность перемещения в зависимости от зума карты
	precision := 0.001 / math.Pow(2, float64(zoom-14))

	// Определяем направление перемещения
	var latDiff, lngDiff float64
	switch direction {
	case DirectionUp:
		latDiff = precision
	case DirectionDown:
		latDiff = -precision
	case DirectionLeft:
		lngDiff = -precision
	case DirectionRight:
		lngDiff = precision
	}

	// Вычисляем новые координаты курьера
	newLat := courier.Location.Lat + latDiff
	newLng := courier.Location.Lng + lngDiff

	// Создаем новую точку для проверки выхода за границы зоны
	newLocation := models.Point{Lat: newLat, Lng: newLng}

	// Проверяем, находится ли курьер в разрешенной зоне
	if !c.allowedZone.Contains(geo.Point(newLocation)) {
		// Курьер вышел за границы разрешенной зоны
		// Перемещаем его в случайную точку внутри зоны
		newLocation = models.Point(c.allowedZone.RandomPoint())
	}

	// Обновляем координаты курьера
	courier.Location = newLocation

	// Сохраняем изменения в хранилище
	err := c.courierStorage.Save(context.Background(), courier)
	if err != nil {
		return err
	}

	return nil
}
