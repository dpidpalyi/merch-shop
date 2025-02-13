package handlers

import (
	"merch-shop/internal/models"
	"net/http"
)

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	data := models.ErrorResponse{Errors: message}

	err := h.writeJSON(w, status, data, nil)
	if err != nil {
		h.logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Print(err)
	message := "internal server error"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusUnauthorized, err.Error())
}
