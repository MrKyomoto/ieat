package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MrKyomoto/ieat/internal/auth"
	"github.com/MrKyomoto/ieat/internal/catalog"
	"github.com/MrKyomoto/ieat/internal/config"
	"github.com/MrKyomoto/ieat/internal/database"
	"github.com/MrKyomoto/ieat/internal/httpapi"
	"github.com/MrKyomoto/ieat/internal/mockdata"
)

func main() {
	log.SetOutput(os.Stdout)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o750); err != nil {
		log.Fatalf("create upload directory: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var dependencies httpapi.Dependencies
	var closeStore = func() {}
	switch cfg.DataBackend {
	case "mock":
		store, err := mockdata.New(cfg.DevSeedPassword)
		if err != nil {
			log.Fatal(err)
		}
		dependencies.Auth = store
		dependencies.Catalog = store
	case "postgres":
		pool, err := database.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
		closeStore = pool.Close
		dependencies.Auth = auth.NewPostgresStore(pool)
		dependencies.Catalog = catalog.NewPostgresStore(pool)
		dependencies.Ready = pool.Ping
	}
	defer closeStore()

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.NewRouter(cfg, dependencies),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown server: %v", err)
		}
	}()

	log.Printf("api listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
