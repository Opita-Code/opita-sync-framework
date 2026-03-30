package surfaceservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/intake"
	"opita-sync-framework/internal/engine/proposal"
)

type Handler struct {
	Intake   intake.Service
	Proposal proposal.Service
	Events   EventWriter
}

type workspaceSummary struct {
	WorkspaceState string            `json:"workspace_state"`
	ChatBoundary   string            `json:"chat_boundary"`
	NextGates      []string          `json:"next_gates"`
	ArtifactRefs   map[string]string `json:"artifact_refs"`
	Intake         *intakeCard       `json:"intake,omitempty"`
	Proposal       *proposalCard     `json:"proposal,omitempty"`
	Patchset       *patchsetCard     `json:"patchset,omitempty"`
}

type intakeCard struct {
	SessionState     intake.SessionState `json:"session_state"`
	LastDecision     intake.Decision     `json:"last_decision"`
	AmbiguityLevel   string              `json:"ambiguity_level"`
	OpenQuestions    []string            `json:"open_questions"`
	ReadyForProposal bool                `json:"ready_for_proposal"`
	ReadyForIntent   bool                `json:"ready_for_intent"`
	Summary          string              `json:"summary"`
}

type proposalCard struct {
	State           proposal.DraftState `json:"state"`
	Title           string              `json:"title"`
	Summary         string              `json:"summary"`
	ConfidenceLevel string              `json:"confidence_level"`
	HumanDiffRef    string              `json:"human_diff_ref,omitempty"`
	MaterialDiffRef string              `json:"material_diff_ref,omitempty"`
	OpenQuestions   []string            `json:"open_questions"`
	SummaryText     string              `json:"summary_text"`
}

type patchsetCard struct {
	ReadyForPreview  bool     `json:"ready_for_preview"`
	ReadyForApply    bool     `json:"ready_for_apply_candidate"`
	TargetArtifacts  []string `json:"target_artifacts"`
	MaterialOps      []string `json:"material_operations"`
	MaterialDiffHash string   `json:"material_diff_hash"`
	Summary          string   `json:"summary"`
}

type EventWriter interface {
	Append(record events.Record) error
}

type createIntakeRequest struct {
	TenantID  string `json:"tenant_id"`
	SessionID string `json:"session_id"`
	SubjectID string `json:"subject_id"`
	RawText   string `json:"raw_text"`
	TraceID   string `json:"trace_id,omitempty"`
}

type createProposalRequest struct {
	TenantID         string   `json:"tenant_id"`
	SessionID        string   `json:"session_id"`
	SubjectID        string   `json:"subject_id"`
	TraceID          string   `json:"trace_id,omitempty"`
	IntakeSessionID  string   `json:"intake_session_id,omitempty"`
	SourceIntentRefs []string `json:"source_intent_refs"`
	Title            string   `json:"title"`
	Summary          string   `json:"summary"`
	ArtifactsInScope []string `json:"artifacts_in_scope,omitempty"`
	ProposedChanges  []string `json:"proposed_changes,omitempty"`
	Constraints      []string `json:"constraints,omitempty"`
	Assumptions      []string `json:"assumptions,omitempty"`
	OpenQuestions    []string `json:"open_questions,omitempty"`
	HumanDiffRef     string   `json:"human_diff_ref,omitempty"`
	MaterialDiffRef  string   `json:"material_diff_ref,omitempty"`
}

type createPatchsetRequest struct {
	TraceID                        string   `json:"trace_id,omitempty"`
	ProposalDraftID                string   `json:"proposal_draft_id"`
	TargetArtifacts                []string `json:"target_artifacts"`
	MaterialOperations             []string `json:"material_operations"`
	MaterialDiffHash               string   `json:"material_diff_hash"`
	PolicyPreviewInputsRef         string   `json:"policy_preview_inputs_ref"`
	ApprovalPreviewInputsRef       string   `json:"approval_preview_inputs_ref"`
	ClassificationPreviewInputsRef string   `json:"classification_preview_inputs_ref"`
}

