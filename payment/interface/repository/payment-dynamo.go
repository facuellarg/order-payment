package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/facuellarg/payment/domain/entities"
	"github.com/google/uuid"
)

type PaymentDynamoRepository struct {
	awsSession         dynamodbiface.DynamoDBAPI
	returnUpdateString string
	updateExpression   string
	queryExpression    string
	tableName          string
}

func NewPaymentDynamoRepository(awsSession dynamodbiface.DynamoDBAPI, tableName string) PaymentDynamoRepository {
	return PaymentDynamoRepository{
		awsSession:         awsSession,
		returnUpdateString: "NONE",
		updateExpression:   "set #st = :newStatus",
		queryExpression:    "order_id = :order_id",
		tableName:          tableName,
	}
}

func (pdr *PaymentDynamoRepository) SavePayment(payment entities.Payment) (string, error) {
	payment.PaymentID = uuid.NewString()
	av, err := dynamodbattribute.MarshalMap(payment)
	if err != nil {
		return "", err
	}
	_, err = pdr.awsSession.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &pdr.tableName,
	})
	if err != nil {
		return "", err
	}

	return payment.PaymentID, nil
}

func (pdr *PaymentDynamoRepository) UpdatePaymentStatus(paymentID string, newStatus entities.PaymentStatus) error {
	key := map[string]*dynamodb.AttributeValue{
		"payment_id": {
			S: &paymentID,
		},
	}

	strStatus := (string)(newStatus)
	values := map[string]*dynamodb.AttributeValue{
		":newStatus": {
			S: &strStatus,
		},
	}
	input := dynamodb.UpdateItemInput{
		TableName:                 &pdr.tableName,
		Key:                       key,
		UpdateExpression:          &pdr.updateExpression,
		ExpressionAttributeValues: values,
		ExpressionAttributeNames: map[string]*string{
			"#st": aws.String("status"),
		},
		ReturnValues: &pdr.returnUpdateString,
	}

	_, err := pdr.awsSession.UpdateItem(&input)

	return err
}

func (pdr *PaymentDynamoRepository) GetPaymentByOrderID(orderID string) (entities.Payment, error) {
	fmt.Printf("orderID: %v\n", orderID)
	input := &dynamodb.ScanInput{
		TableName:        &pdr.tableName,
		FilterExpression: &pdr.queryExpression,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":order_id": {
				S: &orderID,
			},
		},
		Limit: aws.Int64(1),
	}

	result, err := pdr.awsSession.Scan(input)

	payment := entities.Payment{}
	if err != nil {
		return payment, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &payment)
	return payment, err

}
