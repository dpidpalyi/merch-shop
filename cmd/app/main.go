package main

import (
	"log"
	"merch-shop/internal/config"
	"merch-shop/internal/dbinit"
	"merch-shop/internal/handlers"
	"merch-shop/internal/repository"
	"merch-shop/internal/service"
	"net"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	cfg, err := config.New(".")
	if err != nil {
		logger.Fatal(err)
	}

	db, err := dbinit.OpenDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	userRepo := repository.NewPostgresUserRepo(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userService, cfg, logger)

	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Server.Port),
		Handler: userHandler.Routes(),
	}

	logger.Printf("starting backend on %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
