package service

import "github.com/facuellarg/payment/domain/entities"

type PaymentServiceI interface {
	ProcessPayment(entities.ProcessPaymentRequest) error
	CreatePayment(entities.CreatedOrderEvent) (string, error)
}
