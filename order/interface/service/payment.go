package service

import "github.com/facuellarg/order/domain/entities"

type PaymentServiceI interface {
	OrderCreateEvent(entities.CreateOrderEvent) error
}
