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

func (s *Service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) SendCoin(ctx context.Context, senderID int, receiverName string, amount int) error {
	senderBalance, err := s.repo.GetBalance(ctx, senderID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return fmt.Errorf("sender user: %w", err)
		}
		return err
	}

	if err := checkBalance(senderBalance, amount); err != nil {
		return err
	}

	receiver, err := s.repo.GetByUsername(ctx, receiverName)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return fmt.Errorf("receiver user: %w", err)
		}
		return err
	}

	if receiver.ID == senderID {
		return ErrSendToYourself
	}

	err = s.repo.SendCoin(ctx, senderID, receiver.ID, amount)
	if err != nil {
		return err
	}

	return nil
}
