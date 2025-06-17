package services

import (
	"payments-service/internal/models"
	"payments-service/internal/storage"
)

type AccountService struct{ store storage.AccountStorage }

func NewAccountService(store storage.AccountStorage) *AccountService {
	return &AccountService{store: store}
}
func (s *AccountService) Create(userID int64) (int64, error) {
	id, err := s.store.CreateAccount(userID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AccountService) Deposit(userID int64, amount float64) error {
	return s.store.Deposit(userID, amount)
}

func (s *AccountService) Get(userID int64) (*models.Account, error) {
	return s.store.Get(userID)
}
