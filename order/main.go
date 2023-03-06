package main

import (
	"log"

	"github.com/facuellarg/order/external-service/server"
	"github.com/facuellarg/order/interface/controller"
	"github.com/facuellarg/order/interface/event"
	"github.com/facuellarg/order/interface/repository"
	"github.com/facuellarg/order/use-case/service"
)

func main() {
	orderRepository := repository.NewOderRepositoryMemory()

	sendOrderChan := make(chan bool)
	completeOrderChan := make(chan string)
	orderEvent := event.NewORderChanEvent(sendOrderChan, completeOrderChan)
	orderService := service.NewOrderService(&orderRepository, &orderEvent)
	orderController := controller.NewOrderController(&orderService)
	s := server.NewServer(&orderController, 8080)
	// go func() {
	// 	for {
	// 		<-sendOrderChan
	// 		fmt.Println("Order event sent")
	// 	}
	// }()
	// go func() {
	// 	rand.Seed(time.Now().Unix())
	// 	for {
	// 		<-time.After(1 * time.Second)
	// 		if rand.Float32() < 0.1 {
	// 			completeOrderChan <- "Some ID"
	// 		}
	// 	}
	// }()
	log.Fatal(s.Start())
}
