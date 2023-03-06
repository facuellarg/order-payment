package controller

import "github.com/facuellarg/payment/domain/entities"

type PaymentControllerI interface {
	ProcessPaymentRequest(entities.ProcessPaymentRequest) error
	CreatePayment() (string, error)
}
