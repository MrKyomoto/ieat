package main

import (
	"context"
	"log"

	"github.com/MrKyomoto/ieat/internal/config"
	"github.com/MrKyomoto/ieat/internal/database"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DataBackend == "mock" {
		log.Print("mock data enabled; database migration skipped")
		return
	}
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := database.Migrate(ctx, pool); err != nil {
		log.Fatal(err)
	}
	log.Print("database migrations applied")
}
