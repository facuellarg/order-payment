package event

import "github.com/facuellarg/payment/domain/entities"

type PaymentEventHandlerI interface {
	SendOrderCompleteEvent(string) error
	ListenOrderCreatedEvent() (entities.CreatedOrderEvent, error)
}
