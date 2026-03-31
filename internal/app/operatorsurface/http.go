package operatorsurface

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/inspection"
	"opita-sync-framework/internal/engine/runtime"
)

type RuntimeReader interface {
	GetExecution(executionID string) (runtime.ExecutionRecord, bool, error)
	UpdateExecutionState(executionID string, state runtime.ExecutionState) (runtime.ExecutionRecord, error)
}

type EventReader interface {
	RecordsByExecution(executionID string) []events.Record
	Append(record events.Record) error
}

type RunReader interface {
	GetByExecutionID(executionID string) (foundation.FoundationRunResult, bool, error)
}

type ApprovalReader interface {
	GetByID(approvalRequestID string) (approvals.Request, bool, error)
	Decide(approvalRequestID string, decision approvals.Decision) (approvals.Request, error)
}

type Handler struct {
	Runtime   RuntimeReader
	Events    EventReader
	Runs      RunReader
	Approvals ApprovalReader
	Recovery  inspection.RecoveryStore
}

type operatorExecutionWorkspace struct {
	ExecutionID   string                    `json:"execution_id"`
	TenantID      string                    `json:"tenant_id"`
	TraceID       string                    `json:"trace_id"`
	Lifecycle     operatorLifecycleCard     `json:"lifecycle"`
	Outcome       operatorOutcomeCard       `json:"outcome"`
	EvidenceTrail operatorEvidenceTrailCard `json:"evidence_trail"`
	Recovery      operatorRecoveryCard      `json:"recovery"`
	Audit         operatorAuditCard         `json:"audit"`
	Boundary      string                    `json:"boundary"`
}

type operatorLifecycleCard struct {
	CurrentRuntimeState string   `json:"current_runtime_state"`
	CurrentOutcomeState string   `json:"current_outcome_state"`
	LifecycleLabel      string   `json:"lifecycle_label"`
	VisibleStates       []string `json:"visible_states"`
	Summary             string   `json:"summary"`
}

type operatorOutcomeCard struct {
	OutcomeLabel     string   `json:"outcome_label"`
	OperatorFindings []string `json:"operator_findings"`
	Summary          string   `json:"summary"`
}

type operatorEvidenceTrailCard struct {
	EventCount           int      `json:"event_count"`
	PolicyDecisionRefs   []string `json:"policy_decision_refs"`
	ApprovalRequestRefs  []string `json:"approval_request_refs"`
	ConversationTurnRefs []string `json:"conversation_turn_refs"`
	ProposalDraftRefs    []string `json:"proposal_draft_refs"`
	PreviewCandidateRefs []string `json:"preview_candidate_refs"`
	SimulationResultRefs []string `json:"simulation_result_refs"`
	Summary              string   `json:"summary"`
}

type operatorRecoveryCard struct {
	CanTriggerRecovery bool                                 `json:"can_trigger_recovery"`
	SuggestedActions   []string                             `json:"suggested_actions"`
	Candidates         []inspection.RecoveryActionCandidate `json:"candidates"`
	Summary            string                               `json:"summary"`
}

type operatorAuditCard struct {
	AllActionsAudited bool   `json:"all_actions_audited"`
	MutationPolicy    string `json:"mutation_policy"`
	Summary           string `json:"summary"`
}

type createRecoveryRequest struct {
	ExecutionID          string `json:"execution_id"`
	RequestedAction      string `json:"requested_action"`
	RequestedBySubjectID string `json:"requested_by_subject_id"`
	ApprovalRequestID    string `json:"approval_request_id,omitempty"`
}

func NewHandler(runtimeReader RuntimeReader, eventReader EventReader, runReader RunReader, approvalReader ApprovalReader, recoveryStore inspection.RecoveryStore) *Handler {
	return &Handler{Runtime: runtimeReader, Events: eventReader, Runs: runReader, Approvals: approvalReader, Recovery: recoveryStore}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/inspection/executions/", h.handleInspectionView)
	mux.HandleFunc("GET /v1/operator/executions/", h.handleOperatorExecutionWorkspace)
	mux.HandleFunc("POST /v1/recovery-actions", h.handleCreateRecoveryCandidate)
	mux.HandleFunc("GET /v1/recovery-actions/", h.handleGetRecoveryCandidate)
	mux.HandleFunc("POST /v1/recovery-actions/", h.handleExecuteRecoveryCandidate)
	return mux
}

