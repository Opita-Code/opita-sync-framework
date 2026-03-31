package inspection

import "time"

type ExecutionInspectionView struct {
	InspectionViewID      string    `json:"inspection_view_id"`
	ExecutionID           string    `json:"execution_id"`
	TenantID              string    `json:"tenant_id"`
	TraceID               string    `json:"trace_id"`
	ContractID            string    `json:"contract_id"`
	CurrentRuntimeState   string    `json:"current_runtime_state"`
	CurrentOutcomeState   string    `json:"current_outcome_state"`
	PolicyDecisionRefs    []string  `json:"policy_decision_refs"`
	ApprovalRequestRefs   []string  `json:"approval_request_refs"`
	EventRefs             []string  `json:"event_refs"`
	ConversationTurnRefs  []string  `json:"conversation_turn_refs"`
	IntakeSessionRefs     []string  `json:"intake_session_refs"`
	IntentCandidateRefs   []string  `json:"intent_candidate_refs"`
	ProposalDraftRefs     []string  `json:"proposal_draft_refs"`
	PatchsetCandidateRefs []string  `json:"patchset_candidate_refs"`
	PreviewCandidateRefs  []string  `json:"preview_candidate_refs"`
	SimulationResultRefs  []string  `json:"simulation_result_refs"`
	ResolvedCapability    string    `json:"resolved_capability_ref"`
	ResolvedBinding       string    `json:"resolved_binding_ref"`
	ResolvedProvider      string    `json:"resolved_provider_ref"`
	OperatorSummary       string    `json:"operator_summary"`
	OperatorFindings      []string  `json:"operator_findings"`
	GeneratedAt           time.Time `json:"generated_at"`
}

type RecoveryAction string

const (
	RecoveryRetryTechnicalStep     RecoveryAction = "retry_technical_step"
	RecoveryResumeAfterApproval    RecoveryAction = "resume_after_approval"
	RecoveryRequestManualReview    RecoveryAction = "request_human_review"
	RecoveryRequestManualComp      RecoveryAction = "request_manual_compensation"
	RecoveryAcknowledgeUnknown     RecoveryAction = "acknowledge_unknown_outcome"
	RecoveryEscalateForHumanReview RecoveryAction = "escalate_for_human_review"
)

type RecoveryCandidateState string

const (
	RecoveryCandidatePending  RecoveryCandidateState = "pending"
	RecoveryCandidateExecuted RecoveryCandidateState = "executed"
	RecoveryCandidateBlocked  RecoveryCandidateState = "blocked"
)

type RecoveryActionCandidate struct {
	RecoveryActionCandidateID  string                 `json:"recovery_action_candidate_id"`
	ExecutionID                string                 `json:"execution_id"`
	RequestedAction            RecoveryAction         `json:"requested_action"`
	RequestedBySubjectID       string                 `json:"requested_by_subject_id"`
	CurrentRuntimeState        string                 `json:"current_runtime_state"`
	ApprovalRequestID          string                 `json:"approval_request_id,omitempty"`
	PreconditionsRefs          []string               `json:"preconditions_refs"`
	BlockingConstraints        []string               `json:"blocking_constraints"`
	ReasonCodes                []string               `json:"reason_codes"`
	ReadyForExecution          bool                   `json:"ready_for_execution"`
	RequiresAdditionalApproval bool                   `json:"requires_additional_approval"`
	State                      RecoveryCandidateState `json:"state"`
	CreatedAt                  time.Time              `json:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at"`
}

type RecoveryStore interface {
	Create(candidate RecoveryActionCandidate) error
	GetByID(recoveryActionCandidateID string) (RecoveryActionCandidate, bool, error)
	ListByExecution(executionID string) ([]RecoveryActionCandidate, error)
	Update(candidate RecoveryActionCandidate) error
}