func NewHandler(intakeStore intake.Service, proposalStore proposal.Service, eventWriter EventWriter) *Handler {
	return &Handler{Intake: intakeStore, Proposal: proposalStore, Events: eventWriter}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/intake/turns", h.handleCreateIntakeTurn)
	mux.HandleFunc("GET /v1/intake/sessions/", h.handleGetIntakeSession)
	mux.HandleFunc("GET /v1/intake/candidates/", h.handleGetIntentCandidate)
	mux.HandleFunc("POST /v1/proposals", h.handleCreateProposal)
	mux.HandleFunc("GET /v1/proposals/", h.handleGetProposal)
	mux.HandleFunc("POST /v1/patchsets", h.handleCreatePatchset)
	mux.HandleFunc("GET /v1/patchsets/", h.handleGetPatchset)
	mux.HandleFunc("GET /v1/workspaces/intake-proposal", h.handleGetIntakeProposalWorkspace)
	return mux
}

func (h *Handler) handleCreateIntakeTurn(w http.ResponseWriter, r *http.Request) {
	var req createIntakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "intake.invalid_json", "message": err.Error()})
		return
	}
	now := time.Now().UTC()
	turnID := fmt.Sprintf("turn-%d", now.UnixNano())
	intakeSessionID := fmt.Sprintf("intake-%d", now.UnixNano()+1)
	intentCandidateID := fmt.Sprintf("intent-%d", now.UnixNano()+2)
	decision := intake.DecisionContinueFreeChat
	state := intake.SessionStateFreeChat
	readyForProposal := false
	readyForIntent := false
	openQuestions := []string{}
	constraints := []string{}
	if strings.TrimSpace(req.RawText) != "" {
		decision = intake.DecisionEmitProposalDraft
		state = intake.SessionStateProposalReady
		readyForProposal = true
		readyForIntent = true
	}
	turn := intake.ConversationTurn{
		ConversationTurnID: turnID,
		SessionID:          req.SessionID,
		TenantID:           req.TenantID,
		SubjectID:          req.SubjectID,
		MessageRole:        "user",
		RawText:            req.RawText,
		Timestamp:          now,
		TurnClassification: "intent_signal",
		TraceID:            req.TraceID,
	}
	session := intake.Session{
		IntakeSessionID: intakeSessionID,
		SessionID:       req.SessionID,
		TenantID:        req.TenantID,
		SubjectID:       req.SubjectID,
		CurrentState:    state,
		OpenQuestions:   openQuestions,
		ResolvedFacts:   []string{req.RawText},
		AmbiguityLevel:  "tolerable",
		LastDecision:    decision,
		TraceID:         req.TraceID,
		UpdatedAt:       now,
	}
	candidate := intake.IntentCandidate{
		IntentCandidateID:     intentCandidateID,
		SourceTurnIDs:         []string{turnID},
		ObjetivoCandidate:     req.RawText,
		AlcanceCandidate:      "surface-default",
		ArtifactsCandidate:    []string{},
		ConstraintsCandidate:  constraints,
		Assumptions:           []string{"surface inferred a draft-capable intent from raw text"},
		OpenQuestions:         openQuestions,
		ConfidenceLevel:       "medium",
		ReadyForIntentInput:   readyForIntent,
		ReadyForProposalDraft: readyForProposal,
	}
	if err := h.Intake.CreateTurn(turn); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "intake.turn_create_failed", "message": err.Error()})
		return
	}
	if err := h.Intake.CreateSession(session); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "intake.session_create_failed", "message": err.Error()})
		return
	}
	if err := h.Intake.SaveIntentCandidate(candidate); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "intake.candidate_create_failed", "message": err.Error()})
		return
	}
	_ = h.appendEvent(events.Record{
		EventID:            fmt.Sprintf("event-%d", now.UnixNano()+3),
		EventType:          "intake.turn_recorded",
		TenantID:           req.TenantID,
		TraceID:            req.TraceID,
		ConversationTurnID: turnID,
		IntakeSessionID:    intakeSessionID,
		IntentCandidateID:  intentCandidateID,
		OccurredAt:         now,
		Payload: map[string]any{
			"conversation_turn_id": turnID,
			"intake_session_id":    intakeSessionID,
			"intent_candidate_id":  intentCandidateID,
			"last_decision":        decision,
		},
	})
	writeJSON(w, http.StatusCreated, map[string]any{"conversation_turn_id": turnID, "intake_session": session, "intent_candidate": candidate})
}

