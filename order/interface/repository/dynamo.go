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
	tableName          string
	updateExpression   string
	returnUpdateString string
	updateReplacement  map[string]*string
	awsSession         dynamodbiface.DynamoDBAPI
}

func NewORderRepositoryDynamo(awsSession dynamodbiface.DynamoDBAPI) OrderRepositoryDynamo {
	return OrderRepositoryDynamo{
		awsSession:         awsSession,
		tableName:          "orders",
		updateExpression:   "set #st = :newStatus",
		returnUpdateString: "NONE",
		updateReplacement: map[string]*string{
			"#st": aws.String("status"),
		},
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
	key := map[string]*dynamodb.AttributeValue{
		"order_id": {
			S: &orderID,
		},
	}
	strStatus := (string)(newStatus)
	value := map[string]*dynamodb.AttributeValue{
		":newStatus": {
			S: &strStatus,
		},
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 &ord.tableName,
		Key:                       key,
		UpdateExpression:          &ord.updateExpression,
		ExpressionAttributeValues: value,
		ExpressionAttributeNames:  ord.updateReplacement,
		ReturnValues:              &ord.returnUpdateString,
	}
	_, err := ord.awsSession.UpdateItem(&input)
	return err
}
