package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/inspection"
)

type RecoveryStore struct {
	mu         sync.RWMutex
	candidates map[string]inspection.RecoveryActionCandidate
}

func NewRecoveryStore() *RecoveryStore {
	return &RecoveryStore{candidates: map[string]inspection.RecoveryActionCandidate{}}
}

func (s *RecoveryStore) Create(candidate inspection.RecoveryActionCandidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.RecoveryActionCandidateID]; exists {
		return errors.New("recovery candidate already exists")
	}
	s.candidates[candidate.RecoveryActionCandidateID] = candidate
	return nil
}

func (s *RecoveryStore) GetByID(id string) (inspection.RecoveryActionCandidate, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[id]
	return candidate, found, nil
}

func (s *RecoveryStore) ListByExecution(executionID string) ([]inspection.RecoveryActionCandidate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]inspection.RecoveryActionCandidate, 0)
	for _, candidate := range s.candidates {
		if candidate.ExecutionID == executionID {
			out = append(out, candidate)
		}
	}
	return out, nil
}

func (s *RecoveryStore) Update(candidate inspection.RecoveryActionCandidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.RecoveryActionCandidateID]; !exists {
		return errors.New("recovery candidate not found")
	}
	s.candidates[candidate.RecoveryActionCandidateID] = candidate
	return nil
}
