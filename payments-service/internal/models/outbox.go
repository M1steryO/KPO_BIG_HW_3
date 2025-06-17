package models

type paymentEventType = string
type paymentMessage = string

const (
	PaymentStatusFailed  paymentEventType = "failed"
	PaymentStatusSuccess paymentEventType = "success"

	UserNotFoundMessage   paymentMessage = "user not found"
	NotEnoughMoneyMessage paymentMessage = "not enough money on balance"
	SuccessMessage        paymentMessage = "success"
)

type OutboxEntry struct {
	ID         int64            `json:"id"`
	MessageKey int64            `json:"message_key"`
	EventType  paymentEventType `json:"event_type"`
	Message    string           `json:"message"`
	Payload    []byte           `json:"payload"`
}
