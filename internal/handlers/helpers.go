package handlers

import (
	"encoding/json"
	"merch-shop/internal/models"
	"net/http"
)

func (h *Handler) readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

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
	if sendCoinRequest.ReceiverName == "" {
		return ErrEmptyToUser
	}

	if sendCoinRequest.Amount <= 0 {
		return ErrZeroOrNegativeAmount
	}

	return nil
}
