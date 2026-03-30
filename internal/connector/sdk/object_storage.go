package sdk

import "opita-sync-framework/internal/artifacts/storage"

type ArtifactAwareProvider interface {
	Provider
	SetArtifactStore(store storage.Service)
}
