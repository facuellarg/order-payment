package event

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/facuellarg/order/domain/entities"
)

type OrderSQSHandler struct {
	sqsService sqsiface.SQSAPI
	queueUrl   string
}

func NewOrderSQSHandler(sqsService sqsiface.SQSAPI, queueUrl string) OrderSQSHandler {
	return OrderSQSHandler{
		sqsService,
		queueUrl,
	}
}

func (psh *OrderSQSHandler) SendOrderCreatedEvent(event entities.CreateOrderEvent) error {
	var buff strings.Builder

	err := json.NewEncoder(&buff).Encode(event)
	if err != nil {
		return err
	}
	message := buff.String()
	_, err = psh.sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &psh.queueUrl,
		MessageBody: &message,
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (psh *OrderSQSHandler) ListenCompleteOrderEvent() (string, error) {
	panic("not implemented") // TODO: Implement
}
