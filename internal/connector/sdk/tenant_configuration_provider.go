package sdk

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type TenantConfigurationProvider struct {
	mu              sync.Mutex
	executeByKey    map[string]ExecuteResponse
	compensateByKey map[string]CompensateResponse
}

func NewTenantConfigurationProvider() *TenantConfigurationProvider {
	return &TenantConfigurationProvider{
		executeByKey:    map[string]ExecuteResponse{},
		compensateByKey: map[string]CompensateResponse{},
	}
}

func (p *TenantConfigurationProvider) Inspect(req InspectRequest) (InspectResponse, error) {
	classification, risk, restrictions := classifyTenantConfigOperation(req.Meta.CapabilityID, req.Meta.TargetRef)
	return InspectResponse{
		Available:            true,
		RelevantMetadata:     []string{fmt.Sprintf("target:%s", req.Meta.TargetRef), fmt.Sprintf("binding:%s", req.Meta.BindingRef)},
		RestrictionsDetected: restrictions,
		ClassificationHint:   classification,
		RiskHint:             risk,
		Evidence:             tenantEvidence(req.Meta, "inspect", "inspect"),
	}, nil
}

func (p *TenantConfigurationProvider) DryRun(req DryRunRequest) (DryRunResponse, error) {
	classification, risk, restrictions := classifyTenantConfigOperation(req.Meta.CapabilityID, req.Meta.TargetRef)
	return DryRunResponse{
		ExpectedChanges:   []string{fmt.Sprintf("tenant config change for %s", req.Meta.TargetRef), fmt.Sprintf("capability %s prepared", req.Meta.CapabilityID)},
		AffectedObjects:   []string{req.Meta.TargetRef},
		RisksDetected:     append([]string{fmt.Sprintf("risk:%s", risk), fmt.Sprintf("classification:%s", classification)}, restrictions...),
		NormalizedPreview: fmt.Sprintf("preview tenant change for %s under scope %s with classification %s", req.Meta.TargetRef, req.Meta.RequestedScope, classification),
		Evidence:          tenantEvidence(req.Meta, "dry-run", "dry_run"),
	}, nil
}

