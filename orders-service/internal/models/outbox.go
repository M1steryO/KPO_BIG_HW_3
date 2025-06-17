package models

import "time"

type OutboxEntry struct {
	ID          int64      `json:"id"`
	AggregateID int64      `json:"aggregate_id"`
	Payload     []byte     `json:"payload"`
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`
}
