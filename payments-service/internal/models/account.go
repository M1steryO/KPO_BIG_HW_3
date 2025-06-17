package models

type Account struct {
	UserID    int64   `json:"user_id"`
	Balance   float64 `json:"balance"`
	UpdatedAt int64   `json:"updated_at"`
}
