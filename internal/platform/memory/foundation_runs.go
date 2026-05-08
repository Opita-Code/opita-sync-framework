package memory

import (
	"context"
	"sync"

	"opita-sync-framework/internal/engine/foundation"
)

var _ foundation.RunRepository = (*FoundationRunStore)(nil)

type FoundationRunStore struct {
	mu     sync.RWMutex
	byExec map[string]foundation.FoundationRunResult
}

func NewFoundationRunStore() *FoundationRunStore {
	return &FoundationRunStore{byExec: map[string]foundation.FoundationRunResult{}}
}

func (s *FoundationRunStore) Save(ctx context.Context, result foundation.FoundationRunResult) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byExec[result.Execution.ExecutionID] = result
	return nil
}

func (s *FoundationRunStore) GetByExecutionID(ctx context.Context, executionID string) (foundation.FoundationRunResult, bool, error) {
	select {
	case <-ctx.Done():
		return foundation.FoundationRunResult{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, found := s.byExec[executionID]
	return result, found, nil
}
