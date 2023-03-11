package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type PaymentSQSHandler struct {
	sqsService sqsiface.SQSAPI
	queueUrl   string
}

func NewPaymentSQSHandler(sqsService sqsiface.SQSAPI, queueUrl string) PaymentSQSHandler {
	return PaymentSQSHandler{
		sqsService,
		queueUrl,
	}
}

func (psh *PaymentSQSHandler) SendOrderCompleteEvent(orderID string) error {
	_, err := psh.sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &psh.queueUrl,
		MessageBody: &orderID,
	})
	if err != nil {
		fmt.Println(err)
	}
	return err
}
