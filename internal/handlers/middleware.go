package handlers

import (
	"context"
	"merch-shop/internal/utils"
	"net/http"
)

func (h *UserHandler) MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := utils.ExtractTokenFromHeader(r)
		if err != nil {
			h.unauthorizedResponse(w, r, err)
			return
		}

		userID, err := utils.ValidateToken(tokenStr, h.cfg.JWT.SecretKey)
		if err != nil {
			h.unauthorizedResponse(w, r, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)

		next(w, r.WithContext(ctx))
	}
}
