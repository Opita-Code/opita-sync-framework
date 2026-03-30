package surfaceservice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opita-sync-framework/internal/app/surfaceservice"
	"opita-sync-framework/internal/platform/memory"
)

func TestCreateIntakeTurnProducesSessionAndCandidate(t *testing.T) {
	handler := surfaceservice.NewHandler(memory.NewIntakeStore(), memory.NewProposalStore(), memory.NewEventLog())
	body, _ := json.Marshal(map[string]any{
		"tenant_id":  "tenant-1",
		"session_id": "session-1",
		"subject_id": "user-1",
		"raw_text":   "quiero cambiar una configuración",
		"trace_id":   "trace-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/intake/turns", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCreateProposalReturnsDraft(t *testing.T) {
	handler := surfaceservice.NewHandler(memory.NewIntakeStore(), memory.NewProposalStore(), memory.NewEventLog())
	body, _ := json.Marshal(map[string]any{
		"tenant_id":          "tenant-1",
		"session_id":         "session-1",
		"subject_id":         "user-1",
		"source_intent_refs": []string{"intent-1"},
		"title":              "cambio de config",
		"summary":            "resumen",
		"human_diff_ref":     "human-diff-1",
		"material_diff_ref":  "material-diff-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/proposals", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestGetIntakeProposalWorkspaceShowsNextGates(t *testing.T) {
	intakeStore := memory.NewIntakeStore()
	proposalStore := memory.NewProposalStore()
	handler := surfaceservice.NewHandler(intakeStore, proposalStore, memory.NewEventLog())

	turnBody, _ := json.Marshal(map[string]any{
		"tenant_id":  "tenant-1",
		"session_id": "session-1",
		"subject_id": "user-1",
		"raw_text":   "quiero preparar un cambio gobernado",
		"trace_id":   "trace-1",
	})
	turnReq := httptest.NewRequest(http.MethodPost, "/v1/intake/turns", bytes.NewReader(turnBody))
	turnW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(turnW, turnReq)
	if turnW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", turnW.Code, turnW.Body.String())
	}
	var turnResp map[string]any
	if err := json.Unmarshal(turnW.Body.Bytes(), &turnResp); err != nil {
		t.Fatalf("unmarshal turn response: %v", err)
	}
	intakeSession := turnResp["intake_session"].(map[string]any)
	intentCandidate := turnResp["intent_candidate"].(map[string]any)
	workspaceReq := httptest.NewRequest(http.MethodGet, "/v1/workspaces/intake-proposal?intake_session_id="+getStringField(t, intakeSession, "intake_session_id", "IntakeSessionID")+"&intent_candidate_id="+getStringField(t, intentCandidate, "intent_candidate_id", "IntentCandidateID"), nil)
	workspaceW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(workspaceW, workspaceReq)
	if workspaceW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", workspaceW.Code, workspaceW.Body.String())
	}
	var workspace map[string]any
	if err := json.Unmarshal(workspaceW.Body.Bytes(), &workspace); err != nil {
		t.Fatalf("unmarshal workspace response: %v", err)
	}
	if workspace["chat_boundary"] != "free_chat_never_applies_directly" {
		t.Fatalf("unexpected chat boundary: %#v", workspace["chat_boundary"])
	}
	nextGates := workspace["next_gates"].([]any)
	if len(nextGates) == 0 || nextGates[0] != "create_proposal_draft" {
		t.Fatalf("expected create_proposal_draft gate, got %#v", nextGates)
	}
}

func TestGetIntakeProposalWorkspaceIncludesProposalAndPatchsetSummary(t *testing.T) {
	intakeStore := memory.NewIntakeStore()
	proposalStore := memory.NewProposalStore()
	handler := surfaceservice.NewHandler(intakeStore, proposalStore, memory.NewEventLog())

	proposalBody, _ := json.Marshal(map[string]any{
		"tenant_id":          "tenant-1",
		"session_id":         "session-1",
		"subject_id":         "user-1",
		"trace_id":           "trace-1",
		"source_intent_refs": []string{"intent-1"},
		"title":              "cambio de config",
		"summary":            "resumen",
		"human_diff_ref":     "human-diff-1",
		"material_diff_ref":  "material-diff-1",
	})
	proposalReq := httptest.NewRequest(http.MethodPost, "/v1/proposals", bytes.NewReader(proposalBody))
	proposalW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(proposalW, proposalReq)
	if proposalW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", proposalW.Code, proposalW.Body.String())
	}
	var proposalResp map[string]any
	if err := json.Unmarshal(proposalW.Body.Bytes(), &proposalResp); err != nil {
		t.Fatalf("unmarshal proposal response: %v", err)
	}
	proposalID := getStringField(t, proposalResp, "proposal_draft_id", "ProposalDraftID")

	patchsetBody, _ := json.Marshal(map[string]any{
		"trace_id":                          "trace-1",
		"proposal_draft_id":                 proposalID,
		"target_artifacts":                  []string{"artifact-1"},
		"material_operations":               []string{"update config"},
		"material_diff_hash":                "diff-1",
		"policy_preview_inputs_ref":         "policy-preview-1",
		"approval_preview_inputs_ref":       "approval-preview-1",
		"classification_preview_inputs_ref": "classification-preview-1",
	})
	patchsetReq := httptest.NewRequest(http.MethodPost, "/v1/patchsets", bytes.NewReader(patchsetBody))
	patchsetW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(patchsetW, patchsetReq)
	if patchsetW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", patchsetW.Code, patchsetW.Body.String())
	}
	var patchsetResp map[string]any
	if err := json.Unmarshal(patchsetW.Body.Bytes(), &patchsetResp); err != nil {
		t.Fatalf("unmarshal patchset response: %v", err)
	}
	patchsetID := getStringField(t, patchsetResp, "patchset_candidate_id", "PatchsetCandidateID")

	workspaceReq := httptest.NewRequest(http.MethodGet, "/v1/workspaces/intake-proposal?proposal_draft_id="+proposalID+"&patchset_candidate_id="+patchsetID, nil)
	workspaceW := httptest.NewRecorder()
	handler.Routes().ServeHTTP(workspaceW, workspaceReq)
	if workspaceW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", workspaceW.Code, workspaceW.Body.String())
	}
	var workspace map[string]any
	if err := json.Unmarshal(workspaceW.Body.Bytes(), &workspace); err != nil {
		t.Fatalf("unmarshal workspace response: %v", err)
	}
	if workspace["workspace_state"] != "ready_for_preview" {
		t.Fatalf("expected ready_for_preview workspace state, got %#v", workspace["workspace_state"])
	}
	if _, ok := workspace["proposal"]; !ok {
		t.Fatalf("expected proposal summary in workspace")
	}
	if _, ok := workspace["patchset"]; !ok {
		t.Fatalf("expected patchset summary in workspace")
	}
}

func getStringField(t *testing.T, payload map[string]any, keys ...string) string {
	t.Helper()
	for _, key := range keys {
		if value, ok := payload[key]; ok && value != nil {
			if s, ok := value.(string); ok {
				return s
			}
		}
	}
	t.Fatalf("missing string field in payload, tried keys: %v payload=%#v", keys, payload)
	return ""
}
