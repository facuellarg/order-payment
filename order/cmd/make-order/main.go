package main

import (
	"github.com/facuellarg/order/external-service/server"
	myAws "github.com/facuellarg/order/interface/aws"
	"github.com/facuellarg/order/interface/controller"
	"github.com/facuellarg/order/interface/event"
	"github.com/facuellarg/order/interface/repository"
	"github.com/facuellarg/order/use-case/service"
)

func main() {
	// orderRepository := repository.NewOderRepositoryMemory()
	orderRepository := repository.NewORderRepositoryDynamo(
		myAws.Dynamodb(),
	)

	sendOrderChan := make(chan bool)
	completeOrderChan := make(chan string)
	orderEvent := event.NewORderChanEvent(sendOrderChan, completeOrderChan)
	orderService := service.NewOrderService(&orderRepository, &orderEvent)
	orderController := controller.NewOrderController(&orderService)
	s := server.NewServer(&orderController)
	s.ServeOrder()
}
