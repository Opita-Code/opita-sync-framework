package memory

import (
	"fmt"
	"strings"

	"opita-sync-framework/internal/engine/registry"
)

type RegistryResolver struct{}

func NewRegistryResolver() *RegistryResolver {
	return &RegistryResolver{}
}

func (r *RegistryResolver) Resolve(req registry.ResolutionRequest) (registry.ResolutionResult, error) {
	capabilityID := strings.TrimSpace(req.CapabilityID)
	if capabilityID == "" {
		capabilityID = "capability.plan.default"
	}

	return registry.ResolutionResult{
		CapabilityManifestRef:  fmt.Sprintf("manifest://%s", capabilityID),
		BundleDigest:           fmt.Sprintf("sha256:%s", strings.Repeat("a", 64)),
		BindingID:              fmt.Sprintf("binding-%s-%s", strings.ReplaceAll(capabilityID, ".", "-"), req.Environment),
		ProviderRef:            fmt.Sprintf("provider://%s/default", capabilityID),
		ProviderRuntimeVersion: "provider-runtime-v1",
		Resolved:               true,
	}, nil
}
