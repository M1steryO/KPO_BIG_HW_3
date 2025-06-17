package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"payments-service/internal/models"
)

type AccountStorage interface {
	CreateAccount(userID int64) (int64, error)
	Get(userID int64) (*models.Account, error)
	Deposit(userID int64, amount float64) error
	ProcessPaymentRequest(key string, payload []byte) (int64, error)
	FetchUnprocessedOutbox() ([]models.OutboxEntry, error)
	MarkOutboxProcessed(id int64) error
	MarkInboxProcessed(id int64) error
}

type PostgresAccountStorage struct{ db *sql.DB }

func (p *PostgresAccountStorage) CreateAccount(userID int64) (int64, error) {
	var id int64
	err := p.db.QueryRow(`INSERT INTO accounts(user_id,balance,updated_at) VALUES($1,0,now()) ON CONFLICT DO NOTHING RETURNING id`, userID).Scan(&id)
	return id, err
}

func (p *PostgresAccountStorage) ProcessPaymentRequest(key string, payload []byte) (inboxId int64, err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var exists bool
	if err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM inbox WHERE message_key=$1)", key).Scan(&exists); err != nil {
		return 0, err
	}
	if exists {
		return 0, tx.Commit()
	}

	var order models.OrderEntry

	if err = json.Unmarshal(payload, &order); err != nil {
		return 0, err
	}
	res, err := tx.Exec(`UPDATE accounts SET balance=balance-$1,updated_at=now() WHERE user_id=$2 AND balance>=$1`, order.Amount, order.UserID)
	if err != nil {
		return 0, err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		// проверяем причину ошибки
		res, err = tx.Exec(`SELECT id FROM accounts WHERE user_id=$1`, order.UserID)
		if err != nil {

			return 0, err
		}
		ra, _ = res.RowsAffected()
		failureReason := models.NotEnoughMoneyMessage
		if ra == 0 {
			failureReason = models.UserNotFoundMessage
		}

		if err = p.SaveOutboxEvent(tx, key, payload, models.PaymentStatusFailed, failureReason); err != nil {
			return 0, err
		}
		if inboxId, err = p.SaveInboxEvent(tx, key, payload); err != nil {
			return 0, err
		}
		return inboxId, nil
	}
	if inboxId, err = p.SaveInboxEvent(tx, key, payload); err != nil {
		return 0, err
	}

	data, _ := json.Marshal(order)
	if err = p.SaveOutboxEvent(tx, key, data, models.PaymentStatusSuccess, models.SuccessMessage); err != nil {
		return 0, err
	}

	return inboxId, nil
}

func (p *PostgresAccountStorage) SaveOutboxEvent(tx *sql.Tx, messageKey string, payload []byte, eventType, message string) error {
	if _, err := tx.Exec(`INSERT INTO payment_outbox(message_key,payload, event_type, message) VALUES($1,$2, $3, $4)`, messageKey, payload, eventType, message); err != nil {
		return err
	}
	return nil
}

func (p *PostgresAccountStorage) SaveInboxEvent(tx *sql.Tx, messageKey string, payload []byte) (int64, error) {
	var id int64
	if err := tx.QueryRow(`INSERT INTO inbox(message_key,payload) VALUES($1,$2) RETURNING id`, messageKey, payload).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PostgresAccountStorage) FetchUnprocessedOutbox() ([]models.OutboxEntry, error) {
	rows, err := p.db.Query(`SELECT id,message_key,payload, event_type, message FROM payment_outbox WHERE processed_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.OutboxEntry
	for rows.Next() {
		var e models.OutboxEntry
		if err = rows.Scan(&e.ID, &e.MessageKey, &e.Payload, &e.EventType, &e.Message); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (p *PostgresAccountStorage) MarkOutboxProcessed(id int64) error {
	_, err := p.db.Exec(`UPDATE payment_outbox SET processed_at=now() WHERE id=$1`, id)
	return err
}
func (p *PostgresAccountStorage) MarkInboxProcessed(id int64) error {
	_, err := p.db.Exec(`UPDATE inbox SET processed_at=now() WHERE id=$1`, id)
	return err
}

func (p *PostgresAccountStorage) Get(userID int64) (*models.Account, error) {
	row := p.db.QueryRow(`SELECT user_id,balance, extract(epoch FROM updated_at)::BIGINT FROM accounts WHERE user_id=$1`, userID)
	var a models.Account
	if err := row.Scan(&a.UserID, &a.Balance, &a.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (p *PostgresAccountStorage) Deposit(userID int64, amount float64) error {
	res, err := p.db.Exec(`UPDATE accounts SET balance=accounts.balance+$2,updated_at=now() WHERE user_id=$1`, userID, amount)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
