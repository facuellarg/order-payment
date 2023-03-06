package main

import (
	"log"

	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/external-service/server"
	"github.com/facuellarg/payment/interface/controller"
	"github.com/facuellarg/payment/interface/event"
	"github.com/facuellarg/payment/interface/repository"
	"github.com/facuellarg/payment/use-case/service"
)

func main() {
	listenChan := make(chan entities.CreatedOrderEvent)
	sendChannel := make(chan string)
	paymentEventHandler := event.NewPaymentEventChannel(
		sendChannel,
		listenChan,
	)
	paymentRepository := repository.NewPaymentMemoryRepository()
	paymentService := service.NewPaymentService(&paymentRepository, &paymentEventHandler)
	paymentController := controller.NewPaymentController(&paymentService)
	s := server.NewPaymentServer(&paymentController, 8080)
	go func() {
		listenChan <- entities.CreatedOrderEvent{
			OrderID:    "order_test",
			TotalPrice: 40,
		}
	}()

	log.Fatal(s.Start())
}
