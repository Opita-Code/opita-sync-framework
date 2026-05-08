package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/preview"
)

var _ preview.Service = (*PreviewStore)(nil)

type PreviewStore struct {
	mu         sync.RWMutex
	candidates map[string]preview.Candidate
	results    map[string][]preview.Result
}

func NewPreviewStore() *PreviewStore {
	return &PreviewStore{
		candidates: map[string]preview.Candidate{},
		results:    map[string][]preview.Result{},
	}
}

func (s *PreviewStore) CreateCandidate(ctx context.Context, candidate preview.Candidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.PreviewCandidateID]; exists {
		return errors.New("preview candidate already exists")
	}
	s.candidates[candidate.PreviewCandidateID] = candidate
	return nil
}

func (s *PreviewStore) GetCandidate(ctx context.Context, previewCandidateID string) (preview.Candidate, bool, error) {
	select {
	case <-ctx.Done():
		return preview.Candidate{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[previewCandidateID]
	return candidate, found, nil
}

func (s *PreviewStore) SaveResult(ctx context.Context, result preview.Result) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.PreviewCandidateID] = append(s.results[result.PreviewCandidateID], result)
	return nil
}

func (s *PreviewStore) ListResults(ctx context.Context, previewCandidateID string) ([]preview.Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := s.results[previewCandidateID]
	out := make([]preview.Result, len(results))
	copy(out, results)
	return out, nil
}
