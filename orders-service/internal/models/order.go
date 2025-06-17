package models

type Status string

const (
	New       Status = "new"
	Finished  Status = "finished"
	Cancelled Status = "cancelled"
)

type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    Status  `json:"status"`
	Message   string  `json:"message"`
	CreatedAt int64   `json:"created_at"`
}
