package maintenance

import "time"

type SemanticDebugView struct {
	DebugViewID      string    `json:"debug_view_id"`
	TenantID         string    `json:"tenant_id"`
	TraceID          string    `json:"trace_id"`
	SourceRefs       []string  `json:"source_refs"`
	ArtifactRefs     []string  `json:"artifact_refs"`
	ContractRefs     []string  `json:"contract_refs"`
	ExecutionRefs    []string  `json:"execution_refs"`
	PolicyRefs       []string  `json:"policy_refs"`
	EventRefs        []string  `json:"event_refs"`
	SemanticSummary  string    `json:"semantic_summary"`
	ExplanationChain []string  `json:"explanation_chain"`
	Ambiguities      []string  `json:"ambiguities"`
	EvidenceRefs     []string  `json:"evidence_refs"`
	GeneratedAt      time.Time `json:"generated_at"`
}

type ActionType string

const (
	ActionRecompileContract  ActionType = "request_recompile_contract"
	ActionRebuildPreview     ActionType = "request_rebuild_preview"
	ActionRefreshRegistry    ActionType = "request_refresh_registry_resolution"
	ActionReprojectEventView ActionType = "request_reproject_event_view"
	ActionRequestHumanReview ActionType = "request_human_review"
)

type ActionCandidate struct {
	MaintenanceActionCandidateID string     `json:"maintenance_action_candidate_id"`
	TenantID                     string     `json:"tenant_id"`
	RequestedBySubjectID         string     `json:"requested_by_subject_id"`
	ActionType                   ActionType `json:"action_type"`
	TargetRefs                   []string   `json:"target_refs"`
	PreconditionsRefs            []string   `json:"preconditions_refs"`
	GovernanceRequirements       []string   `json:"governance_requirements"`
	ReasonCodes                  []string   `json:"reason_codes"`
	ReadyForExecution            bool       `json:"ready_for_execution"`
	RequiresHumanReview          bool       `json:"requires_human_review"`
	CreatedAt                    time.Time  `json:"created_at"`
}

type Service interface {
	Create(candidate ActionCandidate) error
	GetByID(id string) (ActionCandidate, bool, error)
}
