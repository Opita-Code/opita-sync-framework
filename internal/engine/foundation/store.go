package foundation

import (
	"context"
	"errors"
)

var ErrMissingRuns = errors.New("foundation orchestrator requires run repository")

type RunRepository interface {
	Save(ctx context.Context, result FoundationRunResult) error
	GetByExecutionID(ctx context.Context, executionID string) (FoundationRunResult, bool, error)
}

func (o *FoundationOrchestrator) validateRuns() error {
	if o.Runs == nil {
		return ErrMissingRuns
	}
	return nil
}
