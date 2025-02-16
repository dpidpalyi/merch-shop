package models

import "time"

type User struct {
	ID           int       `json:"-"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"-"`
}
