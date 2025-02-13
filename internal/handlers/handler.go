package handlers

import (
	"context"
	"errors"
	"log"
	"merch-shop/internal/config"
	"merch-shop/internal/models"
	"merch-shop/internal/repository"
	"merch-shop/internal/service"
	"merch-shop/internal/utils"
	"net/http"
	"time"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
	logger  *log.Logger
}

func NewHandler(service *service.Service, cfg *config.Config, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
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

	token, err := h.service.Login(ctx, authRequest.Username, authRequest.Password)
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

func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	sendCoinRequest := &models.SendCoinRequest{}
	err := h.readJSON(r, sendCoinRequest)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	if err := sendCoinRequestValid(sendCoinRequest); err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	senderID := ctx.Value("userID").(int)

	err = h.service.SendCoin(ctx, senderID, sendCoinRequest.ReceiverName, sendCoinRequest.Amount)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound),
			errors.Is(err, service.ErrNotEnoughCoins),
			errors.Is(err, service.ErrSendToYourself):
			h.badRequestResponse(w, r, err)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, "GOOD", nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
