package pilotservice

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/events"
)

type EventReader interface {
	Records() []events.Record
}

type Handler struct {
	Events EventReader
}

type timingMetrics struct {
	IntentionToProposalSeconds        float64 `json:"intention_to_proposal_seconds"`
	ProposalToPreviewSeconds          float64 `json:"proposal_to_preview_seconds"`
	PreviewToApprovalSeconds          float64 `json:"preview_to_approval_seconds"`
	ApprovalToExecutionSeconds        float64 `json:"approval_to_execution_seconds"`
	IncidentToRecoveryDecisionSeconds float64 `json:"incident_to_recovery_decision_seconds"`
}

type governanceMetrics struct {
	GovernanceBlocks      int `json:"governance_blocks"`
	ApprovalsRequired     int `json:"approvals_required"`
	SuccessfulReleases    int `json:"successful_releases"`
	FingerprintMismatches int `json:"fingerprint_mismatches"`
}

type operabilityMetrics struct {
	CasesWithFullEvidenceTrail int `json:"cases_with_full_evidence_trail"`
	EndToEndReconstructable    int `json:"end_to_end_reconstructable"`
	RecoveryBlocked            int `json:"recovery_blocked"`
	RecoveryExecuted           int `json:"recovery_executed"`
	ExecutionsObserved         int `json:"executions_observed"`
}

type pilotScorecard struct {
	TenantID          string             `json:"tenant_id,omitempty"`
	ScenarioID        string             `json:"scenario_id,omitempty"`
	GeneratedAt       time.Time          `json:"generated_at"`
	Timing            timingMetrics      `json:"timing"`
	Governance        governanceMetrics  `json:"governance"`
	Operability       operabilityMetrics `json:"operability"`
	EventCount        int                `json:"event_count"`
	TrackedExecutions []string           `json:"tracked_executions"`
	Boundary          string             `json:"boundary"`
}

type incidentCandidate struct {
	TenantID          string   `json:"tenant_id,omitempty"`
	ScenarioID        string   `json:"scenario_id,omitempty"`
	ExecutionID       string   `json:"execution_id,omitempty"`
	TraceID           string   `json:"trace_id,omitempty"`
	Severity          string   `json:"severity"`
	Category          string   `json:"category"`
	Summary           string   `json:"summary"`
	ReasonCodes       []string `json:"reason_codes,omitempty"`
	EventType         string   `json:"event_type"`
	ApprovalRequestID string   `json:"approval_request_id,omitempty"`
	RecoveryActionID  string   `json:"recovery_action_id,omitempty"`
}

func NewHandler(events EventReader) *Handler {
	return &Handler{Events: events}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/pilot/scorecard", h.handleScorecard)
	mux.HandleFunc("GET /v1/pilot/scorecard/scenarios", h.handleScenarioScorecards)
	mux.HandleFunc("GET /v1/pilot/incident-candidates", h.handleIncidentCandidates)
	return mux
}

func (h *Handler) handleScorecard(w http.ResponseWriter, r *http.Request) {
	if h.Events == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "pilot.events_not_ready"})
		return
	}
	tenantID := strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	records := h.Events.Records()
	if tenantID != "" {
		filtered := make([]events.Record, 0, len(records))
		for _, record := range records {
			if record.TenantID == tenantID {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}
	scorecard := buildScorecard(records, tenantID)
	writeJSON(w, http.StatusOK, scorecard)
}

func (h *Handler) handleScenarioScorecards(w http.ResponseWriter, r *http.Request) {
	if h.Events == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "pilot.events_not_ready"})
		return
	}
	tenantID := strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	records := h.Events.Records()
	if tenantID != "" {
		filtered := make([]events.Record, 0, len(records))
		for _, record := range records {
			if record.TenantID == tenantID {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}
	byTrace := map[string][]events.Record{}
	for _, record := range records {
		traceID := strings.TrimSpace(record.TraceID)
		if traceID == "" {
			continue
		}
		byTrace[traceID] = append(byTrace[traceID], record)
	}
	out := make([]pilotScorecard, 0, len(byTrace))
	for traceID, grouped := range byTrace {
		scorecard := buildScorecard(grouped, tenantID)
		scorecard.ScenarioID = traceID
		out = append(out, scorecard)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ScenarioID < out[j].ScenarioID })
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant_id": tenantID,
		"scenarios": out,
		"boundary":  "scenario_scorecards_derived_from_trace_id",
	})
}

