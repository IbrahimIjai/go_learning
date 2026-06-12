package helpers

import (
	"encoding/json"
	"net/http"
)

// WriteJSON sets the content type, writes the status, and encodes the body.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
