package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"opita-sync-framework/internal/engine/runtime"
)

var _ runtime.RuntimeService = (*RuntimeService)(nil)

type RuntimeService struct {
	mu         sync.RWMutex
	executions map[string]runtime.ExecutionRecord
}

func NewRuntimeService() *RuntimeService {
	return &RuntimeService{executions: map[string]runtime.ExecutionRecord{}}
}

func (s *RuntimeService) CreateExecution(ctx context.Context, record runtime.ExecutionRecord) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.executions[record.ExecutionID]; exists {
		return errors.New("execution already exists")
	}
	s.executions[record.ExecutionID] = record
	return nil
}

func (s *RuntimeService) GetExecution(ctx context.Context, executionID string) (runtime.ExecutionRecord, bool, error) {
	select {
	case <-ctx.Done():
		return runtime.ExecutionRecord{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.executions[executionID]
	return record, ok, nil
}

func (s *RuntimeService) UpdateExecutionState(ctx context.Context, executionID string, state runtime.ExecutionState) (runtime.ExecutionRecord, error) {
	select {
	case <-ctx.Done():
		return runtime.ExecutionRecord{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.executions[executionID]
	if !ok {
		return runtime.ExecutionRecord{}, errors.New("execution not found")
	}
	record.State = state
	record.UpdatedAt = time.Now().UTC()
	s.executions[executionID] = record
	return record, nil
}
