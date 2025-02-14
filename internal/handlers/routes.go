package handlers

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", h.Auth)
	mux.HandleFunc("GET /api/info", h.MiddlewareAuth(h.Info))
	mux.HandleFunc("POST /api/sendCoin", h.MiddlewareAuth(h.SendCoin))
	mux.HandleFunc("GET /api/buy/{item}", h.MiddlewareAuth(h.BuyItem))

	return mux
}
