package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/facuellarg/order/domain/entities"
	"github.com/google/uuid"
)

var (
	tableName          = "orders"
	updateExpression   = "set #st = :newStatus"
	returnUpdateString = "NONE"
	updateReplacement  = map[string]*string{
		"#st": aws.String("status"),
	}
)

type OrderRepositoryDynamo struct {
	awsSession dynamodbiface.DynamoDBAPI
}

func NewORderRepositoryDynamo(awsSession dynamodbiface.DynamoDBAPI) OrderRepositoryDynamo {
	return OrderRepositoryDynamo{
		awsSession: awsSession,
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

func buildInput(orderID string, status entities.OrderStatus) dynamodb.UpdateItemInput {

	key := map[string]*dynamodb.AttributeValue{
		"order_id": {
			S: &orderID,
		},
	}
	strStatus := (string)(status)
	value := map[string]*dynamodb.AttributeValue{
		":newStatus": {
			S: &strStatus,
		},
	}

	return dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       key,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: value,
		ExpressionAttributeNames:  updateReplacement,
		ReturnValues:              &returnUpdateString,
	}
}

func (ord *OrderRepositoryDynamo) UpdateStatus(orderID string, newStatus entities.OrderStatus) error {

	input := buildInput(orderID, newStatus)
	_, err := ord.awsSession.UpdateItem(&input)
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
		return nil
	}
	return err
}
