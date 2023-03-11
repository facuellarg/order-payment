package repository

import (
	"errors"
	"testing"

	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/interface/aws"
	"github.com/stretchr/testify/assert"
)

type TestTableUpdateStatus struct {
	inputUpdate  mockInputUpdate
	outputUpdate mockOutputUpdate
}

type mockInputUpdate struct {
	orderID string
	status  entities.OrderStatus
}
type mockOutputUpdate struct {
	err error
}

func TestUpdateOrder(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestTableUpdateStatus{
		{
			inputUpdate: mockInputUpdate{
				orderID: "",
				status:  entities.Incomplete,
			},

			outputUpdate: mockOutputUpdate{errors.New("mocked error")},
		},
		{
			inputUpdate: mockInputUpdate{
				orderID: "order_test",
				status:  entities.Incomplete,
			},
			outputUpdate: mockOutputUpdate{nil},
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

		assert.ErrorIs(err, testCase.outputUpdate.err)

	}

}
