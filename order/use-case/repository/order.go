package repository

import "github.com/facuellarg/order/domain/entities"

type OrderRepositoryI interface {
	SaveOrder(entities.Order) (string, error)
	UpdateStatus(string, entities.OrderStatus) error
}
