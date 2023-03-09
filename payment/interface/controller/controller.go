package controller

import (
	"github.com/facuellarg/payment/domain/entities"
	"github.com/facuellarg/payment/interface/service"
)

type PaymentController struct {
	paymentService service.PaymentServiceI
}

func NewPaymentController(paymentService service.PaymentServiceI) PaymentController {
	return PaymentController{
		paymentService,
	}
}

func (pc *PaymentController) ProcessPaymentRequest(paymentRequest entities.ProcessPaymentRequest) error {
	return pc.paymentService.ProcessPayment(paymentRequest)
}

func (pc *PaymentController) CreatePayment(orderCreatedEvent entities.CreatedOrderEvent) (string, error) {
	return pc.paymentService.CreatePayment(orderCreatedEvent)
}
