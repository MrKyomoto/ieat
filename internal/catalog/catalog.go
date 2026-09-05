package catalog

import (
	"context"
	"encoding/json"
	"net/http"
)

type Window struct {
	ID            string `json:"id"`
	ExternalCode  string `json:"externalCode"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BusinessHours string `json:"businessHours"`
}

type Floor struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Windows []Window `json:"windows"`
}

type Canteen struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Floors []Floor `json:"floors"`
}

type Store interface {
	List(ctx context.Context) ([]Canteen, error)
}

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	canteens, err := h.store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "读取食堂目录失败")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(canteens)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
