package foundation

import "errors"

var ErrMissingRuns = errors.New("foundation orchestrator requires run repository")

type RunRepository interface {
	Save(result FoundationRunResult) error
	GetByExecutionID(executionID string) (FoundationRunResult, bool, error)
}

func (o *FoundationOrchestrator) validateRuns() error {
	if o.Runs == nil {
		return ErrMissingRuns
	}
	return nil
}
