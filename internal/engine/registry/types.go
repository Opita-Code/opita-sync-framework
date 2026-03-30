package registry

type ResolutionRequest struct {
	CapabilityID          string
	ContractSchemaVersion string
	SupportedResultType   string
	Environment           string
}

type ResolutionResult struct {
	CapabilityManifestRef  string
	BundleDigest           string
	BindingID              string
	ProviderRef            string
	ProviderRuntimeVersion string
	Resolved               bool
}

type Resolver interface {
	Resolve(req ResolutionRequest) (ResolutionResult, error)
}