func (h *Handler) handleGetIntakeSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/intake/sessions/")
	session, found, err := h.Intake.GetSession(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "intake.session_lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "intake.session_not_found"})
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *Handler) handleGetIntentCandidate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/intake/candidates/")
	candidate, found, err := h.Intake.GetIntentCandidate(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "intake.candidate_lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "intake.candidate_not_found"})
		return
	}
	writeJSON(w, http.StatusOK, candidate)
}

func (h *Handler) handleCreateProposal(w http.ResponseWriter, r *http.Request) {
	var req createProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "proposal.invalid_json", "message": err.Error()})
		return
	}
	now := time.Now().UTC()
	draft := proposal.Draft{
		ProposalDraftID:  fmt.Sprintf("proposal-%d", now.UnixNano()),
		TenantID:         req.TenantID,
		SessionID:        req.SessionID,
		SubjectID:        req.SubjectID,
		SourceIntentRefs: req.SourceIntentRefs,
		Title:            req.Title,
		Summary:          req.Summary,
		ArtifactsInScope: req.ArtifactsInScope,
		ProposedChanges:  req.ProposedChanges,
		Constraints:      req.Constraints,
		Assumptions:      req.Assumptions,
		OpenQuestions:    req.OpenQuestions,
		CurrentState:     proposal.DraftStateReadyForPreview,
		ConfidenceLevel:  "medium",
		HumanDiffRef:     req.HumanDiffRef,
		MaterialDiffRef:  req.MaterialDiffRef,
		CreatedAt:        now,
	}
	if err := h.Proposal.CreateDraft(draft); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "proposal.create_failed", "message": err.Error()})
		return
	}
	_ = h.appendEvent(events.Record{
		EventID:           fmt.Sprintf("event-%d", now.UnixNano()+1),
		EventType:         "proposal.created",
		TenantID:          req.TenantID,
		TraceID:           req.TraceID,
		IntakeSessionID:   req.IntakeSessionID,
		IntentCandidateID: firstOrEmpty(req.SourceIntentRefs),
		ProposalDraftID:   draft.ProposalDraftID,
		OccurredAt:        now,
		Payload: map[string]any{
			"proposal_draft_id":  draft.ProposalDraftID,
			"source_intent_refs": req.SourceIntentRefs,
			"state":              draft.CurrentState,
		},
	})
	writeJSON(w, http.StatusCreated, draft)
}

func (h *Handler) handleGetProposal(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/proposals/")
	draft, found, err := h.Proposal.GetDraft(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "proposal.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "proposal.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, draft)
}

