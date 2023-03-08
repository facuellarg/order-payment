package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/external-service/controller"
)

type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

type Server struct {
	orderController controller.ControllerOrderI
	// logger	log.Logger
}

func NewServer(orderController controller.ControllerOrderI) Server {
	s := Server{
		orderController,
	}
	return s
}

func (s *Server) ServeOrder() error {
	lambda.Start(s.createOrder)
	return nil
}
func (s *Server) createOrder(ctx context.Context, event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	fmt.Println("creating order")
	var order entities.CreateOrderRequest
	err := json.NewDecoder(strings.NewReader(event.Body)).Decode(&order)
	if err != nil {

		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	orderID, err := s.orderController.CreateOrderEvent(order)
	if err != nil {
		fmt.Println(err)
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
		Body:       orderID,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
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
