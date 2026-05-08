package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/maintenance"
)

var _ maintenance.Service = (*MaintenanceStore)(nil)

type MaintenanceStore struct {
	mu         sync.RWMutex
	candidates map[string]maintenance.ActionCandidate
}

func NewMaintenanceStore() *MaintenanceStore {
	return &MaintenanceStore{candidates: map[string]maintenance.ActionCandidate{}}
}

func (s *MaintenanceStore) Create(ctx context.Context, candidate maintenance.ActionCandidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.MaintenanceActionCandidateID]; exists {
		return errors.New("maintenance action candidate already exists")
	}
	s.candidates[candidate.MaintenanceActionCandidateID] = candidate
	return nil
}

func (s *MaintenanceStore) GetByID(ctx context.Context, id string) (maintenance.ActionCandidate, bool, error) {
	select {
	case <-ctx.Done():
		return maintenance.ActionCandidate{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[id]
	return candidate, found, nil
}
