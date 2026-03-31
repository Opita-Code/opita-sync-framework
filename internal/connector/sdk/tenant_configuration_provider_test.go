package sdk_test

import (
	"testing"
	"time"

	"opita-sync-framework/internal/connector/sdk"
)

func TestTenantConfigurationProviderExecuteRequiresIdempotencyKey(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	_, err := p.Execute(sdk.ExecuteRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "tenant.execution.compile_governed_intent", TraceRef: "trace-1", TargetRef: "tenant-1/catalog/capability-x", OccurredAt: time.Now().UTC()}})
	if err == nil {
		t.Fatalf("expected idempotency key error")
	}
}

func TestTenantConfigurationProviderExecuteIsIdempotentAndEmitsEvidence(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	req := sdk.ExecuteRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "tenant.execution.compile_governed_intent", BindingRef: "binding-capability-execution-default-dev", IdempotencyKey: "idem-1", TraceRef: "trace-1", TargetRef: "tenant-1/catalog/capability-x", RequestedScope: "tenant_config_change", OccurredAt: time.Now().UTC()}}
	first, err := p.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	second, err := p.Execute(req)
	if err != nil {
		t.Fatalf("expected no error on repeated execute, got %v", err)
	}
	if first.Evidence.ProviderCallRef == "" || first.Evidence.ArtifactRef == "" {
		t.Fatalf("expected evidence refs, got %+v", first.Evidence)
	}
	if first.Evidence.ProviderCallRef != second.Evidence.ProviderCallRef {
		t.Fatalf("expected idempotent provider call ref, got %s and %s", first.Evidence.ProviderCallRef, second.Evidence.ProviderCallRef)
	}
}

func TestTenantConfigurationProviderCapabilitiesAndRiskProfile(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	capabilities, err := p.GetCapabilities()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(capabilities.SupportedCapabilities) == 0 {
		t.Fatalf("expected supported capabilities")
	}
	risk, err := p.GetRiskProfile("tenant.recovery.request_manual_compensation", "tenant-1/classification")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if risk.BusinessRisk == "" || risk.SuggestedApproval == "" {
		t.Fatalf("expected risk profile, got %+v", risk)
	}
}

func TestTenantConfigurationProviderCompensateReturnsEvidence(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	resp, err := p.Compensate(sdk.CompensateRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "tenant.recovery.request_manual_compensation", IdempotencyKey: "comp-1", TraceRef: "trace-1", TargetRef: "tenant-1/baselines/classification", OccurredAt: time.Now().UTC()}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.CompensationState == "" || resp.Evidence.ProviderCallRef == "" {
		t.Fatalf("expected compensation evidence, got %+v", resp)
	}
}

func TestTenantConfigurationProviderRestrictedConnectorRaisesRisk(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	risk, err := p.GetRiskProfile("tenant.execution.compile_governed_intent", "tenant-1/connectors/connector.execution.restricted")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if risk.BusinessRisk != "high" || risk.SecurityRisk != "high" {
		t.Fatalf("expected high risk for restricted connector, got %+v", risk)
	}
	resp, err := p.Execute(sdk.ExecuteRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "tenant.execution.compile_governed_intent", BindingRef: "binding-capability-execution-restricted-dev", IdempotencyKey: "idem-restricted", TraceRef: "trace-1", TargetRef: "tenant-1/connectors/connector.execution.restricted", RequestedScope: "tenant_config_change", OccurredAt: time.Now().UTC()}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.TechnicalState != "success_with_restrictions" || !resp.Retryable {
		t.Fatalf("expected restricted connector execution semantics, got %+v", resp)
	}
}

func TestTenantConfigurationProviderEnterpriseConnectorDemandsVerification(t *testing.T) {
	p := sdk.NewTenantConfigurationProvider()
	risk, err := p.GetRiskProfile("tenant.execution.compile_governed_intent", "tenant-1/connectors/connector.execution.enterprise")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if risk.BusinessRisk != "high" || len(risk.AggravatingFactors) == 0 {
		t.Fatalf("expected high risk with aggravating factors, got %+v", risk)
	}
	scopes, err := p.GetRequiredScopes("tenant.execution.compile_governed_intent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	foundExternal := false
	for _, scope := range scopes.Required {
		if scope == "scope:tenant-config.external-system" {
			foundExternal = true
		}
	}
	if !foundExternal {
		t.Fatalf("expected external system scope in %+v", scopes)
	}
	resp, err := p.Execute(sdk.ExecuteRequest{Meta: sdk.RequestMeta{TenantID: "tenant-1", CapabilityID: "tenant.execution.compile_governed_intent", BindingRef: "binding-capability-execution-enterprise-dev", IdempotencyKey: "idem-enterprise", TraceRef: "trace-1", TargetRef: "tenant-1/connectors/connector.execution.enterprise", RequestedScope: "tenant_config_change", OccurredAt: time.Now().UTC()}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.TechnicalState != "pending_verification" || !resp.Retryable || !resp.Compensable {
		t.Fatalf("expected enterprise verification semantics, got %+v", resp)
	}
}
