package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/premint/premintbot/bot"
)

type AssignRoleReq struct {
	APIKey  string `json:"api_key"`
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
}

type AssignRoleResp struct {
	Message string `json:"message"`
}

func (h *Handler) assignRole(w http.ResponseWriter, r *http.Request) {
	var (
		req     = AssignRoleReq{}
		resp    = AssignRoleResp{}
		ctx     = context.TODO()
		p       = &bot.ConfigParams{}
		err     error
		guildID string
	)

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Message = "Invalid JSON"
	}

	if req.GuildID != "" {
		guild, err := h.bot.State.Guild(req.GuildID)
		if err != nil {
			resp.Message = "Invalid guild"
		}
		guildID = guild.ID
		p = bot.GetConfig(ctx, h.logger, h.database, guildID)
	} else {
		p = bot.GetConfigWithAPIKey(ctx, h.logger, h.database, req.APIKey)
		guildID = p.Config.GuildID
	}

	err = h.bot.GuildMemberRoleAdd(guildID, req.UserID, p.Config.PremintRoleID)
	if err != nil {
		resp.Message = "Failed to assign role"
	} else {
		resp.Message = "Role assigned"
	}

	json.NewEncoder(w).Encode(resp)
}
