// Package httputil provides shared HTTP utility functions for consistent
// JSON encoding/decoding across all handler packages.
package httputil

import (
	"encoding/json"
	"net/http"
)

// WriteJSON encodes payload as JSON and writes it with the given status code.
// It sets the Content-Type header to application/json.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError writes a JSON error response with code and message.
// The response body has the shape: {"error": {"code": code, "message": message}}.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

// EncodeJSON encodes payload as JSON bytes.
func EncodeJSON(payload any) ([]byte, error) {
	return json.Marshal(payload)
}

// ParseJSON decodes the JSON request body into dest.
func ParseJSON(r *http.Request, dest any) error {
	return json.NewDecoder(r.Body).Decode(dest)
}
