package models

import "time"

type Transaction struct {
	ID         int
	SenderID   int
	ReceiverID int
	Amount     int
	CreatedAt  time.Time
}

type Inventory struct {
	ID       int
	UserID   int
	ItemID   int
	Quantity int
}

type Coins struct {
	UserID  int
	Balance int
}