func (h *Handler) handleInspectionView(w http.ResponseWriter, r *http.Request) {
	executionID := strings.TrimPrefix(r.URL.Path, "/v1/inspection/executions/")
	if executionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "inspection.missing_execution_id"})
		return
	}
	run, found, err := h.Runs.GetByExecutionID(executionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "inspection.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "inspection.not_found"})
		return
	}
	records := h.Events.RecordsByExecution(executionID)
	eventRefs := make([]string, 0, len(records))
	conversationTurnRefs := map[string]struct{}{}
	intakeSessionRefs := map[string]struct{}{}
	intentCandidateRefs := map[string]struct{}{}
	proposalDraftRefs := map[string]struct{}{}
	patchsetCandidateRefs := map[string]struct{}{}
	previewCandidateRefs := map[string]struct{}{}
	simulationResultRefs := map[string]struct{}{}
	for _, record := range records {
		eventRefs = append(eventRefs, record.EventID)
		if record.ConversationTurnID != "" {
			conversationTurnRefs[record.ConversationTurnID] = struct{}{}
		}
		if record.IntakeSessionID != "" {
			intakeSessionRefs[record.IntakeSessionID] = struct{}{}
		}
		if record.IntentCandidateID != "" {
			intentCandidateRefs[record.IntentCandidateID] = struct{}{}
		}
		if record.ProposalDraftID != "" {
			proposalDraftRefs[record.ProposalDraftID] = struct{}{}
		}
		if record.PatchsetCandidateID != "" {
			patchsetCandidateRefs[record.PatchsetCandidateID] = struct{}{}
		}
		if record.PreviewCandidateID != "" {
			previewCandidateRefs[record.PreviewCandidateID] = struct{}{}
		}
		if record.SimulationResultID != "" {
			simulationResultRefs[record.SimulationResultID] = struct{}{}
		}
		switch ids := record.Payload["simulation_result_ids"].(type) {
		case []any:
			for _, id := range ids {
				if s, ok := id.(string); ok && s != "" {
					simulationResultRefs[s] = struct{}{}
				}
			}
		case []string:
			for _, id := range ids {
				if id != "" {
					simulationResultRefs[id] = struct{}{}
				}
			}
		}
	}
	approvalRefs := []string{}
	if run.Approval != nil {
		approvalRefs = append(approvalRefs, run.Approval.ApprovalRequestID)
	}
	view := h.buildInspectionView(run, records)
	writeJSON(w, http.StatusOK, view)
}

