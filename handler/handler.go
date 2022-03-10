package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HealthResp struct {
	Status string `json:"status"`
}

type DebugReq struct {
}

type DebugResp struct {
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
	h.router.HandleFunc("/debug", h.health).
		Methods("POST")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	var (
		resp = HealthResp{}
	)

	resp.Status = "OK"

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) debug(w http.ResponseWriter, r *http.Request) {
	var (
		resp = DebugResp{}
	)

	resp.Status = "OK"

	json.NewEncoder(w).Encode(resp)
}
