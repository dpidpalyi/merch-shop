package handlers

import "net/http"

func (h *UserHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", h.Auth)
	mux.HandleFunc("POST /api/sendCoin", h.MiddlewareAuth(h.SendCoin))

	return mux
}
