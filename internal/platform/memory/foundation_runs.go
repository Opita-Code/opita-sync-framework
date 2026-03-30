package memory

import (
	"sync"

	"opita-sync-framework/internal/engine/foundation"
)

type FoundationRunStore struct {
	mu     sync.RWMutex
	byExec map[string]foundation.FoundationRunResult
}

func NewFoundationRunStore() *FoundationRunStore {
	return &FoundationRunStore{byExec: map[string]foundation.FoundationRunResult{}}
}

func (s *FoundationRunStore) Save(result foundation.FoundationRunResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byExec[result.Execution.ExecutionID] = result
	return nil
}

func (s *FoundationRunStore) GetByExecutionID(executionID string) (foundation.FoundationRunResult, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, found := s.byExec[executionID]
	return result, found, nil
}
