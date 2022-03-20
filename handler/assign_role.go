package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mager/premintbot/bot"
)

type AssignRoleReq struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
}

type AssignRoleResp struct {
	Message string `json:"message"`
}

func (h *Handler) assignRole(w http.ResponseWriter, r *http.Request) {
	var (
		req  = AssignRoleReq{}
		resp = AssignRoleResp{}
	)

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Message = "Invalid JSON"
	}

	guild, err := h.bot.State.Guild(req.GuildID)
	if err != nil {
		resp.Message = "Invalid guild"
	}

	p := bot.GetConfig(context.TODO(), h.logger, h.database, guild.ID)
	err = h.bot.GuildMemberRoleAdd(guild.ID, req.UserID, p.Config.PremintRoleID)
	if err != nil {
		resp.Message = "Failed to assign role"
	} else {
		resp.Message = "Role assigned"
	}

	json.NewEncoder(w).Encode(resp)
}
