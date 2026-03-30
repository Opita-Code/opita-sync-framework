package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/maintenance"
)

type MaintenanceStore struct {
	store *Store
}

func NewMaintenanceStore(store *Store) *MaintenanceStore {
	return &MaintenanceStore{store: store}
}

func (s *MaintenanceStore) Create(candidate maintenance.ActionCandidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal maintenance candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into maintenance_action_candidates (maintenance_action_candidate_id, tenant_id, payload, created_at)
		values ($1, $2, $3, $4)
	`, candidate.MaintenanceActionCandidateID, candidate.TenantID, raw, candidate.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert maintenance candidate: %w", err)
	}
	return nil
}

func (s *MaintenanceStore) GetByID(id string) (maintenance.ActionCandidate, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from maintenance_action_candidates where maintenance_action_candidate_id = $1`, id).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return maintenance.ActionCandidate{}, false, nil
		}
		return maintenance.ActionCandidate{}, false, fmt.Errorf("select maintenance candidate: %w", err)
	}
	var candidate maintenance.ActionCandidate
	if err := json.Unmarshal(raw, &candidate); err != nil {
		return maintenance.ActionCandidate{}, false, fmt.Errorf("unmarshal maintenance candidate: %w", err)
	}
	return candidate, true, nil
}
