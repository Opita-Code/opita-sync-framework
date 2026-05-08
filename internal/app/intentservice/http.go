package intentservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/registry"
	"opita-sync-framework/internal/engine/runtime"
	"opita-sync-framework/internal/httputil"
)

type RuntimeReader interface {
	GetExecution(ctx context.Context, executionID string) (runtime.ExecutionRecord, bool, error)
	UpdateExecutionState(ctx context.Context, executionID string, state runtime.ExecutionState) (runtime.ExecutionRecord, error)
}

type EventReader interface {
	RecordsByExecution(ctx context.Context, executionID string) []events.Record
	Append(ctx context.Context, record events.Record) error
}

type ContractReader interface {
	GetByID(ctx context.Context, contractID string) (intent.CompiledContract, bool, error)
}

type RegistryReader interface {
	Resolve(ctx context.Context, req registry.ResolutionRequest) (registry.ResolutionResult, error)
}

type ApprovalReader interface {
	GetByID(ctx context.Context, approvalRequestID string) (approvals.Request, bool, error)
	Decide(ctx context.Context, approvalRequestID string, decision approvals.Decision) (approvals.Request, error)
}

type Handler struct {
	Orchestrator *foundation.FoundationOrchestrator
	Contracts    ContractReader
	Runtime      RuntimeReader
	Events       EventReader
	Runs         RunReader
	Registry     RegistryReader
	Approvals    ApprovalReader
}

type RunReader interface {
	GetByExecutionID(ctx context.Context, executionID string) (foundation.FoundationRunResult, bool, error)
}

type compileRequest struct {
	RequestID             string   `json:"request_id"`
	RequestedCapabilityID string   `json:"requested_capability_id,omitempty"`
	TenantID              string   `json:"tenant_id"`
	WorkspaceID           string   `json:"workspace_id"`
	UserID                string   `json:"user_id"`
	SessionID             string   `json:"session_id"`
	TraceID               string   `json:"trace_id,omitempty"`
	ConversationTurnID    string   `json:"conversation_turn_id,omitempty"`
	IntakeSessionID       string   `json:"intake_session_id,omitempty"`
	IntentCandidateID     string   `json:"intent_candidate_id,omitempty"`
	ProposalDraftID       string   `json:"proposal_draft_id,omitempty"`
	PatchsetCandidateID   string   `json:"patchset_candidate_id,omitempty"`
	PreviewCandidateID    string   `json:"preview_candidate_id,omitempty"`
	SimulationResultIDs   []string `json:"simulation_result_ids,omitempty"`
	Objetivo              string   `json:"objetivo"`
	Alcance               string   `json:"alcance"`
	TipoResultadoEsperado string   `json:"tipo_de_resultado_esperado"`
	Restricciones         []string `json:"restricciones,omitempty"`
	SistemasPosibles      []string `json:"sistemas_posibles,omitempty"`
	DatosPermitidos       []string `json:"datos_permitidos,omitempty"`
	AutonomiaSolicitada   string   `json:"autonomia_solicitada,omitempty"`
	AprobacionRequerida   string   `json:"aprobacion_requerida,omitempty"`
	CriteriosDeExito      []string `json:"criterios_de_exito,omitempty"`
}

type compileResponse struct {
	ContractID     string                   `json:"contract_id"`
	Fingerprint    string                   `json:"fingerprint"`
	Compilation    intent.CompilationStatus `json:"compilation_status"`
	ExecutionID    string                   `json:"execution_id"`
	ExecutionState runtime.ExecutionState   `json:"execution_state"`
	PolicyDecision string                   `json:"policy_decision"`
	ProviderRef    string                   `json:"provider_ref"`
	BundleDigest   string                   `json:"bundle_digest"`
	BindingID      string                   `json:"binding_id"`
	Approval       *approvals.Request       `json:"approval,omitempty"`
	Diagnostics    []intent.Diagnostic      `json:"diagnostics,omitempty"`
	Correlations   map[string]string        `json:"correlations"`
}

