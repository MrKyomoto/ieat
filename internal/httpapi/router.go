package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MrKyomoto/ieat/internal/auth"
	"github.com/MrKyomoto/ieat/internal/catalog"
	"github.com/MrKyomoto/ieat/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	Auth    auth.Store
	Catalog catalog.Store
	Ready   func(context.Context) error
}

func NewRouter(cfg config.Config, dependencies Dependencies) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	authHandler := auth.NewHandler(dependencies.Auth, cfg)
	catalogHandler := catalog.NewHandler(dependencies.Catalog)

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if dependencies.Ready != nil {
			if err := dependencies.Ready(r.Context()); err != nil {
				writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
				return
			}
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	router.Route("/api/v1", func(api chi.Router) {
		api.Use(requireTrustedOrigin(cfg.WebOrigin))
		api.Post("/auth/session", authHandler.Login)

		api.Group(func(protected chi.Router) {
			protected.Use(authHandler.Require)
			protected.Get("/auth/me", authHandler.Me)
			protected.Delete("/auth/session", authHandler.Logout)
			protected.Get("/catalog/canteens", catalogHandler.List)
		})
	})

	return router
}

func requireTrustedOrigin(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				if origin := r.Header.Get("Origin"); origin != "" && origin != expected {
					writeJSON(w, http.StatusForbidden, map[string]string{"error": "请求来源不受信任"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
