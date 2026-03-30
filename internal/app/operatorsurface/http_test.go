package operatorsurface_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opita-sync-framework/internal/app/operatorsurface"
	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/foundation"
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
