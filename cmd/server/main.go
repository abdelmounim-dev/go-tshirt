package main

import (
	"log"

	"github.com/abdelmounim-dev/go-tshirt/internal/api"
	"github.com/abdelmounim-dev/go-tshirt/internal/config"
	"github.com/abdelmounim-dev/go-tshirt/internal/db"
)

func main() {
	cfg := config.Load()

	database, err := db.NewSQLite(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := api.SetupRouter(database)

	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
