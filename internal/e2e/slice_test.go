package e2e_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"opita-sync-framework/internal/app/artifactservice"
	"opita-sync-framework/internal/app/devsurface"
	"opita-sync-framework/internal/app/intentservice"
	"opita-sync-framework/internal/app/operatorsurface"
	"opita-sync-framework/internal/app/previewservice"
	"opita-sync-framework/internal/app/surfaceservice"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/simulation"
	"opita-sync-framework/internal/platform/filesystem"
	"opita-sync-framework/internal/platform/memory"
)

func TestFullSliceCorridor(t *testing.T) {
	resolver, err := filesystem.NewRegistryResolver(filepath.Join("..", "..", "definitions", "capabilities"))
	if err != nil {
		t.Fatalf("resolver bootstrap failed: %v", err)
	}
	artifactStore, err := filesystem.NewArtifactStore(filepath.Join(t.TempDir(), "artifacts"))
	if err != nil {
		t.Fatalf("artifact store bootstrap failed: %v", err)
	}
	contractRepo := memory.NewContractRepository()
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	previewStore := memory.NewPreviewStore()
	intakeStore := memory.NewIntakeStore()
	proposalStore := memory.NewProposalStore()
	recoveryStore := memory.NewRecoveryStore()
	maintenanceStore := memory.NewMaintenanceStore()
	retrievalStore := memory.NewRetrievalStore()
	policyEngine := memory.NewPolicyEngine()

	orchestrator := &foundation.FoundationOrchestrator{
		Compiler:  intent.NewCompiler(contractRepo),
		Policy:    policyEngine,
		Runtime:   runtimeStore,
		Events:    eventLog,
		Registry:  resolver,
		Runs:      runStore,
		Approvals: approvalStore,
	}
	if err := orchestrator.Validate(); err != nil {
		t.Fatalf("invalid orchestrator: %v", err)
	}
	if err := intentservice.Warmup(context.Background(), orchestrator); err != nil {
		t.Fatalf("warmup failed: %v", err)
	}

	intentHandler := intentservice.NewHandler(orchestrator, contractRepo, runtimeStore, eventLog, runStore, resolver, approvalStore)
	previewHandler := previewservice.NewHandler(previewStore, simulation.NewService(policyEngine), eventLog)
	surfaceHandler := surfaceservice.NewHandler(intakeStore, proposalStore, eventLog)
	operatorHandler := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	devHandler := devsurface.NewHandler(runStore, maintenanceStore, eventLog)
	artifactHandler := artifactservice.NewHandler(artifactStore, retrievalStore)

	mux := http.NewServeMux()
	mux.Handle("/", intentHandler.Routes())
	mux.Handle("/v1/intake/", surfaceHandler.Routes())
	mux.Handle("/v1/proposals", surfaceHandler.Routes())
	mux.Handle("/v1/proposals/", surfaceHandler.Routes())
	mux.Handle("/v1/patchsets", surfaceHandler.Routes())
	mux.Handle("/v1/patchsets/", surfaceHandler.Routes())
	mux.Handle("/v1/workspaces/", surfaceHandler.Routes())
	mux.Handle("/v1/previews", previewHandler.Routes())
	mux.Handle("/v1/previews/", previewHandler.Routes())
	mux.Handle("/v1/simulations", previewHandler.Routes())
	mux.Handle("/v1/readable-previews/", previewHandler.Routes())
	mux.Handle("/v1/inspection/", operatorHandler.Routes())
	mux.Handle("/v1/recovery-actions", operatorHandler.Routes())
	mux.Handle("/v1/recovery-actions/", operatorHandler.Routes())
	mux.Handle("/v1/debug/", devHandler.Routes())
	mux.Handle("/v1/maintenance-actions", devHandler.Routes())
	mux.Handle("/v1/maintenance-actions/", devHandler.Routes())
	mux.Handle("/v1/artifacts", artifactHandler.Routes())
	mux.Handle("/v1/artifacts/", artifactHandler.Routes())
	mux.Handle("/v1/retrieval/search", artifactHandler.Routes())

	server := httptest.NewServer(mux)
	defer server.Close()

	turnResp := postJSON(t, server.URL+"/v1/intake/turns", map[string]any{
		"tenant_id":  "tenant-demo",
		"session_id": "session-demo",
		"subject_id": "user-demo",
		"raw_text":   "quiero preparar un cambio gobernado de referencia",
		"trace_id":   "trace-demo-1",
	})
	intakeSession := turnResp["intake_session"].(map[string]any)
	intentCandidate := turnResp["intent_candidate"].(map[string]any)
	if intakeSession["IntakeSessionID"] == "" || intentCandidate["IntentCandidateID"] == "" {
		t.Fatalf("expected intake session and intent candidate ids")
	}

	proposalResp := postJSON(t, server.URL+"/v1/proposals", map[string]any{
		"tenant_id":          "tenant-demo",
		"session_id":         "session-demo",
		"subject_id":         "user-demo",
		"trace_id":           "trace-demo-1",
		"intake_session_id":  getStringField(t, intakeSession, "intake_session_id", "IntakeSessionID"),
		"source_intent_refs": []string{intentCandidate["IntentCandidateID"].(string)},
		"title":              "demo proposal",
		"summary":            "proposal del corredor completo",
		"human_diff_ref":     "human-diff-demo-1",
		"material_diff_ref":  "material-diff-demo-1",
	})
	proposalID := getStringField(t, proposalResp, "proposal_draft_id", "ProposalDraftID")

	patchsetResp := postJSON(t, server.URL+"/v1/patchsets", map[string]any{
		"trace_id":                          "trace-demo-1",
		"proposal_draft_id":                 proposalID,
		"target_artifacts":                  []string{"artifact-demo-1"},
		"material_operations":               []string{"update demo config"},
		"material_diff_hash":                "patchset-hash-demo-1",
		"policy_preview_inputs_ref":         "policy-preview-demo-1",
		"approval_preview_inputs_ref":       "approval-preview-demo-1",
		"classification_preview_inputs_ref": "classification-preview-demo-1",
	})
	patchsetID := getStringField(t, patchsetResp, "patchset_candidate_id", "PatchsetCandidateID")
	workspaceResp := getJSON(t, server.URL+"/v1/workspaces/intake-proposal?intake_session_id="+getStringField(t, intakeSession, "intake_session_id", "IntakeSessionID")+"&intent_candidate_id="+getStringField(t, intentCandidate, "intent_candidate_id", "IntentCandidateID")+"&proposal_draft_id="+proposalID+"&patchset_candidate_id="+patchsetID)
	if len(workspaceResp["next_gates"].([]any)) == 0 {
		t.Fatalf("expected next gates in workspace response")
	}

	previewResp := postJSON(t, server.URL+"/v1/previews", map[string]any{
		"tenant_id":          "tenant-demo",
		"session_id":         "session-demo",
		"subject_id":         "user-demo",
		"trace_id":           "trace-demo-1",
		"proposal_draft_id":  proposalID,
		"contract_id":        "contract-demo-preview",
		"execution_id":       "execution-demo-preview",
		"patchset_ref":       patchsetID,
		"human_diff_ref":     "human-diff-demo-1",
		"material_diff_ref":  "material-diff-demo-1",
		"material_diff_hash": "patchset-hash-demo-1",
		"preview_scope":      "reference-demo",
	})
	previewCandidate := previewResp["preview_candidate"].(map[string]any)
	previewID := getStringField(t, previewCandidate, "preview_candidate_id", "PreviewCandidateID")
	results := previewResp["simulation_results"].([]any)
	if len(results) == 0 {
		t.Fatalf("expected simulation results")
	}
	readablePreview := getJSON(t, server.URL+"/v1/readable-previews/"+previewID)
	if readablePreview["prediction_boundary"] != "simulation_preview_not_kernel_truth" {
		t.Fatalf("expected readable preview boundary, got %#v", readablePreview["prediction_boundary"])
	}

	compileResp := postJSON(t, server.URL+"/v1/intents/compile", map[string]any{
		"request_id":                 "req-demo-1",
		"tenant_id":                  "tenant-demo",
		"workspace_id":               "workspace-demo",
		"user_id":                    "user-demo",
		"session_id":                 "session-demo",
		"trace_id":                   "trace-demo-1",
		"conversation_turn_id":       getStringField(t, turnResp, "conversation_turn_id", "ConversationTurnID"),
		"intake_session_id":          getStringField(t, intakeSession, "intake_session_id", "IntakeSessionID"),
		"intent_candidate_id":        getStringField(t, intentCandidate, "intent_candidate_id", "IntentCandidateID"),
		"proposal_draft_id":          proposalID,
		"patchset_candidate_id":      patchsetID,
		"preview_candidate_id":       previewID,
		"simulation_result_ids":      simulationResultIDs(t, results),
		"objetivo":                   "compilar una intención gobernada de referencia",
		"alcance":                    "reference-demo-scope",
		"tipo_de_resultado_esperado": "execution",
		"autonomia_solicitada":       "assisted",
		"criterios_de_exito":         []string{"evidence trail completo"},
		"restricciones":              []string{"no apply real"},
	})
	executionID := compileResp["execution_id"].(string)
	contractID := compileResp["contract_id"].(string)
	if executionID == "" || contractID == "" {
		t.Fatalf("expected execution and contract ids")
	}

	approvalRaw := compileResp["approval"]
	approvalID := ""
	if approvalRaw != nil {
		approvalID = getStringField(t, approvalRaw.(map[string]any), "approval_request_id", "ApprovalRequestID")
	}
	if approvalID == "" {
		t.Fatalf("expected approval request id for execution flow")
	}

	getJSON(t, server.URL+"/v1/contracts/"+contractID)
	getJSON(t, server.URL+"/v1/executions/"+executionID)
	getJSON(t, server.URL+"/v1/foundation/runs/"+executionID)
	getJSON(t, server.URL+"/v1/approvals/"+approvalID)
	postJSON(t, server.URL+"/v1/approvals/"+approvalID+"/release", map[string]any{
		"decided_by_subject_id": "approver-demo",
		"decision_reason_codes": []string{"approval.release.demo"},
	})
	inspectionResp := getJSON(t, server.URL+"/v1/inspection/executions/"+executionID)
	if len(inspectionResp["conversation_turn_refs"].([]any)) == 0 || len(inspectionResp["proposal_draft_refs"].([]any)) == 0 || len(inspectionResp["preview_candidate_refs"].([]any)) == 0 {
		t.Fatalf("expected correlated refs in inspection view, got %#v", inspectionResp)
	}
	getJSON(t, server.URL+"/v1/debug/semantic?execution_id="+executionID)

	postJSON(t, server.URL+"/v1/maintenance-actions", map[string]any{
		"tenant_id":               "tenant-demo",
		"requested_by_subject_id": "user-demo",
		"action_type":             "request_human_review",
		"target_refs":             []string{executionID},
	})

	bodyB64 := base64.StdEncoding.EncodeToString([]byte("demo reference artifact"))
	postJSON(t, server.URL+"/v1/artifacts", map[string]any{
		"artifact_ref":         "artifact-demo-1",
		"tenant_id":            "tenant-demo",
		"kind":                 "demo",
		"classification_level": "internal",
		"content_type":         "text/plain",
		"body_base64":          bodyB64,
		"title":                "demo artifact",
		"index_text":           "demo reference artifact",
	})
	getJSON(t, server.URL+"/v1/artifacts/artifact-demo-1")
	searchResp := postJSON(t, server.URL+"/v1/retrieval/search", map[string]any{
		"tenant_id": "tenant-demo",
		"text":      "reference",
		"limit":     10,
	})
	if len(searchResp["results"].([]any)) == 0 {
		t.Fatalf("expected retrieval results")
	}

	eventsResp := getJSON(t, server.URL+"/v1/events?execution_id="+executionID)
	if len(eventsResp["records"].([]any)) == 0 {
		t.Fatalf("expected event records for execution")
	}

	_ = previewID
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
	t.Fatalf("missing string field in payload, tried keys: %v, payload: %#v", keys, payload)
	return ""
}

func simulationResultIDs(t *testing.T, results []any) []string {
	t.Helper()
	out := make([]string, 0, len(results))
	for _, item := range results {
		result := item.(map[string]any)
		out = append(out, getStringField(t, result, "simulation_result_id", "SimulationResultID"))
	}
	return out
}

func postJSON(t *testing.T, url string, body any) map[string]any {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var body map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&body)
		t.Fatalf("unexpected status %d for %s: %#v", resp.StatusCode, url, body)
	}
	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return payload
}

func getJSON(t *testing.T, url string) map[string]any {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var body map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&body)
		t.Fatalf("unexpected status %d for %s: %#v", resp.StatusCode, url, body)
	}
	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return payload
}
