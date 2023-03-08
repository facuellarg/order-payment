package service

import (
	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/use-case/event"
	"github.com/facuellarg/payment/use-case/repository"
)

type PaymentService struct {
	paymentRepository   repository.PaymentRepositoryI
	paymentEventHandler event.PaymentEventHandlerI
}

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
		return err
	}

	err = ps.paymentRepository.UpdatePaymentStatus(payment.PaymentID, entities.Complete)
	if err != nil {
		return err
	}

	ps.paymentEventHandler.SendOrderCompleteEvent(payment.OrderID) //TODO: callback if it fails

	return nil
}

func (ps *PaymentService) CreatePayment() (string, error) {
	orderCreatedEvent, err := ps.paymentEventHandler.ListenOrderCreatedEvent()
	if err != nil {
		return "", err
	}

	newPayment := entities.Payment{}
	newPayment.OrderID = orderCreatedEvent.OrderID
	newPayment.Status = entities.Incomplete
	return ps.paymentRepository.SavePayment(newPayment)
}