func NewHandler(orchestrator *foundation.FoundationOrchestrator, contractReader ContractReader, runtimeReader RuntimeReader, eventReader EventReader, runReader RunReader, registryReader RegistryReader, approvalReader ApprovalReader) *Handler {
	return &Handler{Orchestrator: orchestrator, Contracts: contractReader, Runtime: runtimeReader, Events: eventReader, Runs: runReader, Registry: registryReader, Approvals: approvalReader}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("POST /v1/intents/compile", h.handleCompile)
	mux.HandleFunc("GET /v1/contracts/", h.handleGetContract)
	mux.HandleFunc("GET /v1/executions/", h.handleGetExecution)
	mux.HandleFunc("GET /v1/events", h.handleGetEvents)
	mux.HandleFunc("GET /v1/foundation/runs/", h.handleGetRun)
	mux.HandleFunc("GET /v1/registry/resolve", h.handleResolveCapability)
	mux.HandleFunc("GET /v1/approvals/", h.handleGetApproval)
	mux.HandleFunc("POST /v1/approvals/", h.handleApprovalDecision)
	return mux
}

type approvalDecisionRequest struct {
	DecidedBySubjectID  string   `json:"decided_by_subject_id"`
	DecisionComment     string   `json:"decision_comment,omitempty"`
	DecisionReasonCodes []string `json:"decision_reason_codes,omitempty"`
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "intent-service"})
}

func (h *Handler) handleCompile(w http.ResponseWriter, r *http.Request) {
	if h.Orchestrator == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "orchestrator not configured")
		return
	}
	var req compileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "request.invalid_json", err.Error())
		return
	}

	result, err := h.Orchestrator.Run(r.Context(), intent.IntentInput{
		RequestID:             req.RequestID,
		RequestedCapabilityID: req.RequestedCapabilityID,
		TenantID:              req.TenantID,
		WorkspaceID:           req.WorkspaceID,
		UserID:                req.UserID,
		SessionID:             req.SessionID,
		TraceID:               req.TraceID,
		ConversationTurnID:    req.ConversationTurnID,
		IntakeSessionID:       req.IntakeSessionID,
		IntentCandidateID:     req.IntentCandidateID,
		ProposalDraftID:       req.ProposalDraftID,
		PatchsetCandidateID:   req.PatchsetCandidateID,
		PreviewCandidateID:    req.PreviewCandidateID,
		SimulationResultIDs:   req.SimulationResultIDs,
		Objetivo:              req.Objetivo,
		Alcance:               req.Alcance,
		TipoResultadoEsperado: intent.ResultType(req.TipoResultadoEsperado),
		Restricciones:         req.Restricciones,
		SistemasPosibles:      req.SistemasPosibles,
		DatosPermitidos:       req.DatosPermitidos,
		AutonomiaSolicitada:   intent.Autonomy(req.AutonomiaSolicitada),
		AprobacionRequerida:   req.AprobacionRequerida,
		CriteriosDeExito:      req.CriteriosDeExito,
		CreatedAt:             time.Now().UTC(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		code := "compiler.failed"
		if errors.Is(err, intent.ErrInvalidIntentInput) {
			status = http.StatusBadRequest
			code = "compiler.invalid_input"
		}
		httputil.WriteError(w, status, code, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, compileResponse{
		ContractID:     result.Contract.ContractID,
		Fingerprint:    result.Contract.Fingerprint,
		Compilation:    result.Report.Status,
		ExecutionID:    result.Execution.ExecutionID,
		ExecutionState: result.Execution.State,
		PolicyDecision: string(result.PolicyDecision.Decision),
		ProviderRef:    result.Resolution.ProviderRef,
		BundleDigest:   result.Resolution.BundleDigest,
		BindingID:      result.Resolution.BindingID,
		Approval:       result.Approval,
		Diagnostics:    result.Report.Diagnostics,
		Correlations: map[string]string{
			"requested_capability_id": result.Contract.RequestedCapabilityID,
			"tenant_id":               result.Contract.TenantID,
			"contract_id":             result.Contract.ContractID,
			"contract_fingerprint":    result.Contract.Fingerprint,
			"execution_id":            result.Execution.ExecutionID,
			"trace_id":                result.Execution.TraceID,
			"policy_decision_id":      result.PolicyDecision.PolicyDecisionID,
			"conversation_turn_id":    result.Contract.ConversationTurnID,
			"intake_session_id":       result.Contract.IntakeSessionID,
			"intent_candidate_id":     result.Contract.IntentCandidateID,
			"proposal_draft_id":       result.Contract.ProposalDraftID,
			"patchset_candidate_id":   result.Contract.PatchsetCandidateID,
			"preview_candidate_id":    result.Contract.PreviewCandidateID,
		},
	})
}

func (h *Handler) handleGetExecution(w http.ResponseWriter, r *http.Request) {
	if h.Runtime == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "runtime reader not configured")
		return
	}
	executionID := strings.TrimPrefix(r.URL.Path, "/v1/executions/")
	if executionID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "execution.missing_id", "execution id is required")
		return
	}
	record, found, err := h.Runtime.GetExecution(r.Context(), executionID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "execution.lookup_failed", err.Error())
		return
	}
	if !found {
		httputil.WriteError(w, http.StatusNotFound, "execution.not_found", "execution not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, record)
}

func (h *Handler) handleGetContract(w http.ResponseWriter, r *http.Request) {
	if h.Contracts == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "contract reader not configured")
		return
	}
	contractID := strings.TrimPrefix(r.URL.Path, "/v1/contracts/")
	if contractID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "contract.missing_id", "contract id is required")
		return
	}
	contract, found, err := h.Contracts.GetByID(r.Context(), contractID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "contract.lookup_failed", err.Error())
		return
	}
	if !found {
		httputil.WriteError(w, http.StatusNotFound, "contract.not_found", "contract not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, contract)
}

func (h *Handler) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	if h.Events == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "event reader not configured")
		return
	}
	executionID := strings.TrimSpace(r.URL.Query().Get("execution_id"))
	if executionID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "events.missing_execution_id", "execution_id query param is required")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"execution_id": executionID, "records": h.Events.RecordsByExecution(r.Context(), executionID)})
}