func (p *TenantConfigurationProvider) Execute(req ExecuteRequest) (ExecuteResponse, error) {
	if strings.TrimSpace(req.Meta.IdempotencyKey) == "" {
		return ExecuteResponse{}, errors.New("idempotency_key is required for tenant configuration execution")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if existing, ok := p.executeByKey[req.Meta.IdempotencyKey]; ok {
		return existing, nil
	}
	classification, risk, restrictions := classifyTenantConfigOperation(req.Meta.CapabilityID, req.Meta.TargetRef)
	technicalState := "success"
	normalized := fmt.Sprintf("tenant configuration applied for %s using %s", req.Meta.TargetRef, req.Meta.CapabilityID)
	if len(restrictions) > 0 {
		normalized = fmt.Sprintf("tenant configuration applied with restrictions for %s using %s", req.Meta.TargetRef, req.Meta.CapabilityID)
	}
	resp := ExecuteResponse{
		RawResult:          fmt.Sprintf("applied:%s:%s", req.Meta.TargetRef, req.Meta.CapabilityID),
		NormalizedResult:   normalized,
		TechnicalState:     technicalState,
		ClassificationHint: classification,
		Retryable:          false,
		Compensable:        risk != "low",
		Evidence:           tenantEvidence(req.Meta, req.Meta.IdempotencyKey, "execute"),
	}
	p.executeByKey[req.Meta.IdempotencyKey] = resp
	return resp, nil
}

func (p *TenantConfigurationProvider) GetCapabilities() (CapabilitiesResponse, error) {
	return CapabilitiesResponse{
		SupportedCapabilities: []string{
			"tenant.intake.capture_intent",
			"tenant.proposal.create_change_draft",
			"tenant.preview.run_governance_preview",
			"tenant.execution.compile_governed_intent",
			"tenant.execution.inspect_run",
			"tenant.approval.release_blocked_execution",
			"tenant.recovery.resume_after_approval",
			"tenant.recovery.request_manual_compensation",
			"capability.plan.default",
			"capability.execution.default",
		},
		Limitations:        []string{"tenant configuration domain only", "does not replace runtime or policy enforcement"},
		CompatibleVersions: []string{"1.0"},
	}, nil
}

func (p *TenantConfigurationProvider) GetRiskProfile(capabilityID string, targetScope string) (RiskProfileResponse, error) {
	_, risk, restrictions := classifyTenantConfigOperation(capabilityID, targetScope)
	securityRisk := "low"
	suggestedApproval := "auto"
	if risk == "medium" {
		securityRisk = "medium"
		suggestedApproval = "pre_execution"
	}
	if risk == "high" {
		securityRisk = "high"
		suggestedApproval = "pre_execution"
	}
	return RiskProfileResponse{
		BusinessRisk:       risk,
		SecurityRisk:       securityRisk,
		AggravatingFactors: restrictions,
		SuggestedApproval:  suggestedApproval,
	}, nil
}

func (p *TenantConfigurationProvider) GetRequiredScopes(capabilityID string) (RequiredScopesResponse, error) {
	required := []string{"scope:tenant-config.read"}
	if strings.Contains(capabilityID, "execution") || strings.Contains(capabilityID, "approval") || strings.Contains(capabilityID, "classification") {
		required = append(required, "scope:tenant-config.write", "scope:tenant-config.approve")
	} else {
		required = append(required, "scope:tenant-config.write")
	}
	return RequiredScopesResponse{Required: required, Optional: []string{"scope:tenant-config.inspect"}}, nil
}

func (p *TenantConfigurationProvider) NormalizeResult(req NormalizeRequest) (NormalizeResponse, error) {
	return NormalizeResponse{
		NormalizedResult: fmt.Sprintf("normalized tenant-config result: %s", req.RawResult),
		ReasonCodes:      []string{"connector.normalized.tenant_configuration"},
		Severity:         "info",
		Evidence:         tenantEvidence(req.Meta, req.Meta.IdempotencyKey, "normalize"),
	}, nil
}

func (p *TenantConfigurationProvider) Compensate(req CompensateRequest) (CompensateResponse, error) {
	key := strings.TrimSpace(req.Meta.IdempotencyKey)
	if key == "" {
		key = fmt.Sprintf("compensate:%s:%s", req.Meta.ExecutionID, req.Meta.TargetRef)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if existing, ok := p.compensateByKey[key]; ok {
		return existing, nil
	}
	resp := CompensateResponse{
		CompensationState: "compensation_pending",
		EscalationNeeded:  strings.Contains(req.Meta.CapabilityID, "classification") || strings.Contains(req.Meta.CapabilityID, "approval"),
		Evidence:          tenantEvidence(req.Meta, key, "compensate"),
	}
	p.compensateByKey[key] = resp
	return resp, nil
}

func classifyTenantConfigOperation(capabilityID, targetRef string) (classification string, risk string, restrictions []string) {
	classification = "internal"
	risk = "low"
	restrictions = []string{}
	combined := capabilityID + ":" + targetRef
	if strings.Contains(combined, "classification") {
		classification = "restricted"
		risk = "high"
		restrictions = append(restrictions, "classification-change-sensitive")
	}
	if strings.Contains(combined, "approval") || strings.Contains(combined, "execution") || strings.Contains(combined, "connector") {
		if risk != "high" {
			risk = "medium"
		}
		restrictions = append(restrictions, "approval-or-execution-path-change")
	}
	return classification, risk, restrictions
}

func tenantEvidence(meta RequestMeta, dedupeSeed string, operation string) EvidenceRefs {
	seed := dedupeSeed
	if strings.TrimSpace(seed) == "" {
		seed = fmt.Sprintf("%s-%s-%d", meta.TraceRef, operation, time.Now().UTC().UnixNano())
	}
	providerCallRef := fmt.Sprintf("provider://tenant.configuration/%s/%s", operation, seed)
	return EvidenceRefs{
		TraceRef:          meta.TraceRef,
		ProviderCallRef:   providerCallRef,
		InputSnapshotRef:  providerCallRef + ":input",
		OutputSnapshotRef: providerCallRef + ":output",
		ArtifactRef:       providerCallRef + ":artifact",
	}
}
