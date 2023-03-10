package entities

type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

type Order struct {
	CreateOrderRequest
	OrderID string      `json:"order_id"`
	Status  OrderStatus `json:"status"`
}

type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

type ProcessPaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type OrderStatus string

const (
	Incomplete OrderStatus = "incomplete"
	Shipping   OrderStatus = "shipping"
)
