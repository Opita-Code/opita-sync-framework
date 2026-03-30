package previewservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/preview"
	"opita-sync-framework/internal/engine/simulation"
)

type Handler struct {
	PreviewStore preview.Service
	Simulation   *simulation.Service
	Events       EventWriter
}

type EventWriter interface {
	Append(record events.Record) error
}

type createPreviewRequest struct {
	TenantID         string `json:"tenant_id"`
	SessionID        string `json:"session_id"`
	SubjectID        string `json:"subject_id"`
	TraceID          string `json:"trace_id,omitempty"`
	ProposalDraftID  string `json:"proposal_draft_id,omitempty"`
	ContractID       string `json:"contract_id"`
	ExecutionID      string `json:"execution_id"`
	PatchsetRef      string `json:"patchset_ref"`
	HumanDiffRef     string `json:"human_diff_ref,omitempty"`
	MaterialDiffRef  string `json:"material_diff_ref,omitempty"`
	MaterialDiffHash string `json:"material_diff_hash"`
	PreviewScope     string `json:"preview_scope"`
}

type readablePreview struct {
	PreviewCandidateID string                  `json:"preview_candidate_id"`
	PreviewState       preview.Status          `json:"preview_state"`
	PredictionBoundary string                  `json:"prediction_boundary"`
	PatchsetRef        string                  `json:"patchset_ref"`
	Diff               readableDiff            `json:"diff"`
	Approvals          readableApprovalSummary `json:"approvals"`
	Risk               readableSimulationCard  `json:"risk"`
	Policy             readableSimulationCard  `json:"policy"`
	Classification     readableSimulationCard  `json:"classification"`
	PromotionReadiness string                  `json:"promotion_readiness"`
	ReasonCodes        []string                `json:"reason_codes"`
}

type readableDiff struct {
	HumanDiffRef     string `json:"human_diff_ref,omitempty"`
	MaterialDiffRef  string `json:"material_diff_ref,omitempty"`
	MaterialDiffHash string `json:"material_diff_hash"`
	Summary          string `json:"summary"`
}

type readableApprovalSummary struct {
	Status          preview.Status `json:"status"`
	Required        bool           `json:"required"`
	Summary         string         `json:"summary"`
	ReasonCodes     []string       `json:"reason_codes"`
	ConfidenceLevel string         `json:"confidence_level"`
}

type readableSimulationCard struct {
	Status          preview.Status `json:"status"`
	Summary         string         `json:"summary"`
	ReasonCodes     []string       `json:"reason_codes"`
	ConfidenceLevel string         `json:"confidence_level"`
}

func NewHandler(store preview.Service, simulationService *simulation.Service, eventWriter EventWriter) *Handler {
	return &Handler{PreviewStore: store, Simulation: simulationService, Events: eventWriter}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/previews", h.handleCreatePreview)
	mux.HandleFunc("GET /v1/previews/", h.handleGetPreview)
	mux.HandleFunc("GET /v1/simulations", h.handleListSimulations)
	mux.HandleFunc("GET /v1/readable-previews/", h.handleGetReadablePreview)
	return mux
}

func (h *Handler) handleCreatePreview(w http.ResponseWriter, r *http.Request) {
	var req createPreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "request.invalid_json", "message": err.Error()})
		return
	}
	candidate := preview.Candidate{
		PreviewCandidateID: fmt.Sprintf("preview-%d", time.Now().UTC().UnixNano()),
		TenantID:           req.TenantID,
		SessionID:          req.SessionID,
		SubjectID:          req.SubjectID,
		ProposalDraftID:    req.ProposalDraftID,
		ContractID:         req.ContractID,
		ExecutionID:        req.ExecutionID,
		PatchsetRef:        req.PatchsetRef,
		HumanDiffRef:       req.HumanDiffRef,
		MaterialDiffRef:    req.MaterialDiffRef,
		MaterialDiffHash:   req.MaterialDiffHash,
		PreviewScope:       req.PreviewScope,
		State:              preview.StatusPreviewWarning,
		CreatedAt:          time.Now().UTC(),
	}
	if err := h.PreviewStore.CreateCandidate(candidate); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "preview.create_failed", "message": err.Error()})
		return
	}
	results, err := h.Simulation.RunAll(candidate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "preview.simulation_failed", "message": err.Error()})
		return
	}
	for _, result := range results {
		_ = h.PreviewStore.SaveResult(result)
		_ = h.appendEvent(events.Record{
			EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()),
			EventType:           "preview.simulation_recorded",
			TenantID:            candidate.TenantID,
			TraceID:             req.TraceID,
			ContractID:          candidate.ContractID,
			ExecutionID:         candidate.ExecutionID,
			ProposalDraftID:     req.ProposalDraftID,
			PatchsetCandidateID: candidate.PatchsetRef,
			PreviewCandidateID:  candidate.PreviewCandidateID,
			SimulationResultID:  result.SimulationResultID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"preview_candidate_id": candidate.PreviewCandidateID,
				"simulation_result_id": result.SimulationResultID,
				"family":               result.Family,
				"status":               result.Status,
			},
		})
	}
	_ = h.appendEvent(events.Record{
		EventID:             fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()+1),
		EventType:           "preview.created",
		TenantID:            candidate.TenantID,
		TraceID:             req.TraceID,
		ContractID:          candidate.ContractID,
		ExecutionID:         candidate.ExecutionID,
		ProposalDraftID:     req.ProposalDraftID,
		PatchsetCandidateID: candidate.PatchsetRef,
		PreviewCandidateID:  candidate.PreviewCandidateID,
		OccurredAt:          time.Now().UTC(),
		Payload: map[string]any{
			"preview_candidate_id": candidate.PreviewCandidateID,
			"preview_scope":        candidate.PreviewScope,
		},
	})
	writeJSON(w, http.StatusCreated, map[string]any{"preview_candidate": candidate, "simulation_results": results})
}

