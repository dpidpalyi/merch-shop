package handlers

import (
	"log"
	"merch-shop/internal/models"
	"merch-shop/internal/service"
	"net/http"
)

type UserHandler struct {
	Service *service.UserService
	Logger  *log.Logger
}

func NewUserHandler(service *service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	authRequest := &models.AuthRequest{}
	err := h.readJSON(r, authRequest)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	h.Logger.Print(authRequest)
}
