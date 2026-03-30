package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"opita-sync-framework/internal/artifacts/storage"
)

type ArtifactStore struct {
	root string
}

func NewArtifactStore(root string) (*ArtifactStore, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("ensure artifact root: %w", err)
	}
	return &ArtifactStore{root: root}, nil
}

func (s *ArtifactStore) Put(req storage.PutRequest) (storage.Artifact, error) {
	path := s.pathFor(req.Artifact.ArtifactRef)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return storage.Artifact{}, fmt.Errorf("ensure artifact path: %w", err)
	}
	if err := os.WriteFile(path, req.Body, 0o644); err != nil {
		return storage.Artifact{}, fmt.Errorf("write artifact body: %w", err)
	}
	artifact := req.Artifact
	artifact.StorageLocation = path
	artifact.SizeBytes = int64(len(req.Body))
	return artifact, nil
}

func (s *ArtifactStore) Get(artifactRef string) (storage.GetResponse, bool, error) {
	path := s.pathFor(artifactRef)
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.GetResponse{}, false, nil
		}
		return storage.GetResponse{}, false, fmt.Errorf("read artifact body: %w", err)
	}
	return storage.GetResponse{
		Artifact: storage.Artifact{
			ArtifactRef:     artifactRef,
			StorageLocation: path,
			SizeBytes:       int64(len(body)),
		},
		Body: body,
	}, true, nil
}

func (s *ArtifactStore) pathFor(artifactRef string) string {
	normalized := strings.ReplaceAll(artifactRef, "://", "_")
	normalized = strings.ReplaceAll(normalized, "/", "_")
	return filepath.Join(s.root, normalized+".bin")
}

var _ storage.Service = (*ArtifactStore)(nil)
