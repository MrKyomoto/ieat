package main

import (
	"context"
	"log"

	"github.com/MrKyomoto/ieat/internal/config"
	"github.com/MrKyomoto/ieat/internal/database"
	"github.com/MrKyomoto/ieat/internal/devseed"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DataBackend == "mock" {
		log.Print("mock data enabled; database seed skipped")
		return
	}
	if cfg.Environment != "development" {
		log.Fatal("development seed is disabled unless APP_ENV=development")
	}
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := devseed.Run(ctx, pool, cfg.DevSeedPassword); err != nil {
		log.Fatal(err)
	}
	log.Print("development accounts and sample catalog initialized")
}
