package sdk

import "time"

type RequestMeta struct {
	TenantID       string    `json:"tenant_id"`
	CapabilityID   string    `json:"capability_id"`
	BindingRef     string    `json:"binding_ref"`
	ExecutionID    string    `json:"execution_id,omitempty"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	Attempt        int       `json:"attempt,omitempty"`
	TraceRef       string    `json:"trace_ref"`
	ContractRef    string    `json:"compiled_contract_ref,omitempty"`
	TargetRef      string    `json:"target_ref,omitempty"`
	RequestedScope string    `json:"requested_scope,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

type EvidenceRefs struct {
	TraceRef          string `json:"trace_ref"`
	ProviderCallRef   string `json:"provider_call_ref"`
	InputSnapshotRef  string `json:"input_snapshot_ref,omitempty"`
	OutputSnapshotRef string `json:"output_snapshot_ref,omitempty"`
	ArtifactRef       string `json:"artifact_ref,omitempty"`
	ErrorRef          string `json:"error_ref,omitempty"`
}

type InspectRequest struct {
	Meta RequestMeta `json:"meta"`
}

type InspectResponse struct {
	Available            bool         `json:"available"`
	RelevantMetadata     []string     `json:"relevant_metadata"`
	RestrictionsDetected []string     `json:"restrictions_detected"`
	ClassificationHint   string       `json:"classification_hint"`
	RiskHint             string       `json:"risk_hint"`
	Evidence             EvidenceRefs `json:"evidence"`
}

type DryRunRequest struct {
	Meta RequestMeta `json:"meta"`
}

type DryRunResponse struct {
	ExpectedChanges   []string     `json:"expected_changes"`
	AffectedObjects   []string     `json:"affected_objects"`
	RisksDetected     []string     `json:"risks_detected"`
	NormalizedPreview string       `json:"normalized_preview"`
	Evidence          EvidenceRefs `json:"evidence"`
}

type ExecuteRequest struct {
	Meta RequestMeta `json:"meta"`
}

type ExecuteResponse struct {
	RawResult          string       `json:"raw_result"`
	NormalizedResult   string       `json:"normalized_result"`
	TechnicalState     string       `json:"technical_state"`
	ClassificationHint string       `json:"classification_hint"`
	Retryable          bool         `json:"retryable"`
	Compensable        bool         `json:"compensable"`
	Evidence           EvidenceRefs `json:"evidence"`
}

type CapabilitiesResponse struct {
	SupportedCapabilities []string `json:"supported_capabilities"`
	Limitations           []string `json:"limitations"`
	CompatibleVersions    []string `json:"compatible_versions"`
}

type RiskProfileResponse struct {
	BusinessRisk       string   `json:"business_risk"`
	SecurityRisk       string   `json:"security_risk"`
	AggravatingFactors []string `json:"aggravating_factors"`
	SuggestedApproval  string   `json:"suggested_approval"`
}

type RequiredScopesResponse struct {
	Required []string `json:"required"`
	Optional []string `json:"optional"`
}

type NormalizeRequest struct {
	Meta      RequestMeta `json:"meta"`
	RawResult string      `json:"raw_result"`
}

type NormalizeResponse struct {
	NormalizedResult string       `json:"normalized_result"`
	ReasonCodes      []string     `json:"reason_codes"`
	Severity         string       `json:"severity"`
	Evidence         EvidenceRefs `json:"evidence"`
}

type CompensateRequest struct {
	Meta RequestMeta `json:"meta"`
}

type CompensateResponse struct {
	CompensationState string       `json:"compensation_state"`
	EscalationNeeded  bool         `json:"escalation_needed"`
	Evidence          EvidenceRefs `json:"evidence"`
}

type Provider interface {
	Inspect(req InspectRequest) (InspectResponse, error)
	DryRun(req DryRunRequest) (DryRunResponse, error)
	Execute(req ExecuteRequest) (ExecuteResponse, error)
	GetCapabilities() (CapabilitiesResponse, error)
	GetRiskProfile(capabilityID string, targetScope string) (RiskProfileResponse, error)
	GetRequiredScopes(capabilityID string) (RequiredScopesResponse, error)
	NormalizeResult(req NormalizeRequest) (NormalizeResponse, error)
	Compensate(req CompensateRequest) (CompensateResponse, error)
}