func (h *Handler) handleCreatePatchset(w http.ResponseWriter, r *http.Request) {
	var req createPatchsetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "patchset.invalid_json", "message": err.Error()})
		return
	}
	now := time.Now().UTC()
	patchset := proposal.PatchsetCandidate{
		PatchsetCandidateID:            fmt.Sprintf("patchset-%d", now.UnixNano()),
		ProposalDraftID:                req.ProposalDraftID,
		TargetArtifacts:                req.TargetArtifacts,
		MaterialOperations:             req.MaterialOperations,
		MaterialDiffHash:               req.MaterialDiffHash,
		PolicyPreviewInputsRef:         req.PolicyPreviewInputsRef,
		ApprovalPreviewInputsRef:       req.ApprovalPreviewInputsRef,
		ClassificationPreviewInputsRef: req.ClassificationPreviewInputsRef,
		ReadyForPreview:                true,
		ReadyForApplyCandidate:         false,
		CreatedAt:                      now,
	}
	if err := h.Proposal.SavePatchset(patchset); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "patchset.create_failed", "message": err.Error()})
		return
	}
	_ = h.appendEvent(events.Record{
		EventID:             fmt.Sprintf("event-%d", now.UnixNano()+1),
		EventType:           "patchset.created",
		TraceID:             req.TraceID,
		ProposalDraftID:     patchset.ProposalDraftID,
		PatchsetCandidateID: patchset.PatchsetCandidateID,
		OccurredAt:          now,
		Payload: map[string]any{
			"patchset_candidate_id": patchset.PatchsetCandidateID,
			"proposal_draft_id":     patchset.ProposalDraftID,
			"material_diff_hash":    patchset.MaterialDiffHash,
		},
	})
	writeJSON(w, http.StatusCreated, patchset)
}

func (h *Handler) appendEvent(record events.Record) error {
	if h.Events == nil {
		return nil
	}
	return h.Events.Append(record)
}

func firstOrEmpty(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (h *Handler) handleGetPatchset(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/patchsets/")
	patchset, found, err := h.Proposal.GetPatchset(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "patchset.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "patchset.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, patchset)
}

func (h *Handler) handleGetIntakeProposalWorkspace(w http.ResponseWriter, r *http.Request) {
	intakeSessionID := strings.TrimSpace(r.URL.Query().Get("intake_session_id"))
	intentCandidateID := strings.TrimSpace(r.URL.Query().Get("intent_candidate_id"))
	proposalDraftID := strings.TrimSpace(r.URL.Query().Get("proposal_draft_id"))
	patchsetCandidateID := strings.TrimSpace(r.URL.Query().Get("patchset_candidate_id"))
	if intakeSessionID == "" && intentCandidateID == "" && proposalDraftID == "" && patchsetCandidateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "workspace.missing_refs", "message": "at least one artifact ref is required"})
		return
	}

	summary := workspaceSummary{
		WorkspaceState: "intake_only",
		ChatBoundary:   "free_chat_never_applies_directly",
		NextGates:      []string{},
		ArtifactRefs:   map[string]string{},
	}

	var session intake.Session
	var candidate intake.IntentCandidate
	var draft proposal.Draft
	var patchset proposal.PatchsetCandidate
	var foundAny bool

	if intakeSessionID != "" {
		got, found, err := h.Intake.GetSession(intakeSessionID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.intake_session_lookup_failed", "message": err.Error()})
			return
		}
		if !found {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "workspace.intake_session_not_found"})
			return
		}
		session = got
		foundAny = true
		summary.ArtifactRefs["intake_session_id"] = session.IntakeSessionID
		summary.Intake = &intakeCard{
			SessionState:     session.CurrentState,
			LastDecision:     session.LastDecision,
			AmbiguityLevel:   session.AmbiguityLevel,
			OpenQuestions:    session.OpenQuestions,
			ReadyForProposal: session.CurrentState == intake.SessionStateProposalReady,
			ReadyForIntent:   session.CurrentState == intake.SessionStateIntentReady || session.CurrentState == intake.SessionStateProposalReady,
			Summary:          summarizeIntake(session),
		}
		summary.WorkspaceState = string(session.CurrentState)
	}

	if intentCandidateID != "" {
		got, found, err := h.Intake.GetIntentCandidate(intentCandidateID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.intent_candidate_lookup_failed", "message": err.Error()})
			return
		}
		if !found {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "workspace.intent_candidate_not_found"})
			return
		}
		candidate = got
		foundAny = true
		summary.ArtifactRefs["intent_candidate_id"] = candidate.IntentCandidateID
		if summary.Intake == nil {
			summary.Intake = &intakeCard{}
		}
		summary.Intake.ReadyForProposal = candidate.ReadyForProposalDraft
		summary.Intake.ReadyForIntent = candidate.ReadyForIntentInput
		if summary.Intake.Summary == "" {
			summary.Intake.Summary = candidate.ObjetivoCandidate
		}
	}

	if proposalDraftID != "" {
		got, found, err := h.Proposal.GetDraft(proposalDraftID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.proposal_lookup_failed", "message": err.Error()})
			return
		}
		if !found {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "workspace.proposal_not_found"})
			return
		}
		draft = got
		foundAny = true
		summary.ArtifactRefs["proposal_draft_id"] = draft.ProposalDraftID
		summary.Proposal = &proposalCard{
			State:           draft.CurrentState,
			Title:           draft.Title,
			Summary:         draft.Summary,
			ConfidenceLevel: draft.ConfidenceLevel,
			HumanDiffRef:    draft.HumanDiffRef,
			MaterialDiffRef: draft.MaterialDiffRef,
			OpenQuestions:   draft.OpenQuestions,
			SummaryText:     summarizeProposal(draft),
		}
		summary.WorkspaceState = string(draft.CurrentState)
	}

	if patchsetCandidateID != "" {
		got, found, err := h.Proposal.GetPatchset(patchsetCandidateID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.patchset_lookup_failed", "message": err.Error()})
			return
		}
		if !found {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "workspace.patchset_not_found"})
			return
		}
		patchset = got
		foundAny = true
		summary.ArtifactRefs["patchset_candidate_id"] = patchset.PatchsetCandidateID
		summary.Patchset = &patchsetCard{
			ReadyForPreview:  patchset.ReadyForPreview,
			ReadyForApply:    patchset.ReadyForApplyCandidate,
			TargetArtifacts:  patchset.TargetArtifacts,
			MaterialOps:      patchset.MaterialOperations,
			MaterialDiffHash: patchset.MaterialDiffHash,
			Summary:          summarizePatchset(patchset),
		}
		if patchset.ReadyForPreview {
			summary.WorkspaceState = "ready_for_preview"
		} else {
			summary.WorkspaceState = "patchset_prepared"
		}
	}

	if !foundAny {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "workspace.not_found"})
		return
	}

	summary.NextGates = deriveNextGates(summary, session, candidate, draft, patchset)
	writeJSON(w, http.StatusOK, summary)
}

