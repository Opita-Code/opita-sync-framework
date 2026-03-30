package proposal

import "time"

type DraftState string

const (
	DraftStateOpen             DraftState = "draft_open"
	DraftStateRefining         DraftState = "draft_refining"
	DraftStateAwaitingRevision DraftState = "awaiting_revision"
	DraftStateReadyForPreview  DraftState = "ready_for_preview"
	DraftStatePreviewBlocked   DraftState = "preview_blocked"
	DraftStateReadyForApply    DraftState = "ready_for_apply_candidate"
	DraftStateRejected         DraftState = "rejected"
	DraftStateCancelled        DraftState = "cancelled"
	DraftStateClosed           DraftState = "closed"
)

type Draft struct {
	ProposalDraftID  string
	TenantID         string
	SessionID        string
	SubjectID        string
	SourceIntentRefs []string
	Title            string
	Summary          string
	ArtifactsInScope []string
	ProposedChanges  []string
	Constraints      []string
	Assumptions      []string
	OpenQuestions    []string
	CurrentState     DraftState
	ConfidenceLevel  string
	HumanDiffRef     string
	MaterialDiffRef  string
	CreatedAt        time.Time
}

type PatchsetCandidate struct {
	PatchsetCandidateID            string
	ProposalDraftID                string
	TargetArtifacts                []string
	MaterialOperations             []string
	MaterialDiffHash               string
	PolicyPreviewInputsRef         string
	ApprovalPreviewInputsRef       string
	ClassificationPreviewInputsRef string
	ReadyForPreview                bool
	ReadyForApplyCandidate         bool
	CreatedAt                      time.Time
}

type Service interface {
	CreateDraft(draft Draft) error
	GetDraft(proposalDraftID string) (Draft, bool, error)
	SavePatchset(candidate PatchsetCandidate) error
	GetPatchset(patchsetCandidateID string) (PatchsetCandidate, bool, error)
}