func (h *Handler) handleIncidentCandidates(w http.ResponseWriter, r *http.Request) {
	if h.Events == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "pilot.events_not_ready"})
		return
	}
	tenantID := strings.TrimSpace(r.URL.Query().Get("tenant_id"))
	scenarioID := strings.TrimSpace(r.URL.Query().Get("scenario_id"))
	records := h.Events.Records()
	filtered := make([]events.Record, 0, len(records))
	for _, record := range records {
		if tenantID != "" && record.TenantID != tenantID {
			continue
		}
		if scenarioID != "" && strings.TrimSpace(record.TraceID) != scenarioID {
			continue
		}
		filtered = append(filtered, record)
	}
	candidates := buildIncidentCandidates(filtered)
	writeJSON(w, http.StatusOK, map[string]any{
		"tenant_id":   tenantID,
		"scenario_id": scenarioID,
		"candidates":  candidates,
		"boundary":    "incident_candidates_derived_from_canonical_event_log",
	})
}

func buildScorecard(records []events.Record, tenantID string) pilotScorecard {
	byExecution := map[string][]events.Record{}
	tracked := map[string]struct{}{}
	for _, record := range records {
		if record.ExecutionID != "" {
			byExecution[record.ExecutionID] = append(byExecution[record.ExecutionID], record)
			tracked[record.ExecutionID] = struct{}{}
		}
	}
	scorecard := pilotScorecard{
		TenantID:          tenantID,
		GeneratedAt:       time.Now().UTC(),
		EventCount:        len(records),
		TrackedExecutions: mapKeys(tracked),
		Boundary:          "pilot_metrics_derived_from_canonical_event_log",
	}
	sort.Strings(scorecard.TrackedExecutions)
	timing := timingAccumulator{}
	for executionID, executionRecords := range byExecution {
		_ = executionID
		timing.consume(executionRecords)
		if hasFullEvidenceTrail(executionRecords) {
			scorecard.Operability.CasesWithFullEvidenceTrail++
			scorecard.Operability.EndToEndReconstructable++
		}
		scorecard.Operability.ExecutionsObserved++
		for _, record := range executionRecords {
			switch record.EventType {
			case "approval.awaiting":
				scorecard.Governance.ApprovalsRequired++
			case "approval.released":
				scorecard.Governance.SuccessfulReleases++
			case "approval.fingerprint_mismatch":
				scorecard.Governance.FingerprintMismatches++
			case "recovery.execution_blocked":
				scorecard.Operability.RecoveryBlocked++
			case "execution.released", "execution.unknown_outcome", "compensation.requested":
				if record.RecoveryActionID != "" {
					scorecard.Operability.RecoveryExecuted++
				}
			}
			if record.EventType == "approval.awaiting" || (record.EventType == "execution.created" && payloadState(record) == "blocked") {
				scorecard.Governance.GovernanceBlocks++
			}
		}
	}
	scorecard.Timing = timing.average()
	return scorecard
}

type timingAccumulator struct {
	intentionToProposal        []float64
	proposalToPreview          []float64
	previewToApproval          []float64
	approvalToExecution        []float64
	incidentToRecoveryDecision []float64
}

func (a *timingAccumulator) consume(records []events.Record) {
	sort.Slice(records, func(i, j int) bool { return records[i].OccurredAt.Before(records[j].OccurredAt) })
	find := func(eventType string) *events.Record {
		for i := range records {
			if records[i].EventType == eventType {
				return &records[i]
			}
		}
		return nil
	}
	if start, end := find("intake.turn_recorded"), find("proposal.created"); start != nil && end != nil {
		a.intentionToProposal = append(a.intentionToProposal, end.OccurredAt.Sub(start.OccurredAt).Seconds())
	}
	if start, end := find("proposal.created"), find("preview.created"); start != nil && end != nil {
		a.proposalToPreview = append(a.proposalToPreview, end.OccurredAt.Sub(start.OccurredAt).Seconds())
	}
	if start, end := find("preview.created"), find("approval.awaiting"); start != nil && end != nil {
		a.previewToApproval = append(a.previewToApproval, end.OccurredAt.Sub(start.OccurredAt).Seconds())
	}
	if start, end := find("approval.awaiting"), find("approval.released"); start != nil && end != nil {
		a.approvalToExecution = append(a.approvalToExecution, end.OccurredAt.Sub(start.OccurredAt).Seconds())
	}
	if start, end := find("recovery.candidate_created"), firstOf(records, "execution.released", "execution.unknown_outcome", "compensation.requested"); start != nil && end != nil && end.RecoveryActionID != "" {
		a.incidentToRecoveryDecision = append(a.incidentToRecoveryDecision, end.OccurredAt.Sub(start.OccurredAt).Seconds())
	}
}

