package payment

type Payment struct {
	OrderId         string  `json:"orderId"`
	PaymentMethodId string  `json:"paymentMethodId"`
	Amount          float64 `json:"amount"`
}