func (h *Handler) handleGetPreview(w http.ResponseWriter, r *http.Request) {
	previewID := strings.TrimPrefix(r.URL.Path, "/v1/previews/")
	if previewID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "preview.missing_id"})
		return
	}
	candidate, found, err := h.PreviewStore.GetCandidate(previewID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "preview.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "preview.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, candidate)
}

func (h *Handler) handleListSimulations(w http.ResponseWriter, r *http.Request) {
	previewID := strings.TrimSpace(r.URL.Query().Get("preview_id"))
	if previewID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "simulations.missing_preview_id"})
		return
	}
	results, err := h.PreviewStore.ListResults(previewID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "simulations.lookup_failed", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"preview_id": previewID, "results": results})
}

func (h *Handler) handleGetReadablePreview(w http.ResponseWriter, r *http.Request) {
	previewID := strings.TrimPrefix(r.URL.Path, "/v1/readable-previews/")
	previewID = strings.TrimSuffix(previewID, "/")
	if previewID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "readable_preview.missing_id"})
		return
	}
	candidate, found, err := h.PreviewStore.GetCandidate(previewID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "readable_preview.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "readable_preview.not_found"})
		return
	}
	results, err := h.PreviewStore.ListResults(previewID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "readable_preview.results_lookup_failed", "message": err.Error()})
		return
	}
	readable := buildReadablePreview(candidate, results)
	writeJSON(w, http.StatusOK, readable)
}

func buildReadablePreview(candidate preview.Candidate, results []preview.Result) readablePreview {
	policyCard := readableSimulationCard{Status: preview.StatusPreviewIncomplete, Summary: "policy preview unavailable"}
	approvalCard := readableApprovalSummary{Status: preview.StatusPreviewIncomplete, Summary: "approval preview unavailable"}
	classificationCard := readableSimulationCard{Status: preview.StatusPreviewIncomplete, Summary: "classification preview unavailable"}
	riskCard := readableSimulationCard{Status: preview.StatusPreviewIncomplete, Summary: "risk preview unavailable"}
	allReasonCodes := make([]string, 0)
	promotionReadiness := "review_before_promotion"
	for _, result := range results {
		allReasonCodes = append(allReasonCodes, result.ReasonCodes...)
		switch result.Family {
		case preview.SimulationPolicy:
			policyCard = readableSimulationCard{Status: result.Status, Summary: result.OutputsSummary, ReasonCodes: result.ReasonCodes, ConfidenceLevel: result.ConfidenceLevel}
		case preview.SimulationApproval:
			approvalCard = readableApprovalSummary{Status: result.Status, Required: result.Status != preview.StatusPreviewOK, Summary: result.OutputsSummary, ReasonCodes: result.ReasonCodes, ConfidenceLevel: result.ConfidenceLevel}
		case preview.SimulationClassification:
			classificationCard = readableSimulationCard{Status: result.Status, Summary: result.OutputsSummary, ReasonCodes: result.ReasonCodes, ConfidenceLevel: result.ConfidenceLevel}
		case preview.SimulationRisk:
			riskCard = readableSimulationCard{Status: result.Status, Summary: result.OutputsSummary, ReasonCodes: result.ReasonCodes, ConfidenceLevel: result.ConfidenceLevel}
		}
		if result.Status == preview.StatusPreviewBlocked {
			promotionReadiness = "blocked_by_preview"
		}
	}
	if promotionReadiness != "blocked_by_preview" && (approvalCard.Required || riskCard.Status == preview.StatusPreviewWarning || policyCard.Status == preview.StatusPreviewWarning) {
		promotionReadiness = "requires_operator_review"
	}
	return readablePreview{
		PreviewCandidateID: candidate.PreviewCandidateID,
		PreviewState:       candidate.State,
		PredictionBoundary: "simulation_preview_not_kernel_truth",
		PatchsetRef:        candidate.PatchsetRef,
		Diff: readableDiff{
			HumanDiffRef:     candidate.HumanDiffRef,
			MaterialDiffRef:  candidate.MaterialDiffRef,
			MaterialDiffHash: candidate.MaterialDiffHash,
			Summary:          summarizeDiff(candidate),
		},
		Approvals:          approvalCard,
		Risk:               riskCard,
		Policy:             policyCard,
		Classification:     classificationCard,
		PromotionReadiness: promotionReadiness,
		ReasonCodes:        uniqueStrings(allReasonCodes),
	}
}

func summarizeDiff(candidate preview.Candidate) string {
	if candidate.HumanDiffRef != "" || candidate.MaterialDiffRef != "" {
		return fmt.Sprintf("preview uses human diff %s and material diff %s", candidate.HumanDiffRef, candidate.MaterialDiffRef)
	}
	return fmt.Sprintf("preview uses patchset %s and material hash %s", candidate.PatchsetRef, candidate.MaterialDiffHash)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) appendEvent(record events.Record) error {
	if h.Events == nil {
		return nil
	}
	return h.Events.Append(record)
}
