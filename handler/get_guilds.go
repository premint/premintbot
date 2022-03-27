package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mager/premintbot/bot"
)

type GetGuildsResp struct {
	Guilds []bot.Guild `json:"guilds"`
}

func (h *Handler) getGuilds(w http.ResponseWriter, r *http.Request) {
	var (
		resp = GetGuildsResp{}
	)

	// Fetch all the guilds
	docs, err := h.database.Collection("guilds").Documents(r.Context()).GetAll()
	if err != nil {
		h.logger.Errorw("Failed to get guilds", "error", err)
		http.Error(w, "Failed to get guilds", http.StatusInternalServerError)
		return
	}

	// Get the guilds
	for _, doc := range docs {
		guild := bot.Guild{}
		err = doc.DataTo(&guild)
		if err != nil {
			h.logger.Errorw("Failed to get guild", "error", err)
			http.Error(w, "Failed to get guild", http.StatusInternalServerError)
			return
		}

		resp.Guilds = append(resp.Guilds, guild)
	}

	json.NewEncoder(w).Encode(resp)
}
