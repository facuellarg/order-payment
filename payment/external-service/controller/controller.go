package controller

import "github.com/facuellarg/payment/domain/entities"

type PaymentControllerI interface {
	ProcessPaymentRequest(entities.ProcessPaymentRequest) error
	CreatePayment(entities.CreatedOrderEvent) (string, error)
}
