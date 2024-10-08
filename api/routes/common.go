package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	InternalServerError = "Internal server error, err msg: %s"
)

// jsonResponse ... Marshals and writes a JSON response provided arbitrary data
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf(InternalServerError, err.Error()), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, fmt.Sprintf(InternalServerError, err.Error()), http.StatusInternalServerError)
		return err
	}

	return nil
}
