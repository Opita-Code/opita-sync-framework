package artifactservice

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/artifacts/storage"
	"opita-sync-framework/internal/retrieval"
)

type Handler struct {
	Artifacts storage.Service
	Retrieval retrieval.Service
}

type putArtifactRequest struct {
	ArtifactRef         string `json:"artifact_ref"`
	TenantID            string `json:"tenant_id"`
	Kind                string `json:"kind"`
	ClassificationLevel string `json:"classification_level"`
	ContentType         string `json:"content_type"`
	BodyBase64          string `json:"body_base64"`
	Title               string `json:"title,omitempty"`
	IndexText           string `json:"index_text,omitempty"`
}

type searchRequest struct {
	TenantID string `json:"tenant_id"`
	Text     string `json:"text"`
	Limit    int    `json:"limit"`
}

func NewHandler(artifacts storage.Service, retrievalService retrieval.Service) *Handler {
	return &Handler{Artifacts: artifacts, Retrieval: retrievalService}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/artifacts", h.handlePutArtifact)
	mux.HandleFunc("GET /v1/artifacts/", h.handleGetArtifact)
	mux.HandleFunc("POST /v1/retrieval/search", h.handleSearch)
	return mux
}

func (h *Handler) handlePutArtifact(w http.ResponseWriter, r *http.Request) {
	var req putArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "artifact.invalid_json", "message": err.Error()})
		return
	}
	body, err := base64.StdEncoding.DecodeString(req.BodyBase64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "artifact.invalid_body_base64", "message": err.Error()})
		return
	}
	artifact, err := h.Artifacts.Put(storage.PutRequest{
		Artifact: storage.Artifact{
			ArtifactRef:         req.ArtifactRef,
			TenantID:            req.TenantID,
			Kind:                req.Kind,
			ClassificationLevel: req.ClassificationLevel,
			ContentType:         req.ContentType,
			CreatedAt:           time.Now().UTC(),
		},
		Body: body,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "artifact.put_failed", "message": err.Error()})
		return
	}
	if strings.TrimSpace(req.IndexText) != "" {
		_ = h.Retrieval.Index(retrieval.Document{
			DocumentRef:         "doc-" + artifact.ArtifactRef,
			TenantID:            req.TenantID,
			ArtifactRef:         artifact.ArtifactRef,
			ClassificationLevel: req.ClassificationLevel,
			Title:               req.Title,
			Text:                req.IndexText,
			Tags:                []string{req.Kind},
		})
	}
	writeJSON(w, http.StatusCreated, artifact)
}

func (h *Handler) handleGetArtifact(w http.ResponseWriter, r *http.Request) {
	artifactRef := strings.TrimPrefix(r.URL.Path, "/v1/artifacts/")
	if artifactRef == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "artifact.missing_ref"})
		return
	}
	artifact, found, err := h.Artifacts.Get(artifactRef)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "artifact.get_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "artifact.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"artifact":    artifact.Artifact,
		"body_base64": base64.StdEncoding.EncodeToString(artifact.Body),
	})
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	var req searchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "retrieval.invalid_json", "message": err.Error()})
		return
	}
	results, err := h.Retrieval.Search(retrieval.Query{TenantID: req.TenantID, Text: req.Text, Limit: req.Limit})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "retrieval.search_failed", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
