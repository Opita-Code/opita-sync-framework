package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/proposal"
)

type ProposalStore struct {
	mu       sync.RWMutex
	drafts   map[string]proposal.Draft
	patchset map[string]proposal.PatchsetCandidate
}

func NewProposalStore() *ProposalStore {
	return &ProposalStore{
		drafts:   map[string]proposal.Draft{},
		patchset: map[string]proposal.PatchsetCandidate{},
	}
}

func (s *ProposalStore) CreateDraft(draft proposal.Draft) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.drafts[draft.ProposalDraftID]; exists {
		return errors.New("proposal draft already exists")
	}
	s.drafts[draft.ProposalDraftID] = draft
	return nil
}

func (s *ProposalStore) GetDraft(proposalDraftID string) (proposal.Draft, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	draft, found := s.drafts[proposalDraftID]
	return draft, found, nil
}

func (s *ProposalStore) SavePatchset(candidate proposal.PatchsetCandidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.patchset[candidate.PatchsetCandidateID] = candidate
	return nil
}

func (s *ProposalStore) GetPatchset(patchsetCandidateID string) (proposal.PatchsetCandidate, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.patchset[patchsetCandidateID]
	return candidate, found, nil
}
