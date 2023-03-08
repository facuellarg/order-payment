package event

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/facuellarg/payment/domain/entities"
)

type PaymentSQSHandler struct {
	sqsService sqsiface.SQSAPI
	queueUrl   string
}

func (psh *PaymentSQSHandler) SendOrderCompleteEvent(orderID string) error {
	_, err := psh.sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &psh.queueUrl,
		MessageBody: &orderID,
	})
	return err
	// panic("not implemented") // TODO: Implement
}

func (psh *PaymentSQSHandler) ListenOrderCreatedEvent() (entities.CreatedOrderEvent, error) {
	panic("not implemented") // TODO: Implement
}
