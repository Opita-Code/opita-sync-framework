package retrieval

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
	Index(document Document) error
	Search(query Query) ([]Match, error)
}
