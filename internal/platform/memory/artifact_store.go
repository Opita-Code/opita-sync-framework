package memory

import (
	"context"
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

func (s *ArtifactStore) Put(ctx context.Context, req storage.PutRequest) (storage.Artifact, error) {
	select {
	case <-ctx.Done():
		return storage.Artifact{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artifacts[req.Artifact.ArtifactRef] = storage.GetResponse{Artifact: req.Artifact, Body: req.Body}
	return req.Artifact, nil
}

func (s *ArtifactStore) Get(ctx context.Context, artifactRef string) (storage.GetResponse, bool, error) {
	select {
	case <-ctx.Done():
		return storage.GetResponse{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	resp, found := s.artifacts[artifactRef]
	return resp, found, nil
}

var _ storage.Service = (*ArtifactStore)(nil)

func (s *ArtifactStore) MustGet(ctx context.Context, artifactRef string) (storage.GetResponse, error) {
	resp, found, _ := s.Get(ctx, artifactRef)
	if !found {
		return storage.GetResponse{}, errors.New("artifact not found")
	}
	return resp, nil
}
