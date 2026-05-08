// Package testutil provides shared test helper functions for consistent
// test patterns across all handler packages.
package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// GetStringField extracts a string value from a map, trying multiple keys.
// It skips nil and non-string values. It calls t.Fatalf if no key matches.
func GetStringField(t *testing.T, payload map[string]any, keys ...string) string {
	t.Helper()
	for _, key := range keys {
		if value, ok := payload[key]; ok && value != nil {
			if s, ok := value.(string); ok {
				return s
			}
		}
	}
	t.Fatalf("missing string field, tried keys: %v payload=%#v", keys, payload)
	return ""
}

// PostJSON sends a POST request with a JSON-encoded body to the given handler.
func PostJSON(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("PostJSON encode body failed: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// GetJSON sends a GET request to the given handler.
func GetJSON(t *testing.T, handler http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// ErrorContains checks if err contains the expected substring.
func ErrorContains(t *testing.T, err error, expected string) bool {
	t.Helper()
	if err == nil {
		return expected == ""
	}
	return strings.Contains(err.Error(), expected)
}
