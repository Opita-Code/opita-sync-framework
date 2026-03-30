package devsurface_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opita-sync-framework/internal/app/devsurface"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/maintenance"
	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/registry"
	"opita-sync-framework/internal/engine/runtime"
	"opita-sync-framework/internal/platform/memory"
)

func seededRunStore(t *testing.T) *memory.FoundationRunStore {
	t.Helper()
	runStore := memory.NewFoundationRunStore()
	_ = runStore.Save(foundation.FoundationRunResult{
		Contract:       intent.CompiledContract{ContractID: "contract-1", Fingerprint: "fp-1", TenantID: "tenant-1"},
		Execution:      runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", ContractFingerprint: "fp-1", TraceID: "trace-1", State: runtime.ExecutionStateExecutionReleased, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-1", Decision: policy.DecisionAllow},
		Resolution:     registry.ResolutionResult{CapabilityManifestRef: "manifest://capability.plan.default", BundleDigest: "sha256:aaa", BindingID: "binding-1", ProviderRef: "provider://demo", Resolved: true},
	})
	return runStore
}

func TestSemanticDebugViewReturnsOK(t *testing.T) {
	runStore := seededRunStore(t)
	h := devsurface.NewHandler(runStore, memory.NewMaintenanceStore(), memory.NewEventLog())
	req := httptest.NewRequest(http.MethodGet, "/v1/debug/semantic?execution_id=exec-1", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateMaintenanceCandidateReturnsCreated(t *testing.T) {
	runStore := seededRunStore(t)
	h := devsurface.NewHandler(runStore, memory.NewMaintenanceStore(), memory.NewEventLog())
	body, _ := json.Marshal(map[string]any{
		"tenant_id":               "tenant-1",
		"requested_by_subject_id": "user-1",
		"action_type":             string(maintenance.ActionRequestHumanReview),
		"target_refs":             []string{"exec-1"},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/maintenance-actions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
