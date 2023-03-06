package event

import "github.com/facuellarg/payment/domain/entities"

type PaymentEventChannel struct {
	sendOrderCompleteChan  chan<- string
	listenOrderCreatedChan <-chan entities.CreatedOrderEvent
}

func NewPaymentEventChannel(
	sendOrder chan<- string,
	listenOrder <-chan entities.CreatedOrderEvent,
) PaymentEventChannel {
	return PaymentEventChannel{
		sendOrder,
		listenOrder,
	}
}

func (pev *PaymentEventChannel) SendOrderCompleteEvent(orderID string) error {
	pev.sendOrderCompleteChan <- orderID
	return nil
}

func (pev *PaymentEventChannel) ListenOrderCreatedEvent() (entities.CreatedOrderEvent, error) {
	return <-pev.listenOrderCreatedChan, nil
}
