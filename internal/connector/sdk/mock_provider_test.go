package sdk_test

import (
	"testing"
	"time"

	"opita-sync-framework/internal/connector/sdk"
)

func TestMockProviderExecuteReturnsNormalizedResult(t *testing.T) {
	p := sdk.NewMockProvider()
	resp, err := p.Execute(sdk.ExecuteRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "capability.execution.default", TraceRef: "trace-1", OccurredAt: time.Now().UTC()}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.NormalizedResult == "" {
		t.Fatalf("expected normalized result")
	}
	if resp.Evidence.TraceRef != "trace-1" {
		t.Fatalf("expected trace ref to propagate")
	}
}

func TestMockProviderGetCapabilities(t *testing.T) {
	p := sdk.NewMockProvider()
	resp, err := p.GetCapabilities()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.SupportedCapabilities) == 0 {
		t.Fatalf("expected supported capabilities")
	}
}
