package models

type User struct {
	ID int `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"-"`
}

type AuthToken struct {
	ID int `json:"-"`
	UserID int `json:"-"`
	Token string `json:"token"`
	CreatedAt time.Time `json:"-"`
}
