package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/external-service/controller"
)

type Server struct {
	orderController controller.ControllerOrderI
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

func (s *Server) ServeListenOrderComplete() error {
	lambda.Start(s.listenOrderComplete)
	return nil
}

func (s *Server) createOrder(ctx context.Context, event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	var order entities.CreateOrderRequest
	err := json.NewDecoder(strings.NewReader(event.Body)).Decode(&order)
	if err != nil {
		fmt.Println(err)
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

func (s *Server) listenOrderComplete(ctx context.Context, event *events.SQSEvent) error {
	for _, message := range event.Records {
		if message.Body == "" {
			continue
		}
		if err := s.orderController.CompleteOrder(message.Body); err != nil {
			fmt.Printf("error completing order id: %s\nerr: %s\n", message.Body, err)
		}
	}
	return nil
}