func (h *Handler) handleOperatorExecutionWorkspace(w http.ResponseWriter, r *http.Request) {
	executionID := strings.TrimPrefix(r.URL.Path, "/v1/operator/executions/")
	executionID = strings.TrimSuffix(executionID, "/workspace")
	executionID = strings.TrimSuffix(executionID, "/")
	if executionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "operator_workspace.missing_execution_id"})
		return
	}
	run, found, err := h.Runs.GetByExecutionID(executionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "operator_workspace.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "operator_workspace.not_found"})
		return
	}
	records := h.Events.RecordsByExecution(executionID)
	inspectionView := h.buildInspectionView(run, records)
	candidates, err := h.Recovery.ListByExecution(executionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "operator_workspace.recovery_lookup_failed", "message": err.Error()})
		return
	}
	workspace := operatorExecutionWorkspace{
		ExecutionID: run.Execution.ExecutionID,
		TenantID:    run.Execution.TenantID,
		TraceID:     run.Execution.TraceID,
		Lifecycle: operatorLifecycleCard{
			CurrentRuntimeState: inspectionView.CurrentRuntimeState,
			CurrentOutcomeState: inspectionView.CurrentOutcomeState,
			LifecycleLabel:      deriveLifecycleLabel(run.Execution.State),
			VisibleStates:       []string{"blocked", "failed", "compensation_pending", "unknown_outcome", "execution_released"},
			Summary:             fmt.Sprintf("execution %s lifecycle is %s", run.Execution.ExecutionID, deriveLifecycleLabel(run.Execution.State)),
		},
		Outcome: operatorOutcomeCard{
			OutcomeLabel:     inspectionView.CurrentOutcomeState,
			OperatorFindings: inspectionView.OperatorFindings,
			Summary:          inspectionView.OperatorSummary,
		},
		EvidenceTrail: operatorEvidenceTrailCard{
			EventCount:           len(inspectionView.EventRefs),
			PolicyDecisionRefs:   inspectionView.PolicyDecisionRefs,
			ApprovalRequestRefs:  inspectionView.ApprovalRequestRefs,
			ConversationTurnRefs: inspectionView.ConversationTurnRefs,
			ProposalDraftRefs:    inspectionView.ProposalDraftRefs,
			PreviewCandidateRefs: inspectionView.PreviewCandidateRefs,
			SimulationResultRefs: inspectionView.SimulationResultRefs,
			Summary:              fmt.Sprintf("evidence trail exposes %d canonical events", len(inspectionView.EventRefs)),
		},
		Recovery: operatorRecoveryCard{
			CanTriggerRecovery: canTriggerRecovery(run.Execution.State),
			SuggestedActions:   suggestedRecoveryActions(run.Execution.State),
			Candidates:         candidates,
			Summary:            fmt.Sprintf("%d recovery candidates linked to this execution", len(candidates)),
		},
		Audit: operatorAuditCard{
			AllActionsAudited: true,
			MutationPolicy:    "no_direct_canonical_mutation_outside_recovery_and_kernel_routes",
			Summary:           "all operator actions remain auditable and runtime mutations stay inside governed routes",
		},
		Boundary: "operator_surface_reads_and_requests_kernel_executes",
	}
	writeJSON(w, http.StatusOK, workspace)
}

func (h *Handler) buildInspectionView(run foundation.FoundationRunResult, records []events.Record) inspection.ExecutionInspectionView {
	eventRefs := make([]string, 0, len(records))
	conversationTurnRefs := map[string]struct{}{}
	intakeSessionRefs := map[string]struct{}{}
	intentCandidateRefs := map[string]struct{}{}
	proposalDraftRefs := map[string]struct{}{}
	patchsetCandidateRefs := map[string]struct{}{}
	previewCandidateRefs := map[string]struct{}{}
	simulationResultRefs := map[string]struct{}{}
	for _, record := range records {
		eventRefs = append(eventRefs, record.EventID)
		if record.ConversationTurnID != "" {
			conversationTurnRefs[record.ConversationTurnID] = struct{}{}
		}
		if record.IntakeSessionID != "" {
			intakeSessionRefs[record.IntakeSessionID] = struct{}{}
		}
		if record.IntentCandidateID != "" {
			intentCandidateRefs[record.IntentCandidateID] = struct{}{}
		}
		if record.ProposalDraftID != "" {
			proposalDraftRefs[record.ProposalDraftID] = struct{}{}
		}
		if record.PatchsetCandidateID != "" {
			patchsetCandidateRefs[record.PatchsetCandidateID] = struct{}{}
		}
		if record.PreviewCandidateID != "" {
			previewCandidateRefs[record.PreviewCandidateID] = struct{}{}
		}
		if record.SimulationResultID != "" {
			simulationResultRefs[record.SimulationResultID] = struct{}{}
		}
		switch ids := record.Payload["simulation_result_ids"].(type) {
		case []any:
			for _, id := range ids {
				if s, ok := id.(string); ok && s != "" {
					simulationResultRefs[s] = struct{}{}
				}
			}
		case []string:
			for _, id := range ids {
				if id != "" {
					simulationResultRefs[id] = struct{}{}
				}
			}
		}
	}
	approvalRefs := []string{}
	if run.Approval != nil {
		approvalRefs = append(approvalRefs, run.Approval.ApprovalRequestID)
	}
	return inspection.ExecutionInspectionView{
		InspectionViewID:      fmt.Sprintf("inspection-%d", time.Now().UTC().UnixNano()),
		ExecutionID:           run.Execution.ExecutionID,
		TenantID:              run.Execution.TenantID,
		TraceID:               run.Execution.TraceID,
		ContractID:            run.Execution.ContractID,
		CurrentRuntimeState:   string(run.Execution.State),
		CurrentOutcomeState:   string(run.Execution.State),
		PolicyDecisionRefs:    []string{run.PolicyDecision.PolicyDecisionID},
		ApprovalRequestRefs:   approvalRefs,
		EventRefs:             eventRefs,
		ConversationTurnRefs:  mapKeys(conversationTurnRefs),
		IntakeSessionRefs:     mapKeys(intakeSessionRefs),
		IntentCandidateRefs:   mapKeys(intentCandidateRefs),
		ProposalDraftRefs:     mapKeys(proposalDraftRefs),
		PatchsetCandidateRefs: mapKeys(patchsetCandidateRefs),
		PreviewCandidateRefs:  mapKeys(previewCandidateRefs),
		SimulationResultRefs:  mapKeys(simulationResultRefs),
		ResolvedCapability:    run.Resolution.CapabilityManifestRef,
		ResolvedBinding:       run.Resolution.BindingID,
		ResolvedProvider:      run.Resolution.ProviderRef,
		OperatorSummary:       fmt.Sprintf("execution %s is currently %s", run.Execution.ExecutionID, run.Execution.State),
		OperatorFindings:      []string{"foundation run reconstructed from canonical records"},
		GeneratedAt:           time.Now().UTC(),
	}
}

func (h *Handler) handleCreateRecoveryCandidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var req createRecoveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.invalid_json", "message": err.Error()})
		return
	}
	runtimeRecord, found, err := h.Runtime.GetExecution(req.ExecutionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.execution_lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "recovery.execution_not_found"})
		return
	}
	if strings.TrimSpace(req.RequestedBySubjectID) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.missing_requested_by_subject_id"})
		return
	}
	action := inspection.RecoveryAction(req.RequestedAction)
	ready := false
	state := inspection.RecoveryCandidateBlocked
	reasonCodes := []string{}
	blockingConstraints := []string{}
	preconditionsRefs := []string{req.ExecutionID}
	switch action {
	case inspection.RecoveryResumeAfterApproval:
		if req.ApprovalRequestID == "" {
			reasonCodes = []string{"recovery.missing_approval_request_id"}
			blockingConstraints = []string{"approval_request_id_required"}
			break
		}
		preconditionsRefs = append(preconditionsRefs, req.ApprovalRequestID)
		approval, found, err := h.Approvals.GetByID(req.ApprovalRequestID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.approval_lookup_failed", "message": err.Error()})
			return
		}
		if !found {
			reasonCodes = []string{"recovery.approval_not_found"}
			blockingConstraints = []string{"approval_request_missing"}
			break
		}
		if approval.ExecutionID != req.ExecutionID {
			reasonCodes = []string{"recovery.approval_execution_mismatch"}
			blockingConstraints = []string{"approval_must_match_execution"}
			break
		}
		if approval.State != approvals.StateAwaitingApproval {
			reasonCodes = []string{"recovery.approval_not_awaiting"}
			blockingConstraints = []string{"approval_must_be_awaiting"}
			break
		}
		if runtimeRecord.State != runtime.ExecutionStateAwaitingApproval {
			reasonCodes = []string{"recovery.invalid_runtime_state.resume_after_approval"}
			blockingConstraints = []string{"execution_must_be_awaiting_approval"}
			break
		}
		ready = true
		state = inspection.RecoveryCandidatePending
		reasonCodes = []string{"recovery.resume_after_approval"}
	case inspection.RecoveryAcknowledgeUnknown:
		if runtimeRecord.State != runtime.ExecutionStateUnknownOutcome {
			reasonCodes = []string{"recovery.invalid_runtime_state.acknowledge_unknown"}
			blockingConstraints = []string{"execution_must_be_unknown_outcome"}
			break
		}
		ready = true
		state = inspection.RecoveryCandidatePending
		reasonCodes = []string{"recovery.acknowledge_unknown_outcome"}
	case inspection.RecoveryRequestManualComp:
		if runtimeRecord.State != runtime.ExecutionStateFailed && runtimeRecord.State != runtime.ExecutionStateUnknownOutcome {
			reasonCodes = []string{"recovery.invalid_runtime_state.manual_compensation"}
			blockingConstraints = []string{"execution_must_be_failed_or_unknown_outcome"}
			break
		}
		ready = true
		state = inspection.RecoveryCandidatePending
		reasonCodes = []string{"recovery.request_manual_compensation"}
	default:
		reasonCodes = []string{"recovery.action_not_supported_in_v1"}
		blockingConstraints = []string{"action_out_of_scope_for_v1"}
	}
	candidate := inspection.RecoveryActionCandidate{
		RecoveryActionCandidateID:  fmt.Sprintf("recovery-%d", time.Now().UTC().UnixNano()),
		ExecutionID:                req.ExecutionID,
		RequestedAction:            action,
		RequestedBySubjectID:       req.RequestedBySubjectID,
		CurrentRuntimeState:        string(runtimeRecord.State),
		ApprovalRequestID:          req.ApprovalRequestID,
		PreconditionsRefs:          preconditionsRefs,
		BlockingConstraints:        blockingConstraints,
		ReasonCodes:                reasonCodes,
		ReadyForExecution:          ready,
		RequiresAdditionalApproval: false,
		State:                      state,
		CreatedAt:                  time.Now().UTC(),
		UpdatedAt:                  time.Now().UTC(),
	}
	if err := h.Recovery.Create(candidate); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.create_failed", "message": err.Error()})
		return
	}
	_ = h.Events.Append(events.Record{
		EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
		EventType:           "recovery.candidate_created",
		TenantID:            runtimeRecord.TenantID,
		TraceID:             runtimeRecord.TraceID,
		ContractID:          runtimeRecord.ContractID,
		ContractFingerprint: runtimeRecord.ContractFingerprint,
		ExecutionID:         runtimeRecord.ExecutionID,
		ApprovalRequestID:   req.ApprovalRequestID,
		RecoveryActionID:    candidate.RecoveryActionCandidateID,
		OccurredAt:          time.Now().UTC(),
		Payload: map[string]any{
			"recovery_action_candidate_id": candidate.RecoveryActionCandidateID,
			"requested_action":             candidate.RequestedAction,
			"requested_by_subject_id":      candidate.RequestedBySubjectID,
			"current_runtime_state":        candidate.CurrentRuntimeState,
			"ready_for_execution":          candidate.ReadyForExecution,
			"reason_codes":                 candidate.ReasonCodes,
		},
	})
	writeJSON(w, http.StatusCreated, candidate)
}

