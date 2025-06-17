package storage

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"orders-service/internal/models"
)

type OrderStorage interface {
	SaveOrder(order *models.Order) error
	FindByID(id int64) (*models.Order, error)
	FindByUser(userID int64) ([]*models.Order, error)
	FetchUnprocessedOutbox() ([]models.OutboxEntry, error)
	MarkOutboxProcessed(outboxID int64) error
	UpdateStatus(orderID int64, status models.Status, message string) error
}

type Storage struct {
	db *sql.DB
}

func (p *Storage) SaveOrder(order *models.Order) (err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRow(
		`INSERT INTO orders (user_id, amount, status)
     VALUES ($1,$2,$3) RETURNING id`,
		order.UserID, order.Amount, order.Status).Scan(&order.ID)

	if err != nil {
		return err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = p.SaveEvent(tx, order.ID, string(data))

	if err != nil {
		return err
	}

	return nil
}

func (p *Storage) SaveEvent(tx *sql.Tx, aggregateId int64, payload string) error {
	_, err := tx.Exec(
		`INSERT INTO outbox (aggregate_id, payload)
     VALUES ($1,$2)`,
		aggregateId, payload,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Storage) FindByID(id int64) (*models.Order, error) {
	row := p.db.QueryRow(
		"SELECT id, user_id, amount, status,message, extract(epoch FROM created_at)::BIGINT FROM orders WHERE id=$1", id,
	)
	var o models.Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Amount, &o.Status, &o.Message, &o.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (p *Storage) FindByUser(userID int64) ([]*models.Order, error) {
	rows, err := p.db.Query(
		"SELECT id, amount, status, message,extract(epoch FROM created_at)::BIGINT FROM orders WHERE user_id = $1", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.ID, &o.UserID, &o.Amount, &o.Status, &o.Message, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrUserNotFound
	}
	return orders, nil
}

func (p *Storage) FetchUnprocessedOutbox() ([]models.OutboxEntry, error) {
	rows, err := p.db.Query(
		"SELECT id, aggregate_id, payload, created_at FROM outbox WHERE processed_at IS NULL",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.OutboxEntry
	for rows.Next() {
		var e models.OutboxEntry
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (p *Storage) MarkOutboxProcessed(outboxID int64) error {
	_, err := p.db.Exec(
		"UPDATE outbox SET processed_at=now() WHERE id=$1", outboxID,
	)
	return err
}

func (p *Storage) UpdateStatus(orderID int64, status models.Status, message string) error {
	_, err := p.db.Exec(
		`UPDATE orders SET (status, message) = ($2,$3) WHERE id = $1`, orderID, status, message,
	)
	return err
}
