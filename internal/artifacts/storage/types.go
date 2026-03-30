package storage

import "time"

type Artifact struct {
	ArtifactRef         string    `json:"artifact_ref"`
	TenantID            string    `json:"tenant_id"`
	Kind                string    `json:"kind"`
	ClassificationLevel string    `json:"classification_level"`
	ContentType         string    `json:"content_type"`
	StorageLocation     string    `json:"storage_location"`
	SizeBytes           int64     `json:"size_bytes"`
	CreatedAt           time.Time `json:"created_at"`
}

type PutRequest struct {
	Artifact Artifact `json:"artifact"`
	Body     []byte   `json:"-"`
}

type GetResponse struct {
	Artifact Artifact `json:"artifact"`
	Body     []byte   `json:"-"`
}

type Service interface {
	Put(req PutRequest) (Artifact, error)
	Get(artifactRef string) (GetResponse, bool, error)
}