func (h *Handler) handleGetRun(w http.ResponseWriter, r *http.Request) {
	if h.Runs == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "run reader not configured")
		return
	}
	executionID := strings.TrimPrefix(r.URL.Path, "/v1/foundation/runs/")
	if executionID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "foundation_run.missing_execution_id", "execution id is required")
		return
	}
	run, found, err := h.Runs.GetByExecutionID(r.Context(), executionID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "foundation_run.lookup_failed", err.Error())
		return
	}
	if !found {
		httputil.WriteError(w, http.StatusNotFound, "foundation_run.not_found", "foundation run not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, run)
}

func (h *Handler) handleResolveCapability(w http.ResponseWriter, r *http.Request) {
	if h.Registry == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "registry reader not configured")
		return
	}
	capabilityID := strings.TrimSpace(r.URL.Query().Get("capability_id"))
	contractVersion := strings.TrimSpace(r.URL.Query().Get("contract_version"))
	resultType := strings.TrimSpace(r.URL.Query().Get("result_type"))
	environment := strings.TrimSpace(r.URL.Query().Get("environment"))
	if environment == "" {
		environment = "dev"
	}
	if contractVersion == "" || resultType == "" {
		httputil.WriteError(w, http.StatusBadRequest, "registry.invalid_request", "contract_version and result_type are required")
		return
	}
	resolved, err := h.Registry.Resolve(r.Context(), registry.ResolutionRequest{
		CapabilityID:          capabilityID,
		ContractSchemaVersion: contractVersion,
		SupportedResultType:   resultType,
		Environment:           environment,
	})
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "registry.resolve_failed", err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, resolved)
}

