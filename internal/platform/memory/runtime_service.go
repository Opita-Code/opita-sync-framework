package memory

import (
	"errors"
	"sync"
	"time"

	"opita-sync-framework/internal/engine/runtime"
)

type RuntimeService struct {
	mu         sync.RWMutex
	executions map[string]runtime.ExecutionRecord
}

func NewRuntimeService() *RuntimeService {
	return &RuntimeService{executions: map[string]runtime.ExecutionRecord{}}
}

func (s *RuntimeService) CreateExecution(record runtime.ExecutionRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.executions[record.ExecutionID]; exists {
		return errors.New("execution already exists")
	}
	s.executions[record.ExecutionID] = record
	return nil
}

func (s *RuntimeService) GetExecution(executionID string) (runtime.ExecutionRecord, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.executions[executionID]
	return record, ok, nil
}

func (s *RuntimeService) UpdateExecutionState(executionID string, state runtime.ExecutionState) (runtime.ExecutionRecord, error) {
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
