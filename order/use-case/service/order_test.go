package service_test

import (
	"errors"
	"testing"

	"github.com/facuellarg/order/domain/entities"
	"github.com/facuellarg/order/use-case/event"
	"github.com/facuellarg/order/use-case/repository"
	"github.com/facuellarg/order/use-case/service"
	"github.com/stretchr/testify/assert"
)

type mockRepositoryInput struct {
	order entities.Order
}
type mockRepositoryReturn struct {
	orderID string
	err     error
}

type mockEventInput struct {
	createOrderEvent entities.CreateOrderEvent
}

type TestTableSaveOrder struct {
	mockInput      mockRepositoryInput
	mockOutput     mockRepositoryReturn
	mockEventInput mockEventInput
	eventCalled    bool
}

func TestSaveOrderWithSaveError(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestTableSaveOrder{
		{
			mockRepositoryInput{
				service.CreateOrderByOrderRequest(entities.CreateOrderRequest{}),
			},
			mockRepositoryReturn{
				"",
				errors.New("mocked error"),
			},
			mockEventInput{},
			false,
		},
		{
			mockRepositoryInput{
				service.CreateOrderByOrderRequest(entities.CreateOrderRequest{}),
			},
			mockRepositoryReturn{
				"order_id_test",
				nil,
			},
			mockEventInput{
				createOrderEvent: entities.CreateOrderEvent{
					OrderID: "order_id_test",
				},
			},
			true,
		},
	}

	for _, testCase := range testCases {
		repositoryMock := repository.NewMockOrderRepositoryI(t)
		eventHandlerMock := event.NewMockOrderEventHandlerI(t)

		repositoryMock.On("SaveOrder", testCase.mockInput.order).Return(testCase.mockOutput.orderID, testCase.mockOutput.err)
		repositoryMock.EXPECT().
			SaveOrder(
				testCase.mockInput.order,
			).
			Return(
				testCase.mockOutput.orderID,
				testCase.mockOutput.err,
			)

		if testCase.eventCalled {
			eventHandlerMock.EXPECT().
				SendOrderCreatedEvent(
					testCase.mockEventInput.createOrderEvent,
				).
				Return(nil)
		}

		orderService := service.NewOrderService(repositoryMock, eventHandlerMock)
		orderID, errReturned := orderService.SaveOrder(testCase.mockInput.order.CreateOrderRequest)

		assert.ErrorIs(testCase.mockOutput.err, errReturned)
		assert.Equal(testCase.mockOutput.orderID, orderID)
		expectedCalls := 0
		if testCase.eventCalled {
			expectedCalls++

		}
		eventHandlerMock.AssertNumberOfCalls(t, "SendOrderCreatedEvent", expectedCalls)
	}

}
