package service

import (
	"context"
	"errors"
	"fmt"
	"merch-shop/internal/config"
	"merch-shop/internal/models"
	"merch-shop/internal/repository"
	"merch-shop/internal/utils"
)

type Service struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewService(repo repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	if err := utils.ValidatePassword(password); err != nil {
		return "", err
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return s.Add(ctx, username, password)
		}
		return "", err
	}

	if err := utils.CheckPasswordHash(user.PasswordHash, password); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, s.cfg.JWT.SecretKey, s.cfg.JWT.TokenExpiry)
}

func (s *Service) Add(ctx context.Context, username, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err = s.repo.Add(ctx, user)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, s.cfg.JWT.SecretKey, s.cfg.JWT.TokenExpiry)
}

func (s *Service) SendCoin(ctx context.Context, senderID int, receiverName string, amount int) error {
	receiver, err := s.repo.GetByUsername(ctx, receiverName)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("receiver user: %w", err)
		}
		return err
	}

	if receiver.ID == senderID {
		return ErrSendToYourself
	}

	return s.repo.SendCoin(ctx, senderID, receiver.ID, amount)
}

func (s *Service) BuyItem(ctx context.Context, userID int, itemName string) error {
	return s.repo.BuyItem(ctx, userID, itemName)
}

func (s *Service) Info(ctx context.Context, userID int) (*models.InfoResponse, error) {
	coins, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.repo.GetInventory(ctx, userID)
	if err != nil {
		return nil, err
	}

	coinHistory, err := s.repo.GetCoinHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	infoResponse := &models.InfoResponse{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return infoResponse, nil
}
