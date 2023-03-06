package repository

import "github.com/facuellarg/payment/domain/entities"

type PaymentRepositoryI interface {
	SavePayment(entities.Payment) (string, error)
	UpdatePaymentStatus(string, entities.PaymentStatus) error
	GetPaymentByOrderID(string) (entities.Payment, error)
}
