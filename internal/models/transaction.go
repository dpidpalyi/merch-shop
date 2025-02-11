package models

type Transaction struct {
	ID int
	SenderID int
	ReceiverID int
	Amount int 
	CreatedAt time.Time
}

func Inventory struct {
	ID int
	UserID int
	ItemID int
	Quantity int
}

func Coins struct {
	UserID int
	Balance int
}
