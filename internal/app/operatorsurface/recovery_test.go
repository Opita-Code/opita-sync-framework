package operatorsurface_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opita-sync-framework/internal/app/operatorsurface"
	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/inspection"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/registry"
	"opita-sync-framework/internal/engine/runtime"
	"opita-sync-framework/internal/platform/memory"
)

func TestCreateRecoveryCandidateForApprovalResume(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()

	exec := runtime.ExecutionRecord{ExecutionID: "exec-2", TenantID: "tenant-1", ContractID: "contract-2", ContractFingerprint: "fp-2", TraceID: "trace-2", State: runtime.ExecutionStateAwaitingApproval, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = runtimeStore.CreateExecution(exec)
	approval := approvals.Request{ApprovalRequestID: "approval-2", ExecutionID: "exec-2", ContractID: "contract-2", TenantID: "tenant-1", TraceID: "trace-2", State: approvals.StateAwaitingApproval, Mode: "pre_execution", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = approvalStore.Create(approval)
	_ = runStore.Save(foundation.FoundationRunResult{
		Contract:       intent.CompiledContract{ContractID: "contract-2", Fingerprint: "fp-2", TenantID: "tenant-1"},
		Execution:      exec,
		PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-2", Decision: policy.DecisionRequireApproval},
		Resolution:     registry.ResolutionResult{CapabilityManifestRef: "manifest://capability.execution.default", BindingID: "binding-2", ProviderRef: "provider://y"},
		Approval:       &approval,
	})

	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	body, _ := json.Marshal(map[string]any{
		"execution_id":            "exec-2",
		"requested_action":        string(inspection.RecoveryResumeAfterApproval),
		"requested_by_subject_id": "operator-1",
		"approval_request_id":     "approval-2",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestExecuteRecoveryCandidateBlocksWhenNotReady(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-missing", TenantID: "tenant-1", ContractID: "contract-blocked", ContractFingerprint: "fp-blocked", TraceID: "trace-blocked", State: runtime.ExecutionStateBlocked, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	candidate := inspection.RecoveryActionCandidate{
		RecoveryActionCandidateID: "recovery-1",
		ExecutionID:               "exec-missing",
		RequestedAction:           inspection.RecoveryResumeAfterApproval,
		RequestedBySubjectID:      "operator-1",
		CurrentRuntimeState:       string(runtime.ExecutionStateBlocked),
		ReadyForExecution:         false,
		State:                     inspection.RecoveryCandidatePending,
		CreatedAt:                 time.Now().UTC(),
		UpdatedAt:                 time.Now().UTC(),
	}
	_ = recoveryStore.Create(candidate)
	req := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions/recovery-1/execute", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateRecoveryCandidateRequiresRequestedBySubjectID(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-3", TenantID: "tenant-1", ContractID: "contract-3", ContractFingerprint: "fp-3", TraceID: "trace-3", State: runtime.ExecutionStateUnknownOutcome, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	body, _ := json.Marshal(map[string]any{
		"execution_id":     "exec-3",
		"requested_action": string(inspection.RecoveryAcknowledgeUnknown),
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUnsupportedRecoveryCandidateStartsBlocked(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-4", TenantID: "tenant-1", ContractID: "contract-4", ContractFingerprint: "fp-4", TraceID: "trace-4", State: runtime.ExecutionStateBlocked, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	body, _ := json.Marshal(map[string]any{
		"execution_id":            "exec-4",
		"requested_action":        string(inspection.RecoveryRetryTechnicalStep),
		"requested_by_subject_id": "operator-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var candidate inspection.RecoveryActionCandidate
	if err := json.Unmarshal(w.Body.Bytes(), &candidate); err != nil {
		t.Fatalf("unmarshal candidate: %v", err)
	}
	if candidate.State != inspection.RecoveryCandidateBlocked || candidate.ReadyForExecution {
		t.Fatalf("expected blocked unsupported candidate, got %+v", candidate)
	}
}

func TestExecuteAcknowledgeUnknownPreservesUnknownOutcome(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-5", TenantID: "tenant-1", ContractID: "contract-5", ContractFingerprint: "fp-5", TraceID: "trace-5", State: runtime.ExecutionStateUnknownOutcome, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	body, _ := json.Marshal(map[string]any{
		"execution_id":            "exec-5",
		"requested_action":        string(inspection.RecoveryAcknowledgeUnknown),
		"requested_by_subject_id": "operator-1",
	})
	createReq := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	h.Routes().ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createW.Code)
	}
	var candidate inspection.RecoveryActionCandidate
	if err := json.Unmarshal(createW.Body.Bytes(), &candidate); err != nil {
		t.Fatalf("unmarshal candidate: %v", err)
	}
	execReq := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions/"+candidate.RecoveryActionCandidateID+"/execute", nil)
	execW := httptest.NewRecorder()
	h.Routes().ServeHTTP(execW, execReq)
	if execW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", execW.Code, execW.Body.String())
	}
	exec, found, err := runtimeStore.GetExecution("exec-5")
	if err != nil || !found {
		t.Fatalf("expected execution, found=%v err=%v", found, err)
	}
	if exec.State != runtime.ExecutionStateUnknownOutcome {
		t.Fatalf("expected unknown_outcome to be preserved, got %s", exec.State)
	}
}

func TestExecuteManualCompensationFromFailedTransitionsToCompensationPending(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	_ = runtimeStore.CreateExecution(runtime.ExecutionRecord{ExecutionID: "exec-6", TenantID: "tenant-1", ContractID: "contract-6", ContractFingerprint: "fp-6", TraceID: "trace-6", State: runtime.ExecutionStateFailed, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	body, _ := json.Marshal(map[string]any{
		"execution_id":            "exec-6",
		"requested_action":        string(inspection.RecoveryRequestManualComp),
		"requested_by_subject_id": "operator-1",
	})
	createReq := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	h.Routes().ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createW.Code)
	}
	var candidate inspection.RecoveryActionCandidate
	if err := json.Unmarshal(createW.Body.Bytes(), &candidate); err != nil {
		t.Fatalf("unmarshal candidate: %v", err)
	}
	execReq := httptest.NewRequest(http.MethodPost, "/v1/recovery-actions/"+candidate.RecoveryActionCandidateID+"/execute", nil)
	execW := httptest.NewRecorder()
	h.Routes().ServeHTTP(execW, execReq)
	if execW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", execW.Code, execW.Body.String())
	}
	exec, found, err := runtimeStore.GetExecution("exec-6")
	if err != nil || !found {
		t.Fatalf("expected execution, found=%v err=%v", found, err)
	}
	if exec.State != runtime.ExecutionStateCompensationPending {
		t.Fatalf("expected compensation_pending, got %s", exec.State)
	}
}

func TestOperatorWorkspaceShowsBlockedStateClearly(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()
	exec := runtime.ExecutionRecord{ExecutionID: "exec-7", TenantID: "tenant-1", ContractID: "contract-7", ContractFingerprint: "fp-7", TraceID: "trace-7", State: runtime.ExecutionStateBlocked, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = runtimeStore.CreateExecution(exec)
	_ = runStore.Save(foundation.FoundationRunResult{Contract: intent.CompiledContract{ContractID: "contract-7", Fingerprint: "fp-7", TenantID: "tenant-1"}, Execution: exec, PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-7", Decision: policy.DecisionDenyBlock}, Resolution: registry.ResolutionResult{CapabilityManifestRef: "manifest://capability.execution.default", BindingID: "binding-7", ProviderRef: "provider://z"}})
	_ = eventLog.Append(events.Record{EventID: "event-7", EventType: "policy.decision_recorded", ExecutionID: "exec-7", TenantID: "tenant-1", TraceID: "trace-7", OccurredAt: time.Now().UTC()})
	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	req := httptest.NewRequest(http.MethodGet, "/v1/operator/executions/exec-7/workspace", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	lifecycle := resp["lifecycle"].(map[string]any)
	if lifecycle["lifecycle_label"] != "blocked" {
		t.Fatalf("expected blocked lifecycle, got %#v", lifecycle["lifecycle_label"])
	}
}
