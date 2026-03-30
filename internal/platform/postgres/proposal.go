package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"opita-sync-framework/internal/engine/proposal"
)

type ProposalStore struct {
	store *Store
}

func NewProposalStore(store *Store) *ProposalStore {
	return &ProposalStore{store: store}
}

func (s *ProposalStore) CreateDraft(draft proposal.Draft) error {
	raw, err := json.Marshal(draft)
	if err != nil {
		return fmt.Errorf("marshal proposal draft: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into proposal_drafts (proposal_draft_id, tenant_id, session_id, subject_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`, draft.ProposalDraftID, draft.TenantID, draft.SessionID, draft.SubjectID, raw, draft.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert proposal draft: %w", err)
	}
	return nil
}

func (s *ProposalStore) GetDraft(proposalDraftID string) (proposal.Draft, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from proposal_drafts where proposal_draft_id = $1`, proposalDraftID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return proposal.Draft{}, false, nil
		}
		return proposal.Draft{}, false, fmt.Errorf("select proposal draft: %w", err)
	}
	var draft proposal.Draft
	if err := json.Unmarshal(raw, &draft); err != nil {
		return proposal.Draft{}, false, fmt.Errorf("unmarshal proposal draft: %w", err)
	}
	return draft, true, nil
}

func (s *ProposalStore) SavePatchset(candidate proposal.PatchsetCandidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal patchset candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into patchset_candidates (patchset_candidate_id, proposal_draft_id, payload, created_at)
		values ($1, $2, $3, $4)
		on conflict (patchset_candidate_id) do update set payload = excluded.payload
	`, candidate.PatchsetCandidateID, candidate.ProposalDraftID, raw, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("upsert patchset candidate: %w", err)
	}
	return nil
}

func (s *ProposalStore) GetPatchset(patchsetCandidateID string) (proposal.PatchsetCandidate, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from patchset_candidates where patchset_candidate_id = $1`, patchsetCandidateID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return proposal.PatchsetCandidate{}, false, nil
		}
		return proposal.PatchsetCandidate{}, false, fmt.Errorf("select patchset candidate: %w", err)
	}
	var candidate proposal.PatchsetCandidate
	if err := json.Unmarshal(raw, &candidate); err != nil {
		return proposal.PatchsetCandidate{}, false, fmt.Errorf("unmarshal patchset candidate: %w", err)
	}
	return candidate, true, nil
}
