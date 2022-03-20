package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	var (
		resp = HealthResp{}
	)

	resp.Status = "OK"

	json.NewEncoder(w).Encode(resp)
}