func (h *Handler) handleGetRecoveryCandidate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/recovery-actions/")
	if id == "" || strings.HasSuffix(id, "/execute") {
		return
	}
	candidate, found, err := h.Recovery.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "recovery.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, candidate)
}

func (h *Handler) handleExecuteRecoveryCandidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/recovery-actions/")
	if !strings.HasSuffix(id, "/execute") {
		return
	}
	id = strings.TrimSuffix(id, "/execute")
	id = strings.TrimSuffix(id, "/")
	candidate, found, err := h.Recovery.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "recovery.not_found"})
		return
	}
	execution, executionFound, err := h.Runtime.GetExecution(candidate.ExecutionID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.execution_lookup_failed", "message": err.Error()})
		return
	}
	if !executionFound {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "recovery.execution_not_found"})
		return
	}
	if candidate.State != inspection.RecoveryCandidatePending || !candidate.ReadyForExecution {
		candidate = h.blockRecoveryCandidate(candidate, execution, []string{"recovery.not_ready_or_not_pending"}, []string{"candidate_must_be_pending_and_ready"})
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.not_ready", "candidate": candidate})
		return
	}
	if candidate.RequestedAction == inspection.RecoveryResumeAfterApproval {
		if execution.State != runtime.ExecutionStateAwaitingApproval {
			candidate = h.blockRecoveryCandidate(candidate, execution, []string{"recovery.invalid_runtime_state.resume_after_approval"}, []string{"execution_must_be_awaiting_approval"})
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.invalid_runtime_state", "candidate": candidate})
			return
		}
		if _, err := h.Approvals.Decide(candidate.ApprovalRequestID, approvals.Decision{
			State:               approvals.StateReleased,
			DecidedBySubjectID:  candidate.RequestedBySubjectID,
			DecisionComment:     "release via recovery resume_after_approval",
			DecisionReasonCodes: candidate.ReasonCodes,
			DecidedAt:           time.Now().UTC(),
		}); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.approval_release_failed", "message": err.Error()})
			return
		}
		execution, err := h.Runtime.UpdateExecutionState(candidate.ExecutionID, runtime.ExecutionStateExecutionReleased)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.execution_release_failed", "message": err.Error()})
			return
		}
		candidate.State = inspection.RecoveryCandidateExecuted
		candidate.UpdatedAt = time.Now().UTC()
		_ = h.Recovery.Update(candidate)
		_ = h.Events.Append(events.Record{
			EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
			EventType:           "execution.released",
			TenantID:            execution.TenantID,
			TraceID:             execution.TraceID,
			ContractID:          execution.ContractID,
			ContractFingerprint: execution.ContractFingerprint,
			ExecutionID:         execution.ExecutionID,
			ApprovalRequestID:   candidate.ApprovalRequestID,
			RecoveryActionID:    candidate.RecoveryActionCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"recovery_action_candidate_id": candidate.RecoveryActionCandidateID,
				"requested_action":             candidate.RequestedAction,
				"requested_by_subject_id":      candidate.RequestedBySubjectID,
				"current_runtime_state_before": string(runtime.ExecutionStateAwaitingApproval),
				"resulting_runtime_state":      execution.State,
				"reason_codes":                 candidate.ReasonCodes,
			},
		})
		writeJSON(w, http.StatusOK, map[string]any{"candidate": candidate, "execution": execution})
		return
	}
	if candidate.RequestedAction == inspection.RecoveryAcknowledgeUnknown {
		if execution.State != runtime.ExecutionStateUnknownOutcome {
			candidate = h.blockRecoveryCandidate(candidate, execution, []string{"recovery.invalid_runtime_state.acknowledge_unknown"}, []string{"execution_must_be_unknown_outcome"})
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.invalid_runtime_state", "candidate": candidate})
			return
		}
		execution, err := h.Runtime.UpdateExecutionState(candidate.ExecutionID, runtime.ExecutionStateUnknownOutcome)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.unknown_outcome_failed", "message": err.Error()})
			return
		}
		candidate.State = inspection.RecoveryCandidateExecuted
		candidate.UpdatedAt = time.Now().UTC()
		_ = h.Recovery.Update(candidate)
		_ = h.Events.Append(events.Record{
			EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
			EventType:           "execution.unknown_outcome",
			TenantID:            execution.TenantID,
			TraceID:             execution.TraceID,
			ContractID:          execution.ContractID,
			ContractFingerprint: execution.ContractFingerprint,
			ExecutionID:         execution.ExecutionID,
			RecoveryActionID:    candidate.RecoveryActionCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"recovery_action_candidate_id": candidate.RecoveryActionCandidateID,
				"requested_action":             candidate.RequestedAction,
				"requested_by_subject_id":      candidate.RequestedBySubjectID,
				"current_runtime_state_before": string(runtime.ExecutionStateUnknownOutcome),
				"resulting_runtime_state":      execution.State,
				"reason_codes":                 candidate.ReasonCodes,
			},
		})
		writeJSON(w, http.StatusOK, map[string]any{"candidate": candidate, "execution": execution})
		return
	}
	if candidate.RequestedAction == inspection.RecoveryRequestManualComp {
		if execution.State != runtime.ExecutionStateFailed && execution.State != runtime.ExecutionStateUnknownOutcome {
			candidate = h.blockRecoveryCandidate(candidate, execution, []string{"recovery.invalid_runtime_state.manual_compensation"}, []string{"execution_must_be_failed_or_unknown_outcome"})
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.invalid_runtime_state", "candidate": candidate})
			return
		}
		execution, err := h.Runtime.UpdateExecutionState(candidate.ExecutionID, runtime.ExecutionStateCompensationPending)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "recovery.compensation_pending_failed", "message": err.Error()})
			return
		}
		candidate.State = inspection.RecoveryCandidateExecuted
		candidate.UpdatedAt = time.Now().UTC()
		_ = h.Recovery.Update(candidate)
		_ = h.Events.Append(events.Record{
			EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
			EventType:           "compensation.requested",
			TenantID:            execution.TenantID,
			TraceID:             execution.TraceID,
			ContractID:          execution.ContractID,
			ContractFingerprint: execution.ContractFingerprint,
			ExecutionID:         execution.ExecutionID,
			RecoveryActionID:    candidate.RecoveryActionCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"recovery_action_candidate_id": candidate.RecoveryActionCandidateID,
				"requested_action":             candidate.RequestedAction,
				"requested_by_subject_id":      candidate.RequestedBySubjectID,
				"current_runtime_state_before": candidate.CurrentRuntimeState,
				"resulting_runtime_state":      execution.State,
				"reason_codes":                 candidate.ReasonCodes,
			},
		})
		writeJSON(w, http.StatusOK, map[string]any{"candidate": candidate, "execution": execution})
		return
	}
	candidate = h.blockRecoveryCandidate(candidate, execution, []string{"recovery.action_not_supported_in_v1"}, []string{"action_out_of_scope_for_v1"})
	writeJSON(w, http.StatusBadRequest, map[string]any{"error": "recovery.action_not_supported_in_slice", "candidate": candidate})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func mapKeys(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	return out
}

