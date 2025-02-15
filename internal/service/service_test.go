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
			username:   "mao long",
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
			wantErr:  true,
			mockRepoFn: func() {
				mockRepo.On("GetByUsername", ctx, "carl").Return(nil, errors.New("db fails"))
			},
		},
		{
			name:     "user exists, password check fails",
			username: "sarah",
			password: "wrong password",
			wantErr:  true,
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

func Test_Add(t *testing.T) {
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
			name:       "failed hashing, password too long",
			username:   "mao long",
			password:   string(make([]byte, 76)),
			wantErr:    true,
			mockRepoFn: func() {},
		},
		{
			name:     "success hashing, db success",
			username: "alice",
			password: "password",
			mockRepoFn: func() {
				mockRepo.On("Add", ctx, mock.MatchedBy(func(u *models.User) bool {
					return u.Username == "alice" && len(u.PasswordHash) > 0
				})).Return(nil)
			},
		},
		{
			name:     "success hashing, db fails",
			username: "bob",
			password: "password",
			wantErr:  true,
			mockRepoFn: func() {
				mockRepo.On("Add", ctx, mock.Anything).Return(errors.New("db fails"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn()

			token, err := service.Add(ctx, tt.username, tt.password)

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

func Test_SendCoin(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	mockRepo := new(mocks.Repository)
	service := NewService(mockRepo, cfg)

	tests := []struct {
		name       string
		senderID   int
		receiverName   string
		amount int
		wantErr    bool
		mockRepoFn func()
	}{
		{
			name:     "receiver not exists",
			receiverName: "bob",
			wantErr: true,
			mockRepoFn: func() {
				mockRepo.On("GetByUsername", ctx, "bob").Return(nil, repository.ErrRecordNotFound)
			},
		},
		{
			name:     "db fails",
			receiverName: "carl",
			wantErr: true,
			mockRepoFn: func() {
				mockRepo.On("GetByUsername", ctx, "carl").Return(nil, errors.New("db fails"))
			},
		},
		{
			name:     "receiver exists, fail to send yourself",
			senderID: 1,
			receiverName: "alice",
			wantErr: true,
			mockRepoFn: func() {
				hashedPassword, _ := utils.HashPassword("password")
				mockRepo.On("GetByUsername", ctx, "alice").Return(&models.User{ID: 1, Username: "alice", PasswordHash: hashedPassword}, nil)
			},
		},
		{
			name:     "receiver exists, success to send",
			senderID: 1,
			receiverName: "sarah",
			amount: 100,
			mockRepoFn: func() {
				hashedPassword, _ := utils.HashPassword("password")
				mockRepo.On("GetByUsername", ctx, "sarah").Return(&models.User{ID: 2, Username: "sarah", PasswordHash: hashedPassword}, nil)
				mockRepo.On("SendCoin", ctx, 1, 2, 100).Return(nil)
			},
		},
		{
			name:     "receiver exists, failed to send",
			senderID: 1,
			receiverName: "sarah",
			amount: 10000,
			wantErr: true,
			mockRepoFn: func() {
				hashedPassword, _ := utils.HashPassword("password")
				mockRepo.On("GetByUsername", ctx, "sarah").Return(&models.User{ID: 2, Username: "sarah", PasswordHash: hashedPassword}, nil)
				mockRepo.On("SendCoin", ctx, 1, 2, 10000).Return(repository.ErrNotEnoughCoins)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn()

			err := service.SendCoin(ctx, tt.senderID, tt.receiverName, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func Test_BuyItem(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	mockRepo := new(mocks.Repository)
	service := NewService(mockRepo, cfg)

	tests := []struct {
		name       string
		userID   int
		itemName   string
		wantErr    bool
		mockRepoFn func()
	}{
		{
			name:     "success to buy item",
			userID: 1,
			itemName: "cup",
			mockRepoFn: func() {
				mockRepo.On("BuyItem", ctx, 1, "cup").Return(nil)
			},
		},
		{
			name:     "fails to buy item",
			userID: 2,
			itemName: "cup",
			wantErr: true,
			mockRepoFn: func() {
				mockRepo.On("BuyItem", ctx, 2, "cup").Return(repository.ErrNotEnoughCoins)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn()

			err := service.BuyItem(ctx, tt.userID, tt.itemName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func Test_Info(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	mockRepo := new(mocks.Repository)
	service := NewService(mockRepo, cfg)

	tests := []struct {
		name       string
		userID   int
		wantErr    bool
		mockRepoFn func()
	}{
		{
			name:       "user not exists",
			userID:   1,
			wantErr:    true,
			mockRepoFn: func() {
				mockRepo.On("GetBalance", ctx, 1).Return(0, repository.ErrRecordNotFound)
			},
		},
		{
			name:       "user exists, inventory fails",
			userID:   2,
			wantErr:    true,
			mockRepoFn: func() {
				mockRepo.On("GetBalance", ctx, 2).Return(0, nil)
				mockRepo.On("GetInventory", ctx, 2).Return(nil, errors.New("db fails"))
			},
		},
		{
			name:       "user exists, coinHistory fails",
			userID:   3,
			wantErr:    true,
			mockRepoFn: func() {
				mockRepo.On("GetBalance", ctx, 3).Return(0, nil)
				mockRepo.On("GetInventory", ctx, 3).Return(nil, nil)
				mockRepo.On("GetCoinHistory", ctx, 3).Return(nil, errors.New("db fails"))
			},
		},
		{
			name:       "user exists, success",
			userID:   4,
			mockRepoFn: func() {
				mockRepo.On("GetBalance", ctx, 4).Return(0, nil)
				mockRepo.On("GetInventory", ctx, 4).Return(nil, nil)
				mockRepo.On("GetCoinHistory", ctx, 4).Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFn()

			infoResponse, err := service.Info(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, infoResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, infoResponse)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
