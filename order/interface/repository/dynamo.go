package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/facuellarg/order/domain/entities"
	"github.com/google/uuid"
)

type OrderRepositoryDynamo struct {
	awsSession dynamodbiface.DynamoDBAPI
}

func NewORderRepositoryDynamo(awsSession dynamodbiface.DynamoDBAPI) OrderRepositoryDynamo {
	return OrderRepositoryDynamo{
		awsSession,
	}
}

func (ord *OrderRepositoryDynamo) SaveOrder(order entities.Order) (string, error) {
	order.OrderID = uuid.NewString()
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return "", err
	}
	_, err = ord.awsSession.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("orders"),
	})
	if err != nil {
		return "", err
	}
	return order.OrderID, nil
}

func (ord *OrderRepositoryDynamo) UpdateStatus(orderID string, newStatus entities.OrderStatus) error {

	return nil
}
