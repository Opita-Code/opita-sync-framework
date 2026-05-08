package retrieval

import "context"

type Document struct {
	DocumentRef         string   `json:"document_ref"`
	TenantID            string   `json:"tenant_id"`
	ArtifactRef         string   `json:"artifact_ref"`
	ClassificationLevel string   `json:"classification_level"`
	Title               string   `json:"title"`
	Text                string   `json:"text"`
	Tags                []string `json:"tags"`
}

type Query struct {
	TenantID string `json:"tenant_id"`
	Text     string `json:"text"`
	Limit    int    `json:"limit"`
}

type Match struct {
	DocumentRef string  `json:"document_ref"`
	ArtifactRef string  `json:"artifact_ref"`
	Score       float64 `json:"score"`
	Snippet     string  `json:"snippet"`
}

type Service interface {
	Index(ctx context.Context, document Document) error
	Search(ctx context.Context, query Query) ([]Match, error)
}