func deriveLifecycleLabel(state runtime.ExecutionState) string {
	switch state {
	case runtime.ExecutionStateBlocked:
		return "blocked"
	case runtime.ExecutionStateFailed:
		return "failed"
	case runtime.ExecutionStateCompensationPending:
		return "compensating"
	case runtime.ExecutionStateUnknownOutcome:
		return "unknown_outcome"
	default:
		return string(state)
	}
}

func canTriggerRecovery(state runtime.ExecutionState) bool {
	switch state {
	case runtime.ExecutionStateAwaitingApproval, runtime.ExecutionStateFailed, runtime.ExecutionStateUnknownOutcome:
		return true
	default:
		return false
	}
}

func suggestedRecoveryActions(state runtime.ExecutionState) []string {
	switch state {
	case runtime.ExecutionStateAwaitingApproval:
		return []string{string(inspection.RecoveryResumeAfterApproval)}
	case runtime.ExecutionStateUnknownOutcome:
		return []string{string(inspection.RecoveryAcknowledgeUnknown), string(inspection.RecoveryRequestManualComp)}
	case runtime.ExecutionStateFailed:
		return []string{string(inspection.RecoveryRequestManualComp)}
	default:
		return []string{}
	}
}

func (h *Handler) blockRecoveryCandidate(candidate inspection.RecoveryActionCandidate, execution runtime.ExecutionRecord, reasonCodes []string, blockingConstraints []string) inspection.RecoveryActionCandidate {
	candidate.State = inspection.RecoveryCandidateBlocked
	candidate.ReadyForExecution = false
	candidate.ReasonCodes = reasonCodes
	candidate.BlockingConstraints = blockingConstraints
	candidate.UpdatedAt = time.Now().UTC()
	_ = h.Recovery.Update(candidate)
	if h.Events == nil {
		return candidate
	}
	_ = h.Events.Append(events.Record{
		EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
		EventType:           "recovery.execution_blocked",
		TenantID:            execution.TenantID,
		TraceID:             execution.TraceID,
		ContractID:          execution.ContractID,
		ContractFingerprint: execution.ContractFingerprint,
		ExecutionID:         execution.ExecutionID,
		ApprovalRequestID:   candidate.ApprovalRequestID,
		RecoveryActionID:    candidate.RecoveryActionCandidateID,
		OccurredAt:          time.Now().UTC(),
		Payload: map[string]any{
			"recovery_action_candidate_id": candidate.RecoveryActionCandidateID,
			"requested_action":             candidate.RequestedAction,
			"requested_by_subject_id":      candidate.RequestedBySubjectID,
			"current_runtime_state":        execution.State,
			"reason_codes":                 reasonCodes,
			"blocking_constraints":         blockingConstraints,
		},
	})
	return candidate
}
