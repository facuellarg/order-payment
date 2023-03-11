package service

import (
	"log"

	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/use-case/event"
	"github.com/facuellarg/order/use-case/repository"
)

type OrderService struct {
	orderRepository   repository.OrderRepositoryI
	orderEventHandler event.OrderEventHandlerI
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

func CreateOrderByOrderRequest(orderRequest entities.CreateOrderRequest) entities.Order {
	return entities.Order{
		CreateOrderRequest: orderRequest,
		Status:             entities.Incomplete,
	}

}
func (os *OrderService) SaveOrder(orderRequest entities.CreateOrderRequest) (string, error) {
	newOrder := CreateOrderByOrderRequest(orderRequest)
	orderId, err := os.orderRepository.SaveOrder(newOrder)
	if err != nil {
		return "", err
	}

	createEvent := entities.CreateOrderEvent{
		OrderID:    orderId,
		TotalPrice: orderRequest.TotalPrice,
	}

	os.sendOrderCreatedEvent(createEvent)
	return orderId, nil
}
func (oc *OrderService) sendOrderCreatedEvent(createOrderEvent entities.CreateOrderEvent) {
	err := oc.orderEventHandler.SendOrderCreatedEvent(createOrderEvent)
	if err != nil {
		log.Println(err)
	}

}

func (os *OrderService) UpdateStatusOrder(orderID string) error {
	return os.orderRepository.UpdateStatus(orderID, entities.Shipping)
}
