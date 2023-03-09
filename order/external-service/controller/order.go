package controller

import "github.com/facuellarg/order/domain/entities"

type ControllerOrderI interface {
	CreateOrderEvent(entities.CreateOrderRequest) (string, error)
	CompleteOrder(string) error
}
