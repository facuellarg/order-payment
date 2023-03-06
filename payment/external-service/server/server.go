package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/external-service/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PaymentServer struct {
	paymentController controller.PaymentControllerI
	port              uint
}

func NewPaymentServer(
	paymentController controller.PaymentControllerI,
	port uint,
) PaymentServer {
	return PaymentServer{
		paymentController,
		port,
	}
}

func (ps PaymentServer) Start() error {
	server := echo.New()
	server.Use(middleware.Logger())
	server.POST("/payment", ps.processPayment)
	go ps.createPayment()
	return server.Start(fmt.Sprintf(":%d", ps.port))

}

func (ps *PaymentServer) processPayment(ctx echo.Context) error {
	var paymentRequest entities.ProcessPaymentRequest
	if err := ctx.Bind(&paymentRequest); err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}
	if err := ps.paymentController.ProcessPaymentRequest(paymentRequest); err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	return ctx.NoContent(http.StatusOK)
}

func (ps *PaymentServer) createPayment() {
	for {
		paymentID, err := ps.paymentController.CreatePayment()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Payment created %s", paymentID)
	}
}
