package approvals

import "time"

type State string

const (
	StateAwaitingApproval State = "awaiting_approval"
	StateReleased         State = "released"
	StateEscalated        State = "escalated"
	StateRejected         State = "rejected"
)

type Decision struct {
	State               State
	DecidedBySubjectID  string
	DecisionComment     string
	DecisionReasonCodes []string
	DecidedAt           time.Time
}

type Request struct {
	ApprovalRequestID         string
	ExecutionID               string
	ContractID                string
	TenantID                  string
	TraceID                   string
	State                     State
	Mode                      string
	ReasonCodes               []string
	SourceContractFingerprint string
	DecidedBySubjectID        string
	DecisionComment           string
	DecisionReasonCodes       []string
	ReleasedAt                *time.Time
	RejectedAt                *time.Time
	EscalatedAt               *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type Service interface {
	Create(request Request) error
	GetByID(approvalRequestID string) (Request, bool, error)
	Decide(approvalRequestID string, decision Decision) (Request, error)
}