func (h *Handler) handleGetApproval(w http.ResponseWriter, r *http.Request) {
	if h.Approvals == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "approval reader not configured")
		return
	}
	approvalRequestID := strings.TrimPrefix(r.URL.Path, "/v1/approvals/")
	if approvalRequestID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "approval.missing_id", "approval request id is required")
		return
	}
	request, found, err := h.Approvals.GetByID(r.Context(), approvalRequestID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "approval.lookup_failed", err.Error())
		return
	}
	if !found {
		httputil.WriteError(w, http.StatusNotFound, "approval.not_found", "approval request not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, request)
}

func (h *Handler) handleApprovalDecision(w http.ResponseWriter, r *http.Request) {
	if h.Approvals == nil || h.Runtime == nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "service.not_ready", "approval/runtime services not configured")
		return
	}
	approvalRequestID := strings.TrimPrefix(r.URL.Path, "/v1/approvals/")
	decisionState := approvals.State("")
	switch {
	case strings.HasSuffix(approvalRequestID, "/release"):
		decisionState = approvals.StateReleased
		approvalRequestID = strings.TrimSuffix(approvalRequestID, "/release")
	case strings.HasSuffix(approvalRequestID, "/reject"):
		decisionState = approvals.StateRejected
		approvalRequestID = strings.TrimSuffix(approvalRequestID, "/reject")
	case strings.HasSuffix(approvalRequestID, "/escalate"):
		decisionState = approvals.StateEscalated
		approvalRequestID = strings.TrimSuffix(approvalRequestID, "/escalate")
	default:
		httputil.WriteError(w, http.StatusBadRequest, "approval.invalid_action", "expected /v1/approvals/{id}/release|reject|escalate")
		return
	}
	approvalRequestID = strings.TrimSuffix(approvalRequestID, "/")
	if approvalRequestID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "approval.missing_id", "approval request id is required")
		return
	}
	var req approvalDecisionRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	if strings.TrimSpace(req.DecidedBySubjectID) == "" {
		httputil.WriteError(w, http.StatusBadRequest, "approval.missing_decider", "decided_by_subject_id is required")
		return
	}
	current, found, err := h.Approvals.GetByID(r.Context(), approvalRequestID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "approval.lookup_failed", err.Error())
		return
	}
	if !found {
		httputil.WriteError(w, http.StatusNotFound, "approval.not_found", "approval request not found")
		return
	}
	var execution runtime.ExecutionRecord
	if decisionState == approvals.StateReleased {
		execution, found, err = h.Runtime.GetExecution(r.Context(), current.ExecutionID)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "execution.lookup_failed", err.Error())
			return
		}
		if !found {
			httputil.WriteError(w, http.StatusNotFound, "execution.not_found", "execution not found")
			return
		}
		if current.SourceContractFingerprint != "" && execution.ContractFingerprint != current.SourceContractFingerprint {
			if h.Events != nil {
				_ = h.Events.Append(r.Context(), events.Record{
					EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
					EventType:           "approval.fingerprint_mismatch",
					TenantID:            current.TenantID,
					TraceID:             current.TraceID,
					ContractID:          current.ContractID,
					ContractFingerprint: execution.ContractFingerprint,
					ExecutionID:         current.ExecutionID,
					ApprovalRequestID:   current.ApprovalRequestID,
					OccurredAt:          time.Now().UTC(),
					Payload: map[string]any{
						"approved_contract_fingerprint": current.SourceContractFingerprint,
						"current_contract_fingerprint":  execution.ContractFingerprint,
					},
				})
			}
			httputil.WriteJSON(w, http.StatusConflict, map[string]any{
				"error": map[string]any{
					"code":    "approval.fingerprint_mismatch",
					"message": "approval fingerprint no longer matches current execution contract",
				},
				"approval":  current,
				"execution": execution,
			})
			return
		}
	}
	request, err := h.Approvals.Decide(r.Context(), approvalRequestID, approvals.Decision{
		State:               decisionState,
		DecidedBySubjectID:  req.DecidedBySubjectID,
		DecisionComment:     req.DecisionComment,
		DecisionReasonCodes: req.DecisionReasonCodes,
		DecidedAt:           time.Now().UTC(),
	})
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "approval.decision_failed", err.Error())
		return
	}
	if decisionState == approvals.StateReleased {
		execution, err = h.Runtime.UpdateExecutionState(r.Context(), request.ExecutionID, runtime.ExecutionStateExecutionReleased)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "execution.release_failed", err.Error())
			return
		}
	}
	if h.Events != nil {
		decisionEventType := map[approvals.State]string{
			approvals.StateReleased:  "approval.released",
			approvals.StateRejected:  "approval.rejected",
			approvals.StateEscalated: "approval.escalated",
		}[decisionState]
		_ = h.Events.Append(r.Context(), events.Record{
			EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
			EventType:           decisionEventType,
			TenantID:            request.TenantID,
			TraceID:             request.TraceID,
			ContractID:          request.ContractID,
			ContractFingerprint: request.SourceContractFingerprint,
			ExecutionID:         request.ExecutionID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"approval_request_id":   request.ApprovalRequestID,
				"decided_by_subject_id": request.DecidedBySubjectID,
				"decision_reason_codes": request.DecisionReasonCodes,
			},
		})
	}
	payload := map[string]any{"approval": request}
	if decisionState == approvals.StateReleased {
		payload["execution"] = execution
	}
	httputil.WriteJSON(w, http.StatusOK, payload)
}

func Warmup(ctx context.Context, orchestrator *foundation.FoundationOrchestrator) error {
	if orchestrator == nil {
		return errors.New("missing orchestrator")
	}
	_, err := orchestrator.Run(ctx, intent.IntentInput{
		RequestID:             "bootstrap-request-1",
		TenantID:              "tenant-demo",
		WorkspaceID:           "workspace-demo",
		UserID:                "user-demo",
		SessionID:             "session-demo",
		Objetivo:              "Compilar una intención gobernada inicial",
		Alcance:               "foundation-slice",
		TipoResultadoEsperado: intent.ResultTypePlan,
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	})
	return err
}
