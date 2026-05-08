package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/inspection"
)

var _ inspection.RecoveryStore = (*RecoveryStore)(nil)

type RecoveryStore struct {
	mu         sync.RWMutex
	candidates map[string]inspection.RecoveryActionCandidate
}

func NewRecoveryStore() *RecoveryStore {
	return &RecoveryStore{candidates: map[string]inspection.RecoveryActionCandidate{}}
}

func (s *RecoveryStore) Create(ctx context.Context, candidate inspection.RecoveryActionCandidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.RecoveryActionCandidateID]; exists {
		return errors.New("recovery candidate already exists")
	}
	s.candidates[candidate.RecoveryActionCandidateID] = candidate
	return nil
}

func (s *RecoveryStore) GetByID(ctx context.Context, id string) (inspection.RecoveryActionCandidate, bool, error) {
	select {
	case <-ctx.Done():
		return inspection.RecoveryActionCandidate{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[id]
	return candidate, found, nil
}

func (s *RecoveryStore) ListByExecution(ctx context.Context, executionID string) ([]inspection.RecoveryActionCandidate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
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

func (s *RecoveryStore) Update(ctx context.Context, candidate inspection.RecoveryActionCandidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.RecoveryActionCandidateID]; !exists {
		return errors.New("recovery candidate not found")
	}
	s.candidates[candidate.RecoveryActionCandidateID] = candidate
	return nil
}
