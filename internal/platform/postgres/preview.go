package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/preview"
)

type PreviewStore struct {
	store *Store
}

func NewPreviewStore(store *Store) *PreviewStore {
	return &PreviewStore{store: store}
}

func (s *PreviewStore) CreateCandidate(candidate preview.Candidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal preview candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into preview_candidates (preview_candidate_id, tenant_id, session_id, subject_id, contract_id, execution_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, candidate.PreviewCandidateID, candidate.TenantID, candidate.SessionID, candidate.SubjectID, candidate.ContractID, candidate.ExecutionID, raw, candidate.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert preview candidate: %w", err)
	}
	return nil
}

func (s *PreviewStore) GetCandidate(previewCandidateID string) (preview.Candidate, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from preview_candidates where preview_candidate_id = $1`, previewCandidateID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return preview.Candidate{}, false, nil
		}
		return preview.Candidate{}, false, fmt.Errorf("select preview candidate: %w", err)
	}
	var candidate preview.Candidate
	if err := json.Unmarshal(raw, &candidate); err != nil {
		return preview.Candidate{}, false, fmt.Errorf("unmarshal preview candidate: %w", err)
	}
	return candidate, true, nil
}

func (s *PreviewStore) SaveResult(result preview.Result) error {
	raw, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal simulation result: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into simulation_results (simulation_result_id, preview_candidate_id, family, payload, created_at)
		values ($1, $2, $3, $4, $5)
		on conflict (simulation_result_id) do update set payload = excluded.payload
	`, result.SimulationResultID, result.PreviewCandidateID, result.Family, raw, result.CreatedAt)
	if err != nil {
		return fmt.Errorf("upsert simulation result: %w", err)
	}
	return nil
}

func (s *PreviewStore) ListResults(previewCandidateID string) ([]preview.Result, error) {
	rows, err := s.store.DB.QueryContext(contextBackground(), `select payload from simulation_results where preview_candidate_id = $1 order by created_at asc`, previewCandidateID)
	if err != nil {
		return nil, fmt.Errorf("select simulation results: %w", err)
	}
	defer rows.Close()
	out := make([]preview.Result, 0)
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var result preview.Result
		if err := json.Unmarshal(raw, &result); err != nil {
			continue
		}
		out = append(out, result)
	}
	return out, nil
}
