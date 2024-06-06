package routes

import (
	"encoding/json"
	"net/http"
)

type RetrieveRequest struct {
	DAType int         `json:"da_type"`
	Args   interface{} `json:"args"`
}

// RetrieveWithTypePathHandler ... Handles /api/v1/retrieve-with-type Post requests
func (h Routes) RetrieveWithTypePathHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req RetrieveRequest
	err := decoder.Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode retrieve request", "err", err)
		return
	}

	res, err := h.svc.RetrieveFromDAWithType(req.DAType, req.Args)
	if err != nil {
		http.Error(w, "Internal server error rollup with type", http.StatusInternalServerError)
		h.logger.Error("Unable to rollup with type", "err", err.Error())
		return
	}

	err = jsonResponse(w, res, http.StatusOK)
	if err != nil {
		h.logger.Error("Error writing response", "err", err.Error())
	}
}
