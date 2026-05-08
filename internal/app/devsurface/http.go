package devsurface

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/maintenance"
	"opita-sync-framework/internal/httputil"
)

type RunReader interface {
	GetByExecutionID(ctx context.Context, executionID string) (foundation.FoundationRunResult, bool, error)
}

type MaintenanceStore interface {
	Create(ctx context.Context, candidate maintenance.ActionCandidate) error
	GetByID(ctx context.Context, id string) (maintenance.ActionCandidate, bool, error)
}

type Handler struct {
	Runs        RunReader
	Maintenance MaintenanceStore
	Events      EventWriter
}

type EventWriter interface {
	Append(ctx context.Context, record events.Record) error
}

type createMaintenanceRequest struct {
	TenantID             string   `json:"tenant_id"`
	RequestedBySubjectID string   `json:"requested_by_subject_id"`
	ActionType           string   `json:"action_type"`
	TargetRefs           []string `json:"target_refs"`
	ReasonCodes          []string `json:"reason_codes,omitempty"`
}

func NewHandler(runReader RunReader, maintenanceStore MaintenanceStore, eventWriter EventWriter) *Handler {
	return &Handler{Runs: runReader, Maintenance: maintenanceStore, Events: eventWriter}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/debug/semantic", h.handleSemanticDebug)
	mux.HandleFunc("POST /v1/maintenance-actions", h.handleCreateMaintenanceCandidate)
	mux.HandleFunc("GET /v1/maintenance-actions/", h.handleGetMaintenanceCandidate)
	return mux
}

func (h *Handler) handleSemanticDebug(w http.ResponseWriter, r *http.Request) {
	executionID := strings.TrimSpace(r.URL.Query().Get("execution_id"))
	if executionID == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": "debug.missing_execution_id"})
		return
	}
	run, found, err := h.Runs.GetByExecutionID(r.Context(), executionID)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "debug.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": "debug.not_found"})
		return
	}
	view := maintenance.SemanticDebugView{
		DebugViewID:      fmt.Sprintf("debug-%d", time.Now().UTC().UnixNano()),
		TenantID:         run.Execution.TenantID,
		TraceID:          run.Execution.TraceID,
		SourceRefs:       []string{run.Execution.ExecutionID},
		ArtifactRefs:     []string{run.Resolution.CapabilityManifestRef, run.Resolution.BundleDigest, run.Resolution.BindingID, run.Resolution.ProviderRef},
		ContractRefs:     []string{run.Contract.ContractID, run.Contract.Fingerprint},
		ExecutionRefs:    []string{run.Execution.ExecutionID},
		PolicyRefs:       []string{run.PolicyDecision.PolicyDecisionID},
		EventRefs:        []string{"canonical-event-log"},
		SemanticSummary:  fmt.Sprintf("execution %s linked contract %s and provider %s", run.Execution.ExecutionID, run.Contract.ContractID, run.Resolution.ProviderRef),
		ExplanationChain: []string{"intent compiled", "policy evaluated", "capability resolved", "execution created"},
		Ambiguities:      []string{},
		EvidenceRefs:     []string{run.Contract.ContractID, run.Execution.ExecutionID, run.PolicyDecision.PolicyDecisionID},
		GeneratedAt:      time.Now().UTC(),
	}
	_ = h.appendEvent(events.Record{
		EventID:     fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
		EventType:   "debug.semantic_view_generated",
		TenantID:    run.Execution.TenantID,
		TraceID:     run.Execution.TraceID,
		ContractID:  run.Execution.ContractID,
		ExecutionID: run.Execution.ExecutionID,
		OccurredAt:  time.Now().UTC(),
		Payload: map[string]any{
			"debug_view_id": view.DebugViewID,
		},
	})
	httputil.WriteJSON(w, http.StatusOK, view)
}

func (h *Handler) handleCreateMaintenanceCandidate(w http.ResponseWriter, r *http.Request) {
	var req createMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": "maintenance.invalid_json", "message": err.Error()})
		return
	}
	candidate := maintenance.ActionCandidate{
		MaintenanceActionCandidateID: fmt.Sprintf("maintenance-%d", time.Now().UTC().UnixNano()),
		TenantID:                     req.TenantID,
		RequestedBySubjectID:         req.RequestedBySubjectID,
		ActionType:                   maintenance.ActionType(req.ActionType),
		TargetRefs:                   req.TargetRefs,
		PreconditionsRefs:            req.TargetRefs,
		GovernanceRequirements:       []string{"governed-candidate-only"},
		ReasonCodes:                  req.ReasonCodes,
		ReadyForExecution:            false,
		RequiresHumanReview:          true,
		CreatedAt:                    time.Now().UTC(),
	}
	if err := h.Maintenance.Create(r.Context(), candidate); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": "maintenance.create_failed", "message": err.Error()})
		return
	}
	_ = h.appendEvent(events.Record{
		EventID:    fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
		EventType:  "maintenance.candidate_created",
		TenantID:   req.TenantID,
		OccurredAt: time.Now().UTC(),
		Payload: map[string]any{
			"maintenance_action_candidate_id": candidate.MaintenanceActionCandidateID,
			"action_type":                     candidate.ActionType,
		},
	})
	httputil.WriteJSON(w, http.StatusCreated, candidate)
}

func (h *Handler) handleGetMaintenanceCandidate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/maintenance-actions/")
	if id == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": "maintenance.missing_id"})
		return
	}
	candidate, found, err := h.Maintenance.GetByID(r.Context(), id)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "maintenance.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		httputil.WriteJSON(w, http.StatusNotFound, map[string]any{"error": "maintenance.not_found"})
		return
	}
	httputil.WriteJSON(w, http.StatusOK, candidate)
}

func (h *Handler) appendEvent(record events.Record) error {
	if h.Events == nil {
		return nil
	}
	return h.Events.Append(context.Background(), record)
}
