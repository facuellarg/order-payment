package service

import (
	"log"

	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/use-case/event"
	"github.com/facuellarg/order/use-case/repository"
)

type OrderService struct {
	OrderRepository   repository.OrderRepositoryI
	OrderEventHandler event.OrderEventHandlerI
}

func NewOrderService(
	orderRepository repository.OrderRepositoryI,
	orderEventListener event.OrderEventHandlerI,
) OrderService {
	return OrderService{
		orderRepository,
		orderEventListener,
	}
}

func (os *OrderService) SaveOrder(orderRequest entities.CreateOrderRequest) (string, error) {
	newOrder := entities.Order{}
	newOrder.CreateOrderRequest = orderRequest
	newOrder.Status = entities.Incomplete
	orderId, err := os.OrderRepository.SaveOrder(newOrder)
	if err != nil {
		return "", err
	}
	os.sendOrderCreatedEvent(entities.CreateOrderEvent{ //TODO: callback if it fails
		OrderID:    orderId,
		TotalPrice: orderRequest.TotalPrice,
	})
	return orderId, nil
}
func (oc *OrderService) sendOrderCreatedEvent(createOrderEvent entities.CreateOrderEvent) {
	err := oc.OrderEventHandler.SendOrderCreatedEvent(createOrderEvent)
	if err != nil {
		log.Println(err)
	}

}

func (os *OrderService) UpdateStatusOrder(orderID string) error {
	return os.OrderRepository.UpdateStatus(orderID, entities.Shipping)
}
