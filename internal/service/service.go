package service

import (
	"context"
	"errors"
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

func (u *Service) Login(ctx context.Context, username, password string) (string, error) {
	if err := utils.ValidatePassword(password); err != nil {
		return "", err
	}

	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return u.Add(ctx, username, password)
		}
		return "", err
	}

	if err := utils.CheckPasswordHash(user.PasswordHash, password); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, u.cfg.JWT.SecretKey, u.cfg.JWT.TokenExpiry)
}

func (u *Service) Add(ctx context.Context, username, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err = u.repo.Add(ctx, user)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, u.cfg.JWT.SecretKey, u.cfg.JWT.TokenExpiry)
}

func (u *Service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.repo.GetByUsername(ctx, username)
}
