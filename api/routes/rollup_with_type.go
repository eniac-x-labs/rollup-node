package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type RollupRequest struct {
	DAType int    `json:"da_type"`
	Data   string `json:"data"`
}

// RollupWithTypePathHandler ... Handles /api/v1/rollup-with-type Post requests
func (h Routes) RollupWithTypePathHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req RollupRequest
	err := decoder.Decode(&req)
	dataB, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request date, want base64. Err msg: %s", err.Error()), http.StatusBadRequest)
		h.logger.Error("failed to decode rollup request", "err", err)
		return
	}

	res, err := h.svc.RollupWithType(dataB, req.DAType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error rollup with type, err msg: %s", err.Error()), http.StatusInternalServerError)
		h.logger.Error("Unable to rollup with type", "err", err.Error())
		return
	}

	err = jsonResponse(w, res, http.StatusOK)
	if err != nil {
		h.logger.Error("Error writing response", "err", err.Error())
	}
}
