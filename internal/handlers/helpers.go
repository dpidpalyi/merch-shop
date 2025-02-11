package handlers

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (h *UserHandler) readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)
}

func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	return nil
}
