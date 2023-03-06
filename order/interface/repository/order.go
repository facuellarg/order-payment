package repository

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/facuellarg/order/domain/entities"
	"github.com/google/uuid"
)

type OrderRepositoryMemory struct {
	orders map[string]entities.Order
}

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotExists     = errors.New("order doesn't exists")
)

func NewOderRepositoryMemory() OrderRepositoryMemory {
	return OrderRepositoryMemory{
		make(map[string]entities.Order),
	}
}

func (or *OrderRepositoryMemory) SaveOrder(order entities.Order) (string, error) {
	order.OrderID = uuid.NewString()
	// order.OrderID = "Some ID"
	if _, ok := or.orders[order.OrderID]; ok {
		return "", ErrOrderAlreadyExists
	}
	or.orders[order.OrderID] = order
	or.Write()
	return order.OrderID, nil
}

func (or *OrderRepositoryMemory) UpdateStatus(orderID string, newStatus entities.OrderStatus) error {
	order, ok := or.orders[orderID]
	if !ok {
		return ErrOrderNotExists
	}
	order.Status = newStatus
	or.orders[orderID] = order
	or.Write()
	return nil
}

func (or *OrderRepositoryMemory) Write() error {
	file, err := os.OpenFile("memory.csv", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(or.orders)

}
