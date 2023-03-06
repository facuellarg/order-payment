package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/external-service/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	orderController controller.ControllerOrderI
	port            int
	// logger	log.Logger
}

func NewServer(orderController controller.ControllerOrderI, port int) Server {
	s := Server{
		orderController,
		port,
	}
	return s
}

func (s *Server) Start() error {
	server := echo.New()
	server.Use(middleware.Logger())
	server.POST("/order", s.createOrder)
	go s.CompleteOrder()
	return server.Start(fmt.Sprintf(":%d", s.port))

}

func (s *Server) createOrder(ctx echo.Context) error {

	var order entities.CreateOrderRequest
	if err := ctx.Bind(&order); err != nil {
		ctx.Logger().Error(err) //TODO: Change logger to a real one
		return echo.ErrBadRequest
	}
	orderID, err := s.orderController.CreateOrderEvent(order)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	return ctx.String(http.StatusCreated, orderID)
}

func (s *Server) CompleteOrder() {
	for {
		orderID, err := s.orderController.CompleteOrder() //TODO: same logger problem
		if err != nil {
			log.Printf("error completing order orderID:%s\n error:%s", orderID, err)
			continue
		}
		fmt.Printf("order %s completed\n", orderID)

	}
}
