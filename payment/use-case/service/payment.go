package service

import (
	"fmt"

	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/use-case/event"
	"github.com/facuellarg/payment/use-case/repository"
)

type PaymentService struct {
	paymentRepository   repository.PaymentRepositoryI
	paymentEventHandler event.PaymentEventHandlerI
}

const (
	ErrPaymentNotFound = "payment not found"
)

func NewPaymentService(
	paymentRepository repository.PaymentRepositoryI,
	paymentEventHandler event.PaymentEventHandlerI,
) PaymentService {
	return PaymentService{
		paymentRepository,
		paymentEventHandler,
	}
}

func (ps *PaymentService) ProcessPayment(processPaymentRequest entities.ProcessPaymentRequest) error {
	payment, err := ps.paymentRepository.GetPaymentByOrderID(processPaymentRequest.OrderID)
	if err != nil {
		if perr, ok := err.(*entities.Error); ok && perr.Code == ErrPaymentNotFound {
			fmt.Printf(perr.Message)
			return nil
		}
		return err
	}

	err = ps.paymentRepository.UpdatePaymentStatus(payment.PaymentID, entities.Complete)
	if err != nil {
		return err
	}

	ps.paymentEventHandler.SendOrderCompleteEvent(payment.OrderID)

	return nil
}

func (ps *PaymentService) CreatePayment(orderCreatedEvent entities.CreatedOrderEvent) (string, error) {

	newPayment := entities.Payment{}
	newPayment.OrderID = orderCreatedEvent.OrderID
	newPayment.Status = entities.Incomplete
	return ps.paymentRepository.SavePayment(newPayment)
}
