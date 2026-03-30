package previewservice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opita-sync-framework/internal/app/previewservice"
	"opita-sync-framework/internal/engine/preview"
	"opita-sync-framework/internal/engine/simulation"
	"opita-sync-framework/internal/platform/memory"
)

func TestCreatePreviewReturnsCandidateAndSimulations(t *testing.T) {
	store := memory.NewPreviewStore()
	handler := previewservice.NewHandler(store, simulation.NewService(memory.NewPolicyEngine()), memory.NewEventLog())

	body, _ := json.Marshal(map[string]any{
		"tenant_id":          "tenant-1",
		"session_id":         "session-1",
		"subject_id":         "user-1",
		"proposal_draft_id":  "proposal-1",
		"contract_id":        "contract-1",
		"execution_id":       "exec-1",
		"patchset_ref":       "patchset-1",
		"human_diff_ref":     "human-diff-1",
		"material_diff_ref":  "material-diff-1",
		"material_diff_hash": "hash-1",
		"preview_scope":      "default",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/previews", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Routes().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if _, ok := resp["preview_candidate"]; !ok {
		t.Fatalf("expected preview_candidate in response")
	}
	if _, ok := resp["simulation_results"]; !ok {
		t.Fatalf("expected simulation_results in response")
	}
}

func TestGetReadablePreviewReturnsDiffRiskAndApprovals(t *testing.T) {
	store := memory.NewPreviewStore()
	handler := previewservice.NewHandler(store, simulation.NewService(memory.NewPolicyEngine()), memory.NewEventLog())
	body, _ := json.Marshal(map[string]any{
		"tenant_id":          "tenant-1",
		"session_id":         "session-1",
		"subject_id":         "user-1",
		"proposal_draft_id":  "proposal-1",
		"contract_id":        "contract-1",
		"execution_id":       "exec-1",
		"patchset_ref":       "patchset-1",
		"human_diff_ref":     "human-diff-1",
		"material_diff_ref":  "material-diff-1",
		"material_diff_hash": "hash-1",
		"preview_scope":      "default",
	})
	createReq := httptest.NewRequest(http.MethodPost, "/v1/previews", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", createW.Code, createW.Body.String())
	}
	var createResp map[string]any
	if err := json.Unmarshal(createW.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}
	previewCandidate := createResp["preview_candidate"].(map[string]any)
	previewID := getStringField(t, previewCandidate, "preview_candidate_id", "PreviewCandidateID")
	readableReq := httptest.NewRequest(http.MethodGet, "/v1/readable-previews/"+previewID, nil)
	readableW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(readableW, readableReq)
	if readableW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", readableW.Code, readableW.Body.String())
	}
	var readable map[string]any
	if err := json.Unmarshal(readableW.Body.Bytes(), &readable); err != nil {
		t.Fatalf("unmarshal readable preview: %v", err)
	}
	if readable["prediction_boundary"] != "simulation_preview_not_kernel_truth" {
		t.Fatalf("unexpected prediction boundary: %#v", readable["prediction_boundary"])
	}
	if _, ok := readable["diff"]; !ok {
		t.Fatalf("expected diff section")
	}
	if _, ok := readable["approvals"]; !ok {
		t.Fatalf("expected approvals section")
	}
	if _, ok := readable["risk"]; !ok {
		t.Fatalf("expected risk section")
	}
}

func TestListSimulationsRequiresPreviewID(t *testing.T) {
	store := memory.NewPreviewStore()
	handler := previewservice.NewHandler(store, simulation.NewService(memory.NewPolicyEngine()), memory.NewEventLog())
	req := httptest.NewRequest(http.MethodGet, "/v1/simulations", nil)
	w := httptest.NewRecorder()

	handler.Routes().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

var _ preview.Service = (*memory.PreviewStore)(nil)

func getStringField(t *testing.T, payload map[string]any, keys ...string) string {
	t.Helper()
	for _, key := range keys {
		if value, ok := payload[key]; ok && value != nil {
			if s, ok := value.(string); ok {
				return s
			}
		}
	}
	t.Fatalf("missing string field, keys=%v payload=%#v", keys, payload)
	return ""
}
