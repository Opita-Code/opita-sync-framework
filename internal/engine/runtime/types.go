package runtime

import "time"

type ExecutionState string

const (
	ExecutionStateCreated             ExecutionState = "created"
	ExecutionStateEligibilityCheck    ExecutionState = "eligibility_check"
	ExecutionStateAwaitingApproval    ExecutionState = "awaiting_approval"
	ExecutionStateExecutionReleased   ExecutionState = "execution_released"
	ExecutionStateExecuting           ExecutionState = "executing"
	ExecutionStateExecutionCompleted  ExecutionState = "execution_completed"
	ExecutionStateApplicationReleased ExecutionState = "application_released"
	ExecutionStateApplying            ExecutionState = "applying"
	ExecutionStateApplicationDone     ExecutionState = "application_completed"
	ExecutionStateCompensationPending ExecutionState = "compensation_pending"
	ExecutionStateCompensated         ExecutionState = "compensated"
	ExecutionStateUnknownOutcome      ExecutionState = "unknown_outcome"
	ExecutionStateBlocked             ExecutionState = "blocked"
	ExecutionStateFailed              ExecutionState = "failed"
	ExecutionStateClosed              ExecutionState = "closed"
)

type ExecutionRecord struct {
	ExecutionID         string
	TenantID            string
	ContractID          string
	ContractFingerprint string
	TraceID             string
	State               ExecutionState
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type RuntimeService interface {
	CreateExecution(record ExecutionRecord) error
	GetExecution(executionID string) (ExecutionRecord, bool, error)
	UpdateExecutionState(executionID string, state ExecutionState) (ExecutionRecord, error)
}
