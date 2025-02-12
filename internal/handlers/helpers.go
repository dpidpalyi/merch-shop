package handlers

import (
	"encoding/json"
	"errors"
	"merch-shop/internal/models"
	"net/http"
)

var (
	ErrEmptyNamePassword = errors.New("empty name or password specified")
)

func (h *UserHandler) readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)
}

func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
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
