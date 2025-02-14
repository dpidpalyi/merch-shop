package service

import (
	"context"
	"errors"
	"fmt"
	"merch-shop/internal/config"
	"merch-shop/internal/repository"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func Test_Add(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		wantErr  error
	}{
		{
			name:     "success",
			user:     "valid",
			password: "password",
			wantErr:  nil,
		},
		{
			name:     "failure due DB error",
			user:     "wrong",
			password: "password",
			wantErr:  repository.ErrDBError,
		},
		{
			name:     "failure due very long password > 72 bytes",
			user:     "valid",
			password: string(make([]byte, 74)),
			wantErr:  bcrypt.ErrPasswordTooLong,
		},
	}

	repo := &repository.MockRepository{}
	s := NewService(repo, &config.Config{})

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %v", tt.name), func(t *testing.T) {
			_, gotErr := s.Add(context.Background(), tt.user, tt.password)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("got err: %v, want: %v", gotErr, tt.wantErr)
			}
		})
	}
}
