package models

import "time"

type User struct {
	ID           int       `json:"-"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"-"`
}

type AuthToken struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
}

type AuthRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
