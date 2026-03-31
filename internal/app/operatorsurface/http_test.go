package operatorsurface_test

import (
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

func TestInspectionViewReturnsCorrelatedView(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()

	exec := runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", ContractFingerprint: "fp-1", TraceID: "trace-1", State: runtime.ExecutionStateExecutionReleased, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = runtimeStore.CreateExecution(exec)
	approval := approvals.Request{ApprovalRequestID: "approval-1", ExecutionID: "exec-1", ContractID: "contract-1", TenantID: "tenant-1", TraceID: "trace-1", State: approvals.StateReleased, Mode: "pre_execution", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = approvalStore.Create(approval)
	_ = runStore.Save(foundation.FoundationRunResult{
		Contract:       intent.CompiledContract{ContractID: "contract-1", Fingerprint: "fp-1", TenantID: "tenant-1"},
		Execution:      exec,
		PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-1", Decision: policy.DecisionRequireApproval},
		Resolution:     registry.ResolutionResult{CapabilityManifestRef: "manifest://capability.plan.default", BindingID: "binding-1", ProviderRef: "provider://x"},
		Approval:       &approval,
	})

	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	req := httptest.NewRequest(http.MethodGet, "/v1/inspection/executions/exec-1", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["execution_id"] != "exec-1" {
		t.Fatalf("expected execution_id exec-1")
	}
}

func TestOperatorWorkspaceReturnsUsableLifecycleAndRecovery(t *testing.T) {
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	recoveryStore := memory.NewRecoveryStore()

	exec := runtime.ExecutionRecord{ExecutionID: "exec-2", TenantID: "tenant-1", ContractID: "contract-2", ContractFingerprint: "fp-2", TraceID: "trace-2", State: runtime.ExecutionStateUnknownOutcome, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = runtimeStore.CreateExecution(exec)
	approval := approvals.Request{ApprovalRequestID: "approval-2", ExecutionID: "exec-2", ContractID: "contract-2", TenantID: "tenant-1", TraceID: "trace-2", State: approvals.StateAwaitingApproval, Mode: "pre_execution", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	_ = approvalStore.Create(approval)
	_ = recoveryStore.Create(inspection.RecoveryActionCandidate{RecoveryActionCandidateID: "recovery-2", ExecutionID: "exec-2", RequestedAction: inspection.RecoveryAcknowledgeUnknown, RequestedBySubjectID: "operator-1", CurrentRuntimeState: string(runtime.ExecutionStateUnknownOutcome), PreconditionsRefs: []string{"exec-2"}, ReasonCodes: []string{"recovery.acknowledge_unknown_outcome"}, ReadyForExecution: true, State: inspection.RecoveryCandidatePending, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	_ = runStore.Save(foundation.FoundationRunResult{
		Contract:       intent.CompiledContract{ContractID: "contract-2", Fingerprint: "fp-2", TenantID: "tenant-1", ProposalDraftID: "proposal-2", PreviewCandidateID: "preview-2"},
		Execution:      exec,
		PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-2", Decision: policy.DecisionRequireApproval},
		Resolution:     registry.ResolutionResult{CapabilityManifestRef: "manifest://capability.execution.default", BindingID: "binding-2", ProviderRef: "provider://y"},
		Approval:       &approval,
	})
	_ = eventLog.Append(events.Record{EventID: "event-1", EventType: "execution.unknown_outcome", ExecutionID: "exec-2", TenantID: "tenant-1", TraceID: "trace-2", ProposalDraftID: "proposal-2", PreviewCandidateID: "preview-2", SimulationResultID: "sim-2", OccurredAt: time.Now().UTC()})

	h := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	req := httptest.NewRequest(http.MethodGet, "/v1/operator/executions/exec-2/workspace", nil)
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
	if lifecycle["lifecycle_label"] != "unknown_outcome" {
		t.Fatalf("expected unknown_outcome lifecycle, got %#v", lifecycle["lifecycle_label"])
	}
	recovery := resp["recovery"].(map[string]any)
	if recovery["can_trigger_recovery"] != true {
		t.Fatalf("expected can_trigger_recovery true, got %#v", recovery["can_trigger_recovery"])
	}
	evidence := resp["evidence_trail"].(map[string]any)
	if evidence["event_count"].(float64) < 1 {
		t.Fatalf("expected evidence events, got %#v", evidence)
	}
}