func summarizeIntake(session intake.Session) string {
	return fmt.Sprintf("intake session is %s with last decision %s", session.CurrentState, session.LastDecision)
}

func summarizeProposal(draft proposal.Draft) string {
	return fmt.Sprintf("proposal %s is %s and points to preview readiness", draft.ProposalDraftID, draft.CurrentState)
}

func summarizePatchset(patchset proposal.PatchsetCandidate) string {
	return fmt.Sprintf("patchset %s targets %d artifacts and ready_for_preview=%t", patchset.PatchsetCandidateID, len(patchset.TargetArtifacts), patchset.ReadyForPreview)
}

func deriveNextGates(summary workspaceSummary, session intake.Session, candidate intake.IntentCandidate, draft proposal.Draft, patchset proposal.PatchsetCandidate) []string {
	next := make([]string, 0, 4)
	if summary.Intake != nil && summary.Proposal == nil && (summary.Intake.ReadyForProposal || candidate.ReadyForProposalDraft) {
		next = append(next, "create_proposal_draft")
	}
	if summary.Proposal != nil && summary.Patchset == nil && draft.CurrentState == proposal.DraftStateReadyForPreview {
		next = append(next, "prepare_patchset_candidate")
	}
	if summary.Patchset != nil && patchset.ReadyForPreview {
		next = append(next, "create_preview_candidate")
	}
	if summary.Intake != nil && session.CurrentState == intake.SessionStateAwaitingClarification {
		next = append(next, "resolve_open_questions")
	}
	if len(next) == 0 {
		next = append(next, "review_current_artifacts")
	}
	return next
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
