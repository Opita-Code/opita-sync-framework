package cerbos_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/platform/cerbos"
)

func TestCerbosClientEvaluateAllowMapsToRequireApprovalWhenNeeded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/check/resources" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"execute": map[string]any{"effect": "EFFECT_ALLOW"},
			},
		})
	}))
	defer srv.Close()

	client := cerbos.NewClient(srv.URL)
	decision, err := client.Evaluate(policy.Input{
		TenantID:              "tenant-1",
		ContractID:            "contract-1",
		ExecutionID:           "exec-1",
		ResourceKind:          "execution_record",
		Action:                "execute",
		ClassificationLevel:   "internal",
		ApprovalModeEffective: "pre_execution",
		RiskLevel:             "medium",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if decision.Decision != policy.DecisionRequireApproval {
		t.Fatalf("expected require_approval, got %s", decision.Decision)
	}
}

func TestCerbosClientEvaluateDeny(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"execute": map[string]any{"effect": "EFFECT_DENY"},
			},
		})
	}))
	defer srv.Close()

	client := cerbos.NewClient(srv.URL)
	decision, err := client.Evaluate(policy.Input{
		TenantID:            "tenant-1",
		ContractID:          "contract-1",
		ExecutionID:         "exec-1",
		ResourceKind:        "execution_record",
		Action:              "execute",
		ClassificationLevel: "internal",
		RiskLevel:           "low",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if decision.Decision != policy.DecisionDenyBlock {
		t.Fatalf("expected deny_block, got %s", decision.Decision)
	}
}
