package preview

import "time"

type Status string

const (
	StatusPreviewOK         Status = "preview_ok"
	StatusPreviewWarning    Status = "preview_warning"
	StatusPreviewBlocked    Status = "preview_blocked"
	StatusPreviewIncomplete Status = "preview_incomplete"
)

type Candidate struct {
	PreviewCandidateID string
	TenantID           string
	SessionID          string
	SubjectID          string
	ProposalDraftID    string
	ContractID         string
	ExecutionID        string
	PatchsetRef        string
	HumanDiffRef       string
	MaterialDiffRef    string
	MaterialDiffHash   string
	PreviewScope       string
	State              Status
	CreatedAt          time.Time
}

type SimulationFamily string

const (
	SimulationPolicy         SimulationFamily = "policy"
	SimulationApproval       SimulationFamily = "approval"
	SimulationClassification SimulationFamily = "classification"
	SimulationRisk           SimulationFamily = "risk"
)

type Result struct {
	SimulationResultID string
	PreviewCandidateID string
	Family             SimulationFamily
	Status             Status
	ReasonCodes        []string
	InputsRefs         []string
	OutputsSummary     string
	ConfidenceLevel    string
	CreatedAt          time.Time
}

type Service interface {
	CreateCandidate(candidate Candidate) error
	GetCandidate(previewCandidateID string) (Candidate, bool, error)
	SaveResult(result Result) error
	ListResults(previewCandidateID string) ([]Result, error)
}
