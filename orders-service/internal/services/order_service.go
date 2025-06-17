package services

import (
	"orders-service/internal/models"
	"orders-service/internal/storage"
)

type OrderService struct {
	store storage.OrderStorage
}

func NewOrderService(store storage.OrderStorage) *OrderService {
	return &OrderService{store: store}
}

func (s *OrderService) Create(userID int64, amount float64) (*models.Order, error) {
	o := &models.Order{UserID: userID, Amount: amount, Status: models.New}

	if err := s.store.SaveOrder(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) Get(id int64) (*models.Order, error) {
	return s.store.FindByID(id)
}

func (s *OrderService) List(userID int64) ([]*models.Order, error) {
	return s.store.FindByUser(userID)
}
