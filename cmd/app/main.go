package main

import (
	"context"
	"log"
	"merch-shop/internal/config"
	"merch-shop/internal/dbinit"
	"merch-shop/internal/handlers"
	"merch-shop/internal/repository"
	"merch-shop/internal/service"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	repo := repository.NewPostgresRepository(db)
	service := service.NewService(repo, cfg)
	handler := handlers.NewHandler(service, cfg, logger)

	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Server.Port),
		Handler: handler.Routes(),
	}

	logger.Printf("starting backend on %s", srv.Addr)
	go func() {
		logger.Print(srv.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Print("shutting down backend...")
	srv.Shutdown(ctx)
}
