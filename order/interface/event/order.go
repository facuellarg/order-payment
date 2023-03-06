package event

import (
	"github.com/facuellarg/order/domain/entities"
)

type OrderChanEvent struct {
	sendOrderChan     chan<- bool
	completeOrderChan <-chan string
}

func NewORderChanEvent(sendOrderChan chan bool, completeOrderChan chan string) OrderChanEvent {
	return OrderChanEvent{
		sendOrderChan,
		completeOrderChan,
	}
}

func (oce *OrderChanEvent) SendOrderCreatedEvent(createOrderEvent entities.CreateOrderEvent) error {
	oce.sendOrderChan <- true
	return nil
}

func (oce *OrderChanEvent) ListenCompleteOrderEvent() (string, error) {
	return <-oce.completeOrderChan, nil
}
