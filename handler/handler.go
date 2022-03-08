package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Resp struct {
	Status string `json:"status"`
}

// Handler struct for HTTP requests
type Handler struct {
	logger *zap.SugaredLogger
	router *mux.Router
}

// New creates a Handler struct
func New(
	logger *zap.SugaredLogger,
	router *mux.Router,
) *Handler {
	h := Handler{logger, router}
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.router.HandleFunc("/health", h.health).
		Methods("GET")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	var (
		resp = Resp{}
	)

	resp.Status = "OK"

	json.NewEncoder(w).Encode(resp)
}
