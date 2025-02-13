package handlers

import (
	"context"
	"errors"
	"fmt"
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
	fromUserID := ctx.Value("userID").(int)

	toUser, err := h.service.GetByUsername(ctx, sendCoinRequest.ToUser)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			h.badRequestResponse(w, r, fmt.Errorf("%s:%w", sendCoinRequest.ToUser, err))
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	h.logger.Print(fromUserID, toUser)
}
