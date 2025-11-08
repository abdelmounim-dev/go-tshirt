package main

import (
	"log"
	"github.com/abdelmounim-dev/go-tshirt/internal/api"
	"github.com/abdelmounim-dev/go-tshirt/internal/config"
)

func main() {
	cfg := config.Load()
	r := api.NewRouter(cfg)
	log.Printf("Server running on %s", cfg.Address)
	if err := r.Run(cfg.Address); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
