package handler

import (
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
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
	bot      *discordgo.Session
	database *firestore.Client
	logger   *zap.SugaredLogger
	router   *mux.Router
}

// New creates a Handler struct
func New(
	logger *zap.SugaredLogger,
	router *mux.Router,
	database *firestore.Client,
	bot *discordgo.Session,
) *Handler {
	h := Handler{bot, database, logger, router}
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.router.HandleFunc("/health", h.health).
		Methods("GET")
	// h.router.HandleFunc("/guilds", h.getGuilds).
	// 	Methods("GET")
	h.router.HandleFunc("/debug", h.debug).
		Methods("POST")
}
