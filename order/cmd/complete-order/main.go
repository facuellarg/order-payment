package main

import (
	"os"

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

	orderEvent := event.NewOrderSQSHandler(myAws.SQS(), os.Getenv("QUEUE_URL"))
	orderService := service.NewOrderService(&orderRepository, &orderEvent)
	orderController := controller.NewOrderController(&orderService)
	s := server.NewServer(&orderController)
	s.ServeListenOrderComplete()
}
