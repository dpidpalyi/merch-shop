package handlers

import "net/http"

func (h *UserHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", h.Auth)

	return mux
}
