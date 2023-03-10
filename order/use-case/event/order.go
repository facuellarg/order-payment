package event

import "github.com/facuellarg/order/domain/entities"

type OrderEventHandlerI interface {
	SendOrderCreatedEvent(entities.CreateOrderEvent) error
}
