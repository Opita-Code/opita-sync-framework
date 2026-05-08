package events

import (
	"context"
	"time"
)

type Record struct {
	EventID             string
	EventType           string
	TenantID            string
	TraceID             string
	ContractID          string
	ContractFingerprint string
	ExecutionID         string
	PolicyDecisionID    string
	ConversationTurnID  string
	IntakeSessionID     string
	IntentCandidateID   string
	ProposalDraftID     string
	PatchsetCandidateID string
	PreviewCandidateID  string
	SimulationResultID  string
	ApprovalRequestID   string
	RecoveryActionID    string
	OccurredAt          time.Time
	Payload             map[string]any
}

type EventLog interface {
	Append(ctx context.Context, record Record) error
	Records(ctx context.Context) []Record
}
