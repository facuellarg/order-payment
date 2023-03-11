package repository

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/interface/aws"
	"github.com/stretchr/testify/assert"
)

type TestTableUpdateStatus struct {
	inputUpdate  mockInputUpdate
	outputUpdate mockOutputUpdate
	errExpected  error
}

type mockInputUpdate struct {
	orderID string
	status  entities.OrderStatus
}
type mockOutputUpdate struct {
	err error
}

var (
	errMocked = errors.New("mocked error")
)

func TestUpdateOrder(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestTableUpdateStatus{
		{
			inputUpdate: mockInputUpdate{
				orderID: "",
				status:  entities.Incomplete,
			},

			outputUpdate: mockOutputUpdate{errMocked},
			errExpected:  errMocked,
		},
		{
			inputUpdate: mockInputUpdate{
				orderID: "",
				status:  entities.Incomplete,
			},

			outputUpdate: mockOutputUpdate{
				awserr.New(
					dynamodb.ErrCodeConditionalCheckFailedException,
					"",
					nil,
				),
			},
			errExpected: nil,
		},
		{
			inputUpdate: mockInputUpdate{
				orderID: "order_test",
				status:  entities.Incomplete,
			},
			outputUpdate: mockOutputUpdate{nil},
			errExpected:  nil,
		},
	}
	for _, testCase := range testCases {
		awsSessionMock := aws.NewMockDynamoDBI(t)
		dynamoRepository := NewORderRepositoryDynamo(awsSessionMock)
		input := buildInput(
			testCase.inputUpdate.orderID,
			testCase.inputUpdate.status,
		)
		awsSessionMock.EXPECT().UpdateItem(&input).Return(nil, testCase.outputUpdate.err)

		err := dynamoRepository.UpdateStatus(testCase.inputUpdate.orderID, testCase.inputUpdate.status)
		assert.ErrorIs(err, testCase.errExpected)

	}

}
