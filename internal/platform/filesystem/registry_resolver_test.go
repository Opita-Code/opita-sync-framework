package filesystem

import (
	"path/filepath"
	"testing"

	"opita-sync-framework/internal/engine/registry"
)

func TestRegistryResolverResolveSuccess(t *testing.T) {
	root := filepath.Join("..", "..", "..", "definitions", "capabilities")
	resolver, err := NewRegistryResolver(root)
	if err != nil {
		t.Fatalf("expected resolver, got error: %v", err)
	}
	result, err := resolver.Resolve(registry.ResolutionRequest{
		CapabilityID:          "capability.plan.default",
		ContractSchemaVersion: "1.0",
		SupportedResultType:   "plan",
		Environment:           "dev",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Resolved {
		t.Fatalf("expected resolved result")
	}
	if result.BindingID == "" || result.ProviderRef == "" {
		t.Fatalf("expected binding and provider refs")
	}
	if result.ProviderRef != "provider://tenant.configuration/plan" {
		t.Fatalf("expected tenant configuration plan provider, got %s", result.ProviderRef)
	}
}

func TestRegistryResolverRejectsIncompatibleContractVersion(t *testing.T) {
	root := filepath.Join("..", "..", "..", "definitions", "capabilities")
	resolver, err := NewRegistryResolver(root)
	if err != nil {
		t.Fatalf("expected resolver, got error: %v", err)
	}
	_, err = resolver.Resolve(registry.ResolutionRequest{
		CapabilityID:          "capability.plan.default",
		ContractSchemaVersion: "2.0",
		SupportedResultType:   "plan",
		Environment:           "dev",
	})
	if err == nil {
		t.Fatalf("expected incompatibility error")
	}
}
