package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"opita-sync-framework/internal/engine/registry"
)

type RegistryResolver struct {
	capabilities map[string]capabilityDocument
}

type capabilityDocument struct {
	APIVersion string             `yaml:"api_version" json:"api_version"`
	Kind       string             `yaml:"kind" json:"kind"`
	Metadata   capabilityMetadata `yaml:"metadata" json:"metadata"`
	Spec       capabilitySpec     `yaml:"spec" json:"spec"`
}

type capabilityMetadata struct {
	ID      string `yaml:"id" json:"id"`
	Version string `yaml:"version" json:"version"`
}

type capabilitySpec struct {
	SupportedResultTypes  []string      `yaml:"supported_result_types" json:"supported_result_types"`
	SupportedContracts    []string      `yaml:"supported_contract_versions" json:"supported_contract_versions"`
	Packaging             packagingSpec `yaml:"packaging" json:"packaging"`
	Bindings              []bindingSpec `yaml:"bindings" json:"bindings"`
	ProviderCompatibility providerSpec  `yaml:"provider_compatibility" json:"provider_compatibility"`
}

type packagingSpec struct {
	BundleDigest string `yaml:"bundle_digest" json:"bundle_digest"`
}

type bindingSpec struct {
	BindingID                 string   `yaml:"binding_id" json:"binding_id"`
	Environment               string   `yaml:"environment" json:"environment"`
	SupportedResultTypes      []string `yaml:"supported_result_types" json:"supported_result_types"`
	SupportedContractVersions []string `yaml:"supported_contract_versions" json:"supported_contract_versions"`
	ProviderRef               string   `yaml:"provider_ref" json:"provider_ref"`
	ProviderRuntimeVersion    string   `yaml:"provider_runtime_version" json:"provider_runtime_version"`
	Status                    string   `yaml:"status" json:"status"`
}

type providerSpec struct {
	SupportedRuntimeVersions []string `yaml:"supported_runtime_versions" json:"supported_runtime_versions"`
}

func NewRegistryResolver(root string) (*RegistryResolver, error) {
	resolver := &RegistryResolver{capabilities: map[string]capabilityDocument{}}
	files, err := filepath.Glob(filepath.Join(root, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob manifests: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no capability manifests found in %s", root)
	}
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read manifest %s: %w", file, err)
		}
		var doc capabilityDocument
		if err := json.Unmarshal(raw, &doc); err != nil {
			return nil, fmt.Errorf("parse manifest %s: %w", file, err)
		}
		if strings.TrimSpace(doc.Metadata.ID) == "" {
			return nil, fmt.Errorf("manifest %s missing metadata.id", file)
		}
		resolver.capabilities[doc.Metadata.ID] = doc
	}
	return resolver, nil
}

func (r *RegistryResolver) Resolve(req registry.ResolutionRequest) (registry.ResolutionResult, error) {
	doc, ok := r.capabilities[req.CapabilityID]
	if !ok {
		return registry.ResolutionResult{}, fmt.Errorf("capability manifest not found: %s", req.CapabilityID)
	}

	if !contains(doc.Spec.SupportedContracts, req.ContractSchemaVersion) {
		return registry.ResolutionResult{}, fmt.Errorf("contract version %s not supported by capability %s", req.ContractSchemaVersion, req.CapabilityID)
	}

	if !contains(doc.Spec.SupportedResultTypes, req.SupportedResultType) {
		return registry.ResolutionResult{}, fmt.Errorf("result type %s not supported by capability %s", req.SupportedResultType, req.CapabilityID)
	}

	for _, binding := range doc.Spec.Bindings {
		if binding.Environment != req.Environment {
			continue
		}
		if binding.Status != "active" {
			continue
		}
		if !contains(binding.SupportedContractVersions, req.ContractSchemaVersion) {
			continue
		}
		if !contains(binding.SupportedResultTypes, req.SupportedResultType) {
			continue
		}
		return registry.ResolutionResult{
			CapabilityManifestRef:  fmt.Sprintf("manifest://%s@%s", doc.Metadata.ID, doc.Metadata.Version),
			BundleDigest:           doc.Spec.Packaging.BundleDigest,
			BindingID:              binding.BindingID,
			ProviderRef:            binding.ProviderRef,
			ProviderRuntimeVersion: binding.ProviderRuntimeVersion,
			Resolved:               true,
		}, nil
	}

	return registry.ResolutionResult{}, fmt.Errorf("no active binding found for capability=%s env=%s contract=%s result=%s", req.CapabilityID, req.Environment, req.ContractSchemaVersion, req.SupportedResultType)
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
