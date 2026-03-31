package pilotservice_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opita-sync-framework/internal/app/pilotservice"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/platform/memory"
)

func TestPilotScorecardAggregatesMetricsByTenant(t *testing.T) {
	eventLog := memory.NewEventLog()
	base := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)
	records := []events.Record{
		{EventID: "1", EventType: "intake.turn_recorded", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base},
		{EventID: "2", EventType: "proposal.created", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(10 * time.Second)},
		{EventID: "3", EventType: "preview.created", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(20 * time.Second)},
		{EventID: "4", EventType: "contract.compilation_completed", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(25 * time.Second)},
		{EventID: "5", EventType: "policy.decision_recorded", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(30 * time.Second)},
		{EventID: "6", EventType: "execution.created", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(35 * time.Second), Payload: map[string]any{"runtime_state": "awaiting_approval"}},
		{EventID: "7", EventType: "approval.awaiting", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(40 * time.Second)},
		{EventID: "8", EventType: "approval.released", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(50 * time.Second)},
		{EventID: "9", EventType: "recovery.candidate_created", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", RecoveryActionID: "rec-1", OccurredAt: base.Add(60 * time.Second)},
		{EventID: "10", EventType: "execution.released", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", RecoveryActionID: "rec-1", OccurredAt: base.Add(70 * time.Second)},
		{EventID: "11", EventType: "approval.fingerprint_mismatch", TenantID: "tenant-alpha-ops", ExecutionID: "exec-1", OccurredAt: base.Add(80 * time.Second)},
		{EventID: "12", EventType: "execution.created", TenantID: "tenant-beta-governance", ExecutionID: "exec-2", OccurredAt: base.Add(90 * time.Second), Payload: map[string]any{"runtime_state": "blocked"}},
	}
	for _, record := range records {
		if err := eventLog.Append(record); err != nil {
			t.Fatalf("append event: %v", err)
		}
	}
	h := pilotservice.NewHandler(eventLog)
	req := httptest.NewRequest(http.MethodGet, "/v1/pilot/scorecard?tenant_id=tenant-alpha-ops", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	governance := resp["governance"].(map[string]any)
	if governance["approvals_required"].(float64) != 1 || governance["successful_releases"].(float64) != 1 || governance["fingerprint_mismatches"].(float64) != 1 {
		t.Fatalf("unexpected governance metrics: %#v", governance)
	}
	operability := resp["operability"].(map[string]any)
	if operability["recovery_executed"].(float64) != 1 {
		t.Fatalf("expected recovery_executed=1, got %#v", operability)
	}
	if operability["cases_with_full_evidence_trail"].(float64) != 1 {
		t.Fatalf("expected full evidence trail case, got %#v", operability)
	}
}

func TestPilotScenarioScorecardsGroupByTraceID(t *testing.T) {
	eventLog := memory.NewEventLog()
	base := time.Date(2026, 3, 30, 11, 0, 0, 0, time.UTC)
	records := []events.Record{
		{EventID: "1", EventType: "intake.turn_recorded", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base},
		{EventID: "2", EventType: "proposal.created", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base.Add(10 * time.Second)},
		{EventID: "3", EventType: "contract.compilation_completed", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base.Add(20 * time.Second)},
		{EventID: "4", EventType: "policy.decision_recorded", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base.Add(30 * time.Second)},
		{EventID: "5", EventType: "execution.created", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base.Add(40 * time.Second), Payload: map[string]any{"runtime_state": "execution_released"}},
		{EventID: "6", EventType: "intake.turn_recorded", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-02", ExecutionID: "exec-2", OccurredAt: base.Add(50 * time.Second)},
		{EventID: "7", EventType: "proposal.created", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-02", ExecutionID: "exec-2", OccurredAt: base.Add(60 * time.Second)},
		{EventID: "8", EventType: "contract.compilation_completed", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-02", ExecutionID: "exec-2", OccurredAt: base.Add(70 * time.Second)},
		{EventID: "9", EventType: "policy.decision_recorded", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-02", ExecutionID: "exec-2", OccurredAt: base.Add(80 * time.Second)},
		{EventID: "10", EventType: "execution.created", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-02", ExecutionID: "exec-2", OccurredAt: base.Add(90 * time.Second), Payload: map[string]any{"runtime_state": "blocked"}},
	}
	for _, record := range records {
		if err := eventLog.Append(record); err != nil {
			t.Fatalf("append event: %v", err)
		}
	}
	h := pilotservice.NewHandler(eventLog)
	req := httptest.NewRequest(http.MethodGet, "/v1/pilot/scorecard/scenarios?tenant_id=tenant-alpha-ops", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	scenarios := resp["scenarios"].([]any)
	if len(scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(scenarios))
	}
	first := scenarios[0].(map[string]any)
	if first["scenario_id"] != "trace-alpha-01" {
		t.Fatalf("expected first scenario trace-alpha-01, got %#v", first["scenario_id"])
	}
}

func TestPilotIncidentCandidatesDeriveFromCanonicalEvents(t *testing.T) {
	eventLog := memory.NewEventLog()
	base := time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC)
	records := []events.Record{
		{EventID: "1", EventType: "preview.simulation_recorded", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", OccurredAt: base, Payload: map[string]any{"status": "preview_warning"}},
		{EventID: "2", EventType: "recovery.execution_blocked", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", RecoveryActionID: "recovery-1", OccurredAt: base.Add(10 * time.Second), Payload: map[string]any{"reason_codes": []string{"recovery.invalid_runtime_state.resume_after_approval"}}},
		{EventID: "3", EventType: "approval.fingerprint_mismatch", TenantID: "tenant-alpha-ops", TraceID: "trace-alpha-01", ExecutionID: "exec-1", ApprovalRequestID: "approval-1", OccurredAt: base.Add(20 * time.Second)},
	}
	for _, record := range records {
		if err := eventLog.Append(record); err != nil {
			t.Fatalf("append event: %v", err)
		}
	}
	h := pilotservice.NewHandler(eventLog)
	req := httptest.NewRequest(http.MethodGet, "/v1/pilot/incident-candidates?tenant_id=tenant-alpha-ops&scenario_id=trace-alpha-01", nil)
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	candidates := resp["candidates"].([]any)
	if len(candidates) != 3 {
		t.Fatalf("expected 3 incident candidates, got %d", len(candidates))
	}
}
