package handlers

import (
	"context"
	"errors"
	"log"
	"merch-shop/internal/models"
	"merch-shop/internal/service"
	"merch-shop/internal/utils"
	"net/http"
	"time"
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
		h.badRequestResponse(w, r, err)
		return
	}

	if err := authRequestValid(authRequest); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token, err := h.Service.Login(ctx, authRequest.Username, authRequest.Password)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMismatchHashPassword):
			h.unauthorizedResponse(w, r, err)
		case errors.Is(err, utils.ErrTooLongPassword):
			h.badRequestResponse(w, r, err)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	authResponse := &models.AuthResponse{
		Token: token,
	}

	err = h.writeJSON(w, http.StatusOK, authResponse, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
