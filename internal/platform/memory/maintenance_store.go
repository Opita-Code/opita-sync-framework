package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/maintenance"
)

type MaintenanceStore struct {
	mu         sync.RWMutex
	candidates map[string]maintenance.ActionCandidate
}

func NewMaintenanceStore() *MaintenanceStore {
	return &MaintenanceStore{candidates: map[string]maintenance.ActionCandidate{}}
}

func (s *MaintenanceStore) Create(candidate maintenance.ActionCandidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.MaintenanceActionCandidateID]; exists {
		return errors.New("maintenance action candidate already exists")
	}
	s.candidates[candidate.MaintenanceActionCandidateID] = candidate
	return nil
}

func (s *MaintenanceStore) GetByID(id string) (maintenance.ActionCandidate, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[id]
	return candidate, found, nil
}
