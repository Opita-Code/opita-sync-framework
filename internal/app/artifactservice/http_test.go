package artifactservice_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"opita-sync-framework/internal/app/artifactservice"
	"opita-sync-framework/internal/platform/filesystem"
	"opita-sync-framework/internal/platform/memory"
)

func TestPutAndGetArtifact(t *testing.T) {
	root := filepath.Join(t.TempDir(), "artifacts")
	store, err := filesystem.NewArtifactStore(root)
	if err != nil {
		t.Fatalf("artifact store init failed: %v", err)
	}
	h := artifactservice.NewHandler(store, memory.NewRetrievalStore())
	body, _ := json.Marshal(map[string]any{
		"artifact_ref":         "artifact-demo-1",
		"tenant_id":            "tenant-1",
		"kind":                 "demo",
		"classification_level": "internal",
		"content_type":         "text/plain",
		"body_base64":          base64.StdEncoding.EncodeToString([]byte("hello artifact")),
		"title":                "demo artifact",
		"index_text":           "hello artifact",
	})
	putReq := httptest.NewRequest(http.MethodPost, "/v1/artifacts", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, putReq)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/v1/artifacts/artifact-demo-1", nil)
	w2 := httptest.NewRecorder()
	h.Routes().ServeHTTP(w2, getReq)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("expected artifact root to exist: %v", err)
	}
}

func TestRetrievalSearchReturnsMatch(t *testing.T) {
	root := filepath.Join(t.TempDir(), "artifacts")
	store, err := filesystem.NewArtifactStore(root)
	if err != nil {
		t.Fatalf("artifact store init failed: %v", err)
	}
	retrievalStore := memory.NewRetrievalStore()
	h := artifactservice.NewHandler(store, retrievalStore)
	body, _ := json.Marshal(map[string]any{
		"artifact_ref":         "artifact-demo-2",
		"tenant_id":            "tenant-1",
		"kind":                 "demo",
		"classification_level": "internal",
		"content_type":         "text/plain",
		"body_base64":          base64.StdEncoding.EncodeToString([]byte("search me")),
		"title":                "searchable",
		"index_text":           "search me please",
	})
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/v1/artifacts", bytes.NewReader(body)))
	searchBody, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "text": "search", "limit": 10})
	searchReq := httptest.NewRequest(http.MethodPost, "/v1/retrieval/search", bytes.NewReader(searchBody))
	w2 := httptest.NewRecorder()
	h.Routes().ServeHTTP(w2, searchReq)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}
