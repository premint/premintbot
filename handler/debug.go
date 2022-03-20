package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) debug(w http.ResponseWriter, r *http.Request) {
	var (
		resp = DebugResp{}
	)

	resp.Status = "OK"

	json.NewEncoder(w).Encode(resp)
}
