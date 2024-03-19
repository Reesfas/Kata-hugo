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
	var data []byte
	if err := json.Unmarshal(data, &cm); err != nil {
		log.Fatal(err)
	}

	// Вызываем метод MoveCourier у courierService
	c.courierService.MoveCourier(context.Background(), cm.Direction, cm.Zoom)

	log.Println("______Обработаны данные о перемещении курьера")
}
