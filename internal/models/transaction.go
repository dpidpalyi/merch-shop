package models

import "time"

type Transaction struct {
	ID         int
	SenderID   int
	ReceiverID int
	Amount     int
	CreatedAt  time.Time
}

type Coins struct {
	UserID  int
	Balance int
}
