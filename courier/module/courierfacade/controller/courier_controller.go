package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/ptflp/geotask/module/courierfacade/service"
)

type CourierController struct {
	courierService service.CourierFacer
}

func NewCourierController(courierService service.CourierFacer) *CourierController {
	return &CourierController{courierService: courierService}
}

func (c *CourierController) GetStatus(ctx *gin.Context) {
	// Устанавливаем задержку в 50 миллисекунд
	time.Sleep(50 * time.Millisecond)

	// Получаем статус курьера из сервиса courierService, используя метод GetStatus
	status := c.courierService.GetStatus(ctx)

	// Отправляем статус курьера в ответ
	ctx.JSON(200, status)
}

func (c *CourierController) MoveCourier(m webSocketMessage) {
	var cm CourierMove
	var err error

	switch data := m.Data.(type) {
	case string:
		err = json.Unmarshal([]byte(data), &cm)
	case []byte:
		err = json.Unmarshal(data, &cm)
	default:
		log.Println("Error: Unsupported data type in webSocketMessage")
		return
	}
	if err != nil {
		log.Println("Error deserializing courier move data:", err)
		return
	}
	// Вызываем метод MoveCourier у courierService
	c.courierService.MoveCourier(context.TODO(), cm.Direction, cm.Zoom)

	log.Println("______Обработаны данные о перемещении курьера")
}
