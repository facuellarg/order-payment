package event

type PaymentEventHandlerI interface {
	SendOrderCompleteEvent(string) error
}
