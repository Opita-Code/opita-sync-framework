package foundation

import (
	"errors"

	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/runtime"
)

var (
	ErrMissingCompiler  = errors.New("foundation orchestrator requires compiler")
	ErrMissingPolicy    = errors.New("foundation orchestrator requires policy engine")
	ErrMissingRuntime   = errors.New("foundation orchestrator requires runtime service")
	ErrMissingEvents    = errors.New("foundation orchestrator requires event log")
	ErrMissingRegistry  = errors.New("foundation orchestrator requires registry resolver")
	ErrMissingApprovals = errors.New("foundation orchestrator requires approval service")
)

func (o *FoundationOrchestrator) Validate() error {
	if o.Compiler == nil {
		return ErrMissingCompiler
	}
	if o.Policy == nil {
		return ErrMissingPolicy
	}
	if o.Runtime == nil {
		return ErrMissingRuntime
	}
	if o.Events == nil {
		return ErrMissingEvents
	}
	if o.Registry == nil {
		return ErrMissingRegistry
	}
	if o.Approvals == nil {
		return ErrMissingApprovals
	}
	if err := o.validateRuns(); err != nil {
		return err
	}
	return nil
}

type eventLog interface {
	Append(record events.Record) error
}

type runtimeStore interface {
	CreateExecution(record runtime.ExecutionRecord) error
	GetExecution(executionID string) (runtime.ExecutionRecord, bool, error)
}
