package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"opita-sync-framework/internal/engine/foundation"
)

type FoundationRunStore struct {
	store *Store
}

func NewFoundationRunStore(store *Store) *FoundationRunStore {
	return &FoundationRunStore{store: store}
}

func (s *FoundationRunStore) Save(result foundation.FoundationRunResult) error {
	raw, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal foundation run: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into foundation_runs (execution_id, contract_id, trace_id, payload, created_at)
		values ($1, $2, $3, $4, $5)
		on conflict (execution_id) do update set
		  payload = excluded.payload
	`, result.Execution.ExecutionID, result.Contract.ContractID, result.Execution.TraceID, raw, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("upsert foundation run: %w", err)
	}
	return nil
}

func (s *FoundationRunStore) GetByExecutionID(executionID string) (foundation.FoundationRunResult, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from foundation_runs where execution_id = $1`, executionID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return foundation.FoundationRunResult{}, false, nil
		}
		return foundation.FoundationRunResult{}, false, fmt.Errorf("select foundation run: %w", err)
	}
	var result foundation.FoundationRunResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return foundation.FoundationRunResult{}, false, fmt.Errorf("unmarshal foundation run: %w", err)
	}
	return result, true, nil
}
