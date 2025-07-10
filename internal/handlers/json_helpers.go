package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	errorResponse := map[string]string{"error": message}
	WriteJSON(w, status, errorResponse)
}
