package controller

import (
	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/interface/service"
)

type OrderController struct {
	OrderService service.OrderServiceI
}

func NewOrderController(orderService service.OrderServiceI) OrderController {
	return OrderController{
		orderService,
	}
}

func (oc *OrderController) CreateOrderEvent(orderRequest entities.CreateOrderRequest) (string, error) {
	return oc.OrderService.SaveOrder(orderRequest)

}

func (oc *OrderController) CompleteOrder(orderID string) error {
	return oc.OrderService.UpdateStatusOrder(orderID)
}
