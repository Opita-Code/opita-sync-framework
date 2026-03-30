package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/runtime"
)

type RuntimeService struct {
	store *Store
}

func NewRuntimeService(store *Store) *RuntimeService {
	return &RuntimeService{store: store}
}

func (s *RuntimeService) CreateExecution(record runtime.ExecutionRecord) error {
	raw, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal execution record: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into execution_records (execution_id, contract_id, tenant_id, trace_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, record.ExecutionID, record.ContractID, record.TenantID, record.TraceID, record.State, raw, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert execution record: %w", err)
	}
	return nil
}

func (s *RuntimeService) GetExecution(executionID string) (runtime.ExecutionRecord, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from execution_records where execution_id = $1`, executionID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return runtime.ExecutionRecord{}, false, nil
		}
		return runtime.ExecutionRecord{}, false, fmt.Errorf("select execution record: %w", err)
	}
	var record runtime.ExecutionRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return runtime.ExecutionRecord{}, false, fmt.Errorf("unmarshal execution record: %w", err)
	}
	return record, true, nil
}

func (s *RuntimeService) UpdateExecutionState(executionID string, state runtime.ExecutionState) (runtime.ExecutionRecord, error) {
	record, found, err := s.GetExecution(executionID)
	if err != nil {
		return runtime.ExecutionRecord{}, err
	}
	if !found {
		return runtime.ExecutionRecord{}, fmt.Errorf("execution not found")
	}
	record.State = state
	raw, err := json.Marshal(record)
	if err != nil {
		return runtime.ExecutionRecord{}, fmt.Errorf("marshal execution record: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		update execution_records set state = $2, payload = $3, updated_at = now() where execution_id = $1
	`, executionID, state, raw)
	if err != nil {
		return runtime.ExecutionRecord{}, fmt.Errorf("update execution record: %w", err)
	}
	return record, nil
}
