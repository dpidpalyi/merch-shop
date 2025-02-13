package models

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendCoinRequest struct {
	ReceiverName string `json:"toUser"`
	Amount       int    `json:"amount"`
}
