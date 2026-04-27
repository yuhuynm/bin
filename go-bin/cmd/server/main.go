package main

import (
	"log"

	"go-bin/internal/app"
	"go-bin/internal/config"
	"go-bin/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	server := app.NewServer(cfg, db)
	if err := server.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
