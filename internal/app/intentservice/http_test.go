package intentservice_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"opita-sync-framework/internal/app/intentservice"
	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/runtime"
	"opita-sync-framework/internal/platform/filesystem"
	"opita-sync-framework/internal/platform/memory"
)

func newIntentHandler(t *testing.T) *intentservice.Handler {
	t.Helper()
	contractRepo := memory.NewContractRepository()
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	resolver, err := filesystem.NewRegistryResolver(filepath.Join("..", "..", "..", "definitions", "capabilities"))
	if err != nil {
		t.Fatalf("resolver bootstrap failed: %v", err)
	}
	orchestrator := &foundation.FoundationOrchestrator{
		Compiler:  intent.NewCompiler(contractRepo),
		Policy:    memory.NewPolicyEngine(),
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
	return intentservice.NewHandler(orchestrator, contractRepo, runtimeStore, eventLog, runStore, resolver, approvalStore)
}

func TestCompileIntentReturnsCreated(t *testing.T) {
	h := newIntentHandler(t)
	body, _ := json.Marshal(map[string]any{
		"request_id":                 "req-1",
		"tenant_id":                  "tenant-1",
		"workspace_id":               "workspace-1",
		"user_id":                    "user-1",
		"session_id":                 "session-1",
		"objetivo":                   "compilar una intención gobernada",
		"alcance":                    "slice-test",
		"tipo_de_resultado_esperado": "plan",
		"autonomia_solicitada":       "assisted",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/intents/compile", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetContractRequiresID(t *testing.T) {
	h := newIntentHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/contracts/", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReleaseApprovalRequiresDecider(t *testing.T) {
	approvalStore := memory.NewApprovalStore()
	runtimeStore := memory.NewRuntimeService()
	request := approvals.Request{
		ApprovalRequestID:         "approval-1",
		ExecutionID:               "exec-1",
		ContractID:                "contract-1",
		TenantID:                  "tenant-1",
		TraceID:                   "trace-1",
		State:                     approvals.StateAwaitingApproval,
		Mode:                      "pre_execution",
		SourceContractFingerprint: "fp-1",
		CreatedAt:                 time.Now().UTC(),
		UpdatedAt:                 time.Now().UTC(),
	}
	_ = approvalStore.Create(request)
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", ContractFingerprint: "fp-1", TraceID: "trace-1", State: runtime.ExecutionStateAwaitingApproval, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := &intentservice.Handler{Approvals: approvalStore, Runtime: runtimeStore, Events: memory.NewEventLog()}
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals/approval-1/release", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReleaseApprovalRejectsFingerprintMismatch(t *testing.T) {
	approvalStore := memory.NewApprovalStore()
	runtimeStore := memory.NewRuntimeService()
	request := approvals.Request{
		ApprovalRequestID:         "approval-1",
		ExecutionID:               "exec-1",
		ContractID:                "contract-1",
		TenantID:                  "tenant-1",
		TraceID:                   "trace-1",
		State:                     approvals.StateAwaitingApproval,
		Mode:                      "pre_execution",
		SourceContractFingerprint: "fp-approved",
		CreatedAt:                 time.Now().UTC(),
		UpdatedAt:                 time.Now().UTC(),
	}
	_ = approvalStore.Create(request)
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", ContractFingerprint: "fp-current", TraceID: "trace-1", State: runtime.ExecutionStateAwaitingApproval, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := &intentservice.Handler{Approvals: approvalStore, Runtime: runtimeStore, Events: memory.NewEventLog()}
	body, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals/approval-1/release", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestRejectApprovalReturnsDecision(t *testing.T) {
	approvalStore := memory.NewApprovalStore()
	runtimeStore := memory.NewRuntimeService()
	request := approvals.Request{
		ApprovalRequestID:         "approval-1",
		ExecutionID:               "exec-1",
		ContractID:                "contract-1",
		TenantID:                  "tenant-1",
		TraceID:                   "trace-1",
		State:                     approvals.StateAwaitingApproval,
		Mode:                      "pre_execution",
		SourceContractFingerprint: "fp-1",
		CreatedAt:                 time.Now().UTC(),
		UpdatedAt:                 time.Now().UTC(),
	}
	_ = approvalStore.Create(request)
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", ContractFingerprint: "fp-1", TraceID: "trace-1", State: runtime.ExecutionStateAwaitingApproval, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := &intentservice.Handler{Approvals: approvalStore, Runtime: runtimeStore, Events: memory.NewEventLog()}
	body, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1", "decision_reason_codes": []string{"approval.reject.manual"}})
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals/approval-1/reject", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
