package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/preview"
)

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

func (s *PreviewStore) CreateCandidate(candidate preview.Candidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.candidates[candidate.PreviewCandidateID]; exists {
		return errors.New("preview candidate already exists")
	}
	s.candidates[candidate.PreviewCandidateID] = candidate
	return nil
}

func (s *PreviewStore) GetCandidate(previewCandidateID string) (preview.Candidate, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[previewCandidateID]
	return candidate, found, nil
}

func (s *PreviewStore) SaveResult(result preview.Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.PreviewCandidateID] = append(s.results[result.PreviewCandidateID], result)
	return nil
}

func (s *PreviewStore) ListResults(previewCandidateID string) ([]preview.Result, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := s.results[previewCandidateID]
	out := make([]preview.Result, len(results))
	copy(out, results)
	return out, nil
}
