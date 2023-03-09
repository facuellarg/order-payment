package service

import "github.com/facuellarg/order/domain/entities"

type OrderServiceI interface {
	SaveOrder(entities.CreateOrderRequest) (string, error)
	UpdateStatusOrder(string) error
}
