package main

import (
	"log"
	"merch-shop/internal/config"
	"merch-shop/internal/dbinit"
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
	logger.Print(cfg)

	db, err := dbinit.OpenDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	db.Close()

	srv := &http.Server{
		Addr: net.JoinHostPort("", cfg.Server.Port),
	}

	logger.Printf("starting backend on %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
