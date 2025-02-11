package service

import (
	"context"
	"errors"
	"merch-shop/internal/config"
	"merch-shop/internal/models"
	"merch-shop/internal/repository"
	"merch-shop/internal/utils"
)

type UserService struct {
	userRepo repository.UserRepo
	cfg      *config.Config
}

func NewUserService(repo repository.UserRepo, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: repo,
		cfg:      cfg,
	}
}

func (u *UserService) Login(ctx context.Context, username, password string) (string, error) {
	if err := utils.ValidatePassword(password); err != nil {
		return "", err
	}

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return u.Add(ctx, username, password)
		}
		return "", err
	}

	if err := utils.CheckPasswordHash(user.PasswordHash, password); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, u.cfg.JWT.Secret, u.cfg.JWT.TokenExpiry)
}

func (u *UserService) Add(ctx context.Context, username, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err = u.userRepo.Add(ctx, user)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, u.cfg.JWT.Secret, u.cfg.JWT.TokenExpiry)
}
