package notification

type Message struct {
	UserID  string `json:"user_id"`
	OrderId string `json:"order_id"`
	Message string `json:"message"`
}
