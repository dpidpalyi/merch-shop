package repository

import (
	"context"
	"errors"
	"merch-shop/internal/models"
)

var (
	ErrDBError = errors.New("db error")
)

type MockRepository struct{}

func (m *MockRepository) Add(ctx context.Context, user *models.User) error {
	if user.Username == "wrong" {
		return ErrDBError
	}
	return nil
}

func (m *MockRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (m *MockRepository) BuyItem(ctx context.Context, userID int, itemName string) error {
	return nil
}

func (m *MockRepository) SendCoin(ctx context.Context, senderID, receiverID int, amount int) error {
	return nil
}
func (m *MockRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

func (m *MockRepository) GetInventory(ctx context.Context, userID int) ([]*models.InventoryItem, error) {
	return nil, nil
}

func (m *MockRepository) GetCoinHistory(ctx context.Context, userID int) (*models.CoinHistory, error) {
	return nil, nil
}
