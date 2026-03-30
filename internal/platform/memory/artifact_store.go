package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/artifacts/storage"
)

type ArtifactStore struct {
	mu        sync.RWMutex
	artifacts map[string]storage.GetResponse
}

func NewArtifactStore() *ArtifactStore {
	return &ArtifactStore{artifacts: map[string]storage.GetResponse{}}
}

func (s *ArtifactStore) Put(req storage.PutRequest) (storage.Artifact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artifacts[req.Artifact.ArtifactRef] = storage.GetResponse{Artifact: req.Artifact, Body: req.Body}
	return req.Artifact, nil
}

func (s *ArtifactStore) Get(artifactRef string) (storage.GetResponse, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resp, found := s.artifacts[artifactRef]
	return resp, found, nil
}

var _ storage.Service = (*ArtifactStore)(nil)

func (s *ArtifactStore) MustGet(artifactRef string) (storage.GetResponse, error) {
	resp, found, _ := s.Get(artifactRef)
	if !found {
		return storage.GetResponse{}, errors.New("artifact not found")
	}
	return resp, nil
}
