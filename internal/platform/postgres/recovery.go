package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/inspection"
)

type RecoveryStore struct {
	store *Store
}

func NewRecoveryStore(store *Store) *RecoveryStore {
	return &RecoveryStore{store: store}
}

func (s *RecoveryStore) Create(candidate inspection.RecoveryActionCandidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal recovery candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into recovery_action_candidates (recovery_action_candidate_id, execution_id, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
	`, candidate.RecoveryActionCandidateID, candidate.ExecutionID, raw, candidate.CreatedAt, candidate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert recovery candidate: %w", err)
	}
	return nil
}

func (s *RecoveryStore) GetByID(id string) (inspection.RecoveryActionCandidate, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from recovery_action_candidates where recovery_action_candidate_id = $1`, id).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return inspection.RecoveryActionCandidate{}, false, nil
		}
		return inspection.RecoveryActionCandidate{}, false, fmt.Errorf("select recovery candidate: %w", err)
	}
	var candidate inspection.RecoveryActionCandidate
	if err := json.Unmarshal(raw, &candidate); err != nil {
		return inspection.RecoveryActionCandidate{}, false, fmt.Errorf("unmarshal recovery candidate: %w", err)
	}
	return candidate, true, nil
}

func (s *RecoveryStore) Update(candidate inspection.RecoveryActionCandidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal recovery candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		update recovery_action_candidates set payload = $2, updated_at = $3 where recovery_action_candidate_id = $1
	`, candidate.RecoveryActionCandidateID, raw, candidate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update recovery candidate: %w", err)
	}
	return nil
}
