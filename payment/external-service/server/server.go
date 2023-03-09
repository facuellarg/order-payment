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
	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/external-service/controller"
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

func (ps PaymentServer) ProcessPayment() {
	lambda.Start(ps.processPayment)
}
func (ps PaymentServer) CreatePayment() {
	lambda.Start(ps.createPayment)
}

func (ps *PaymentServer) processPayment(ctx context.Context, event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	var paymentRequest entities.ProcessPaymentRequest

	if err := json.NewDecoder(strings.NewReader(event.Body)).Decode(&paymentRequest); err != nil {
		fmt.Println(err)
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	fmt.Printf("paymentRequest: %v\n", paymentRequest)
	if err := ps.paymentController.ProcessPaymentRequest(paymentRequest); err != nil {
		fmt.Println(err)
		return &events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	return &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func (ps *PaymentServer) createPayment(ctx context.Context, event *events.SQSEvent) error {

	for _, message := range event.Records {
		createdOrderEvent := entities.CreatedOrderEvent{}
		if err := json.NewDecoder(strings.NewReader(message.Body)).Decode(&createdOrderEvent); err != nil {
			fmt.Println(err)
			continue
		}
		paymentID, err := ps.paymentController.CreatePayment(createdOrderEvent)
		if err != nil {
			fmt.Println(err)
		}
		log.Println("payment created ", paymentID)
	}
	return nil

}
