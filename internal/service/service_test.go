package service

import (
	"context"
	"errors"
	"merch-shop/internal/config"
	"merch-shop/internal/models"
	"merch-shop/internal/repository"
	"merch-shop/internal/repository/mocks"
	"merch-shop/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Login(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	mockRepo := new(mocks.Repository)
	service := NewService(mockRepo, cfg)

	tests := []struct {
		name       string
		username   string
		password   string
		wantErr    bool
		mockRepoFn func()
	}{
		{
			name:       "password too long",
			username:   "ma long",
			password:   string(make([]byte, 76)),
			wantErr:    true,
			mockRepoFn: func() {},
		},
		{
			name:     "user exists, password check success",
			username: "bob",
			password: "password",
			mockRepoFn: func() {
				hashedPassword, _ := utils.HashPassword("password")
				mockRepo.On("GetByUsername", ctx, "bob").Return(&models.User{Username: "bob", PasswordHash: hashedPassword}, nil)
			},
		},
		{
			name:     "user not exists, db success",
			username: "alice",
			password: "password",
			mockRepoFn: func() {
				mockRepo.On("GetByUsername", ctx, "alice").Return(nil, repository.ErrRecordNotFound)
				mockRepo.On("Add", ctx, mock.MatchedBy(func(u *models.User) bool {
					return u.Username == "alice" && len(u.PasswordHash) > 0
				})).Return(nil)
			},
		},
		{
			name:     "user not exists, db fails",
			username: "carl",
			password: "password",
			wantErr: true,
			mockRepoFn: func() {
				mockRepo.On("GetByUsername", ctx, "carl").Return(nil, errors.New("db fails"))
			},
		},
		{
			name:     "user exists, password check fails",
			username: "sarah",
			password: "wrong password",
			wantErr: true,
			mockRepoFn: func() {
				hashedPassword, _ := utils.HashPassword("password")
				mockRepo.On("GetByUsername", ctx, "sarah").Return(&models.User{Username: "sarah", PasswordHash: hashedPassword}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn()

			token, err := service.Login(ctx, tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
