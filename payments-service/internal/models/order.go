package models

type OrderEntry struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}
