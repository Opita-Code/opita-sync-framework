package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/proposal"
)

var _ proposal.Service = (*ProposalStore)(nil)

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

func (s *ProposalStore) CreateDraft(ctx context.Context, draft proposal.Draft) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.drafts[draft.ProposalDraftID]; exists {
		return errors.New("proposal draft already exists")
	}
	s.drafts[draft.ProposalDraftID] = draft
	return nil
}

func (s *ProposalStore) GetDraft(ctx context.Context, proposalDraftID string) (proposal.Draft, bool, error) {
	select {
	case <-ctx.Done():
		return proposal.Draft{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	draft, found := s.drafts[proposalDraftID]
	return draft, found, nil
}

func (s *ProposalStore) SavePatchset(ctx context.Context, candidate proposal.PatchsetCandidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.patchset[candidate.PatchsetCandidateID] = candidate
	return nil
}

func (s *ProposalStore) GetPatchset(ctx context.Context, patchsetCandidateID string) (proposal.PatchsetCandidate, bool, error) {
	select {
	case <-ctx.Done():
		return proposal.PatchsetCandidate{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.patchset[patchsetCandidateID]
	return candidate, found, nil
}
