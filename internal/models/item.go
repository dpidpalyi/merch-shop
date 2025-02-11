package models

type Item struct {
	ID    int    `json:"-"`
	Name  string `json:"type"`
	Price int    `json:"-"`
}
