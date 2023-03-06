package entities

type ProcessPaymentRequest struct {
	OrderID string        `json:"order_id"`
	Status  PaymentStatus `json:"status"`
}

type CreatedOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

type Payment struct {
	ProcessPaymentRequest
	PaymentID string `json:"payment_id"`
}

type PaymentStatus string

const (
	Incomplete PaymentStatus = "incomplete"
	Complete   PaymentStatus = "complete"
)
