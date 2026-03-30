package sdk

import (
	"fmt"
	"time"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Inspect(req InspectRequest) (InspectResponse, error) {
	return InspectResponse{
		Available:            true,
		RelevantMetadata:     []string{"mock-target-available"},
		RestrictionsDetected: []string{},
		ClassificationHint:   "internal",
		RiskHint:             "low",
		Evidence:             evidence(req.Meta.TraceRef, "inspect"),
	}, nil
}

func (p *MockProvider) DryRun(req DryRunRequest) (DryRunResponse, error) {
	return DryRunResponse{
		ExpectedChanges:   []string{"mock-change"},
		AffectedObjects:   []string{req.Meta.TargetRef},
		RisksDetected:     []string{},
		NormalizedPreview: "mock-preview",
		Evidence:          evidence(req.Meta.TraceRef, "dry_run"),
	}, nil
}

func (p *MockProvider) Execute(req ExecuteRequest) (ExecuteResponse, error) {
	return ExecuteResponse{
		RawResult:          "mock-executed",
		NormalizedResult:   "mock-normalized",
		TechnicalState:     "success",
		ClassificationHint: "internal",
		Retryable:          false,
		Compensable:        true,
		Evidence:           evidence(req.Meta.TraceRef, "execute"),
	}, nil
}

func (p *MockProvider) GetCapabilities() (CapabilitiesResponse, error) {
	return CapabilitiesResponse{
		SupportedCapabilities: []string{"capability.plan.default", "capability.execution.default"},
		Limitations:           []string{},
		CompatibleVersions:    []string{"1.0"},
	}, nil
}

func (p *MockProvider) GetRiskProfile(capabilityID string, targetScope string) (RiskProfileResponse, error) {
	return RiskProfileResponse{
		BusinessRisk:       "low",
		SecurityRisk:       "low",
		AggravatingFactors: []string{},
		SuggestedApproval:  "auto",
	}, nil
}

func (p *MockProvider) GetRequiredScopes(capabilityID string) (RequiredScopesResponse, error) {
	return RequiredScopesResponse{
		Required: []string{"scope:basic"},
		Optional: []string{"scope:extended"},
	}, nil
}

func (p *MockProvider) NormalizeResult(req NormalizeRequest) (NormalizeResponse, error) {
	return NormalizeResponse{
		NormalizedResult: req.RawResult,
		ReasonCodes:      []string{"connector.normalized.mock"},
		Severity:         "info",
		Evidence:         evidence(req.Meta.TraceRef, "normalize"),
	}, nil
}

func (p *MockProvider) Compensate(req CompensateRequest) (CompensateResponse, error) {
	return CompensateResponse{
		CompensationState: "compensated",
		EscalationNeeded:  false,
		Evidence:          evidence(req.Meta.TraceRef, "compensate"),
	}, nil
}

func evidence(traceRef, method string) EvidenceRefs {
	seed := fmt.Sprintf("%s-%s-%d", traceRef, method, time.Now().UTC().UnixNano())
	return EvidenceRefs{
		TraceRef:          traceRef,
		ProviderCallRef:   seed,
		InputSnapshotRef:  seed + ":input",
		OutputSnapshotRef: seed + ":output",
		ArtifactRef:       seed + ":artifact",
	}
}
