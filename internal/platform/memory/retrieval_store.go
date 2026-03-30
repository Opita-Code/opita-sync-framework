package memory

import (
	"strings"
	"sync"

	"opita-sync-framework/internal/retrieval"
)

type RetrievalStore struct {
	mu        sync.RWMutex
	documents map[string]retrieval.Document
}

func NewRetrievalStore() *RetrievalStore {
	return &RetrievalStore{documents: map[string]retrieval.Document{}}
}

func (s *RetrievalStore) Index(document retrieval.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[document.DocumentRef] = document
	return nil
}

func (s *RetrievalStore) Search(query retrieval.Query) ([]retrieval.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	matches := make([]retrieval.Match, 0)
	for _, doc := range s.documents {
		if doc.TenantID != query.TenantID {
			continue
		}
		if strings.Contains(strings.ToLower(doc.Text), strings.ToLower(query.Text)) {
			matches = append(matches, retrieval.Match{
				DocumentRef: doc.DocumentRef,
				ArtifactRef: doc.ArtifactRef,
				Score:       1.0,
				Snippet:     doc.Text,
			})
		}
	}
	if query.Limit > 0 && len(matches) > query.Limit {
		matches = matches[:query.Limit]
	}
	return matches, nil
}

var _ retrieval.Service = (*RetrievalStore)(nil)