func (a timingAccumulator) average() timingMetrics {
	return timingMetrics{
		IntentionToProposalSeconds:        avg(a.intentionToProposal),
		ProposalToPreviewSeconds:          avg(a.proposalToPreview),
		PreviewToApprovalSeconds:          avg(a.previewToApproval),
		ApprovalToExecutionSeconds:        avg(a.approvalToExecution),
		IncidentToRecoveryDecisionSeconds: avg(a.incidentToRecoveryDecision),
	}
}

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func firstOf(records []events.Record, eventTypes ...string) *events.Record {
	set := map[string]struct{}{}
	for _, eventType := range eventTypes {
		set[eventType] = struct{}{}
	}
	for i := range records {
		if _, ok := set[records[i].EventType]; ok {
			return &records[i]
		}
	}
	return nil
}

func hasFullEvidenceTrail(records []events.Record) bool {
	required := map[string]bool{
		"contract.compilation_completed": false,
		"policy.decision_recorded":       false,
		"execution.created":              false,
	}
	for _, record := range records {
		if _, ok := required[record.EventType]; ok {
			required[record.EventType] = true
		}
	}
	for _, present := range required {
		if !present {
			return false
		}
	}
	return true
}

func payloadState(record events.Record) string {
	if state, ok := record.Payload["runtime_state"].(string); ok {
		return state
	}
	return ""
}

func mapKeys(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	return out
}

func buildIncidentCandidates(records []events.Record) []incidentCandidate {
	out := make([]incidentCandidate, 0)
	for _, record := range records {
		candidate, ok := buildIncidentCandidate(record)
		if ok {
			out = append(out, candidate)
		}
	}
	return out
}

func buildIncidentCandidate(record events.Record) (incidentCandidate, bool) {
	candidate := incidentCandidate{
		TenantID:          record.TenantID,
		ScenarioID:        record.TraceID,
		ExecutionID:       record.ExecutionID,
		TraceID:           record.TraceID,
		EventType:         record.EventType,
		ApprovalRequestID: record.ApprovalRequestID,
		RecoveryActionID:  record.RecoveryActionID,
	}
	switch record.EventType {
	case "approval.fingerprint_mismatch":
		candidate.Severity = "high"
		candidate.Category = "fingerprint_mismatch"
		candidate.Summary = "approved fingerprint no longer matches the current execution contract"
		candidate.ReasonCodes = []string{"approval.fingerprint_mismatch"}
		return candidate, true
	case "recovery.execution_blocked":
		candidate.Severity = "high"
		candidate.Category = "recovery_blocked"
		candidate.Summary = "recovery candidate became blocked by runtime or approval constraints"
		candidate.ReasonCodes = payloadStringSlice(record.Payload, "reason_codes")
		return candidate, true
	case "preview.simulation_recorded":
		status, _ := record.Payload["status"].(string)
		if status == "preview_warning" || status == "preview_blocked" {
			candidate.Severity = map[string]string{"preview_warning": "medium", "preview_blocked": "high"}[status]
			candidate.Category = "preview_signal"
			candidate.Summary = "preview simulation emitted a warning or blocked signal"
			candidate.ReasonCodes = []string{status}
			return candidate, true
		}
	}
	return incidentCandidate{}, false
}

func payloadStringSlice(payload map[string]any, key string) []string {
	values, ok := payload[key]
	if !ok {
		return nil
	}
	switch typed := values.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, value := range typed {
			if s, ok := value.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
