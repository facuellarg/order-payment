package repository

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/facuellarg/payment/domain/entities"
	"github.com/google/uuid"
)

type PaymentMemoryRepository struct {
	payments map[string]entities.Payment
}

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

func NewPaymentMemoryRepository() PaymentMemoryRepository {
	return PaymentMemoryRepository{
		payments: make(map[string]entities.Payment),
	}
}

func (pr *PaymentMemoryRepository) SavePayment(payment entities.Payment) (string, error) {
	payment.PaymentID = uuid.NewString()
	pr.payments[payment.PaymentID] = payment
	pr.Write()
	return payment.PaymentID, nil
}

func (pr *PaymentMemoryRepository) UpdatePaymentStatus(paymentID string, paymentStatus entities.PaymentStatus) error {
	payment, ok := pr.payments[paymentID]
	if !ok {
		return ErrPaymentNotFound
	}
	payment.Status = paymentStatus
	pr.payments[paymentID] = payment
	pr.Write()
	return nil
}

func (pr *PaymentMemoryRepository) GetPaymentByOrderID(orderID string) (entities.Payment, error) {
	for _, v := range pr.payments {
		if v.OrderID == orderID {
			return v, nil

		}
	}
	return entities.Payment{}, nil
}

func (pr *PaymentMemoryRepository) Write() {
	file, err := os.OpenFile("memory.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(pr.payments)
}
