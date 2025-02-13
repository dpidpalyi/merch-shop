package handlers

import (
	"encoding/json"
	"errors"
	"merch-shop/internal/models"
	"net/http"
)

var (
	ErrEmptyNamePassword    = errors.New("empty name or password specified")
	ErrEmptyToUser          = errors.New("empty toUser field")
	ErrZeroOrNegativeAmount = errors.New("amount to send should be positive")
)

func (h *Handler) readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func authRequestValid(authRequest *models.AuthRequest) error {
	if authRequest.Username == "" || authRequest.Password == "" {
		return ErrEmptyNamePassword
	}

	return nil
}

func sendCoinRequestValid(sendCoinRequest *models.SendCoinRequest) error {
	if sendCoinRequest.ToUser == "" {
		return ErrEmptyToUser
	}

	if sendCoinRequest.Amount <= 0 {
		return ErrZeroOrNegativeAmount
	}

	return nil
}
