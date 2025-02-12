package models

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
