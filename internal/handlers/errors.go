package handlers

import (
	"merch-shop/internal/models"
	"net/http"
)

func (h *UserHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	data := models.ErrorResponse{Errors: message}

	err := h.writeJSON(w, status, data, nil)
	if err != nil {
		h.Logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Print(err)
	message := "internal server error"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *UserHandler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *UserHandler) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusUnauthorized, err.Error())
}
