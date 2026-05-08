package testutil_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"opita-sync-framework/internal/testutil"
)

func TestGetStringFieldReturnsValue(t *testing.T) {
	payload := map[string]any{"name": "test-value"}
	got := testutil.GetStringField(t, payload, "name")
	if got != "test-value" {
		t.Fatalf("expected 'test-value', got %q", got)
	}
}

func TestGetStringFieldWithFallbackKeys(t *testing.T) {
	payload := map[string]any{"Name": "fallback", "name": "exact"}
	got := testutil.GetStringField(t, payload, "Name", "name")
	if got != "fallback" {
		t.Fatalf("expected first match 'fallback', got %q", got)
	}
}

func TestGetStringFieldSkipsNilValues(t *testing.T) {
	payload := map[string]any{"missing": nil, "actual": "real-value"}
	got := testutil.GetStringField(t, payload, "missing", "actual")
	if got != "real-value" {
		t.Fatalf("expected 'real-value', got %q", got)
	}
}

func TestGetStringFieldSkipsNonStringValues(t *testing.T) {
	payload := map[string]any{"count": 42, "name": "string-value"}
	got := testutil.GetStringField(t, payload, "count", "name")
	if got != "string-value" {
		t.Fatalf("expected 'string-value' after skipping int, got %q", got)
	}
}

func TestPostJSONExecutesRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %q", contentType)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if body["key"] != "value" {
			t.Fatalf("expected key=value, got %#v", body["key"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	input := map[string]string{"key": "value"}
	w := testutil.PostJSON(t, handler, "POST", "/test", input)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetJSONExecutesRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/test" {
			t.Fatalf("expected /api/test, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"found"}`))
	})
	w := testutil.GetJSON(t, handler, "/api/test")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["result"] != "found" {
		t.Fatalf("expected result=found, got %#v", resp["result"])
	}
}

func TestErrorContainsMatchesSubstring(t *testing.T) {
	err := errors.New("this is a test error: something went wrong")
	if !testutil.ErrorContains(t, err, "test error") {
		t.Fatal("ErrorContains should have found 'test error'")
	}
}

func TestErrorContainsNoMatch(t *testing.T) {
	err := errors.New("unrelated error")
	if testutil.ErrorContains(t, err, "test error") {
		t.Fatal("ErrorContains should NOT have found 'test error'")
	}
}

func TestErrorContainsWithNilError(t *testing.T) {
	if testutil.ErrorContains(t, nil, "anything") {
		t.Fatal("ErrorContains should return false for nil error")
	}
}

func TestErrorContainsEmptyExpected(t *testing.T) {
	err := errors.New("some error")
	if !testutil.ErrorContains(t, err, "") {
		t.Fatal("ErrorContains should return true for empty expected string")
	}
}
