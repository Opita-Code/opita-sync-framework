package httputil_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"opita-sync-framework/internal/httputil"
)

func TestWriteJSONSetsContentType(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}

func TestWriteJSONWritesStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteJSON(w, http.StatusCreated, map[string]string{"id": "abc"})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestWriteJSONEncodesPayload(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]any{"name": "test", "count": 42}
	httputil.WriteJSON(w, http.StatusOK, payload)
	var decoded map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded["name"] != "test" {
		t.Fatalf("expected name=test, got %#v", decoded["name"])
	}
	if decoded["count"].(float64) != 42 {
		t.Fatalf("expected count=42, got %#v", decoded["count"])
	}
}

func TestWriteJSONWithNilPayload(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteJSON(w, http.StatusNoContent, nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestWriteErrorSetsContentType(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteError(w, http.StatusBadRequest, "test.code", "test message")
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}

func TestWriteErrorWritesStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestWriteErrorEncodesErrorStructure(t *testing.T) {
	w := httptest.NewRecorder()
	httputil.WriteError(w, http.StatusBadRequest, "request.invalid", "invalid input")

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	errObj, ok := resp["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object, got %#v", resp["error"])
	}
	if errObj["code"] != "request.invalid" {
		t.Fatalf("expected code=request.invalid, got %#v", errObj["code"])
	}
	if errObj["message"] != "invalid input" {
		t.Fatalf("expected message=invalid input, got %#v", errObj["message"])
	}
}

func TestEncodeJSONProducesValidJSON(t *testing.T) {
	payload := map[string]any{"key": "value", "num": 3.14}
	data, err := httputil.EncodeJSON(payload)
	if err != nil {
		t.Fatalf("encodeJSON failed: %v", err)
	}
	if !json.Valid(data) {
		t.Fatal("encodeJSON produced invalid JSON")
	}
}

func TestEncodeJSONProducedExpectedOutput(t *testing.T) {
	payload := map[string]string{"hello": "world"}
	data, err := httputil.EncodeJSON(payload)
	if err != nil {
		t.Fatalf("encodeJSON failed: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded["hello"] != "world" {
		t.Fatalf("expected hello=world, got %#v", decoded["hello"])
	}
}

func TestParseJSONDecodesValidBody(t *testing.T) {
	body := `{"name":"test","value":42}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var dest struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	if err := httputil.ParseJSON(req, &dest); err != nil {
		t.Fatalf("parseJSON failed: %v", err)
	}
	if dest.Name != "test" {
		t.Fatalf("expected name=test, got %q", dest.Name)
	}
	if dest.Value != 42 {
		t.Fatalf("expected value=42, got %d", dest.Value)
	}
}

func TestParseJSONReturnsErrorOnInvalidBody(t *testing.T) {
	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var dest map[string]any
	if err := httputil.ParseJSON(req, &dest); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseJSONEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")

	var dest map[string]any
	if err := httputil.ParseJSON(req, &dest); err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
}

func TestParseJSONNullBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("null"))
	req.Header.Set("Content-Type", "application/json")

	var dest map[string]any
	if err := httputil.ParseJSON(req, &dest); err != nil {
		t.Fatalf("parseJSON('null') failed: %v", err)
	}
	if dest != nil {
		t.Fatalf("expected nil map for 'null' body, got %#v", dest)
	}
}
