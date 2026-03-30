package tenant

import "time"

type State string

const (
	StateDraft       State = "draft"
	StateConfiguring State = "configuring"
	StateConfigured  State = "configured"
	StateOperable    State = "operable"
	StateBlocked     State = "blocked"
)

type OperabilityCheck struct {
	CheckName string `json:"check_name"`
	Passed    bool   `json:"passed"`
	Message   string `json:"message"`
}

type OperabilityReport struct {
	TenantID string             `json:"tenant_id"`
	Operable bool               `json:"operable"`
	Checks   []OperabilityCheck `json:"checks"`
}

type PolicyBaseline struct {
	ProfileRef          string   `json:"profile_ref"`
	DecisionModel       string   `json:"decision_model"`
	DefaultAction       string   `json:"default_action"`
	SupportedDecisions  []string `json:"supported_decisions"`
	GovernanceMode      string   `json:"governance_mode"`
	ClassificationGuard bool     `json:"classification_guard"`
}

type ApprovalBaseline struct {
	ProfileRef            string   `json:"profile_ref"`
	DefaultMode           string   `json:"default_mode"`
	RequiredStates        []string `json:"required_states"`
	DecisionActions       []string `json:"decision_actions"`
	RequireDecisionActor  bool     `json:"require_decision_actor"`
	FingerprintValidation bool     `json:"fingerprint_validation"`
}

type ClassificationBaseline struct {
	ProfileRef         string   `json:"profile_ref"`
	DefaultLevel       string   `json:"default_level"`
	AllowedLevels      []string `json:"allowed_levels"`
	RestrictedViewMode string   `json:"restricted_view_mode"`
	RedactionRequired  bool     `json:"redaction_required"`
}

type Visibility struct {
	Visible       bool `json:"visible"`
	Usable        bool `json:"usable"`
	Approvable    bool `json:"approvable"`
	Assignable    bool `json:"assignable"`
	Administrable bool `json:"administrable"`
}

type CatalogCapability struct {
	CapabilityID     string     `json:"capability_id"`
	Name             string     `json:"name"`
	ResultType       string     `json:"result_type"`
	Sensitivity      string     `json:"sensitivity"`
	RiskLevel        string     `json:"risk_level"`
	RequiresApproval bool       `json:"requires_approval"`
	Status           string     `json:"status"`
	TenantVisibility Visibility `json:"tenant_visibility"`
}

type ConnectorProjection struct {
	ConnectorRef     string   `json:"connector_ref"`
	Enabled          bool     `json:"enabled"`
	Scope            string   `json:"scope"`
	SupportedActions []string `json:"supported_actions"`
}

type BootstrapRecord struct {
	TenantID                  string                 `json:"tenant_id"`
	TenantName                string                 `json:"tenant_name"`
	AdminSubjectID            string                 `json:"admin_subject_id"`
	InitialCatalogRefs        []string               `json:"initial_catalog_refs"`
	InitialConnectorRefs      []string               `json:"initial_connector_refs"`
	PolicyProfileRef          string                 `json:"policy_profile_ref"`
	ApprovalProfileRef        string                 `json:"approval_profile_ref"`
	ClassificationProfileRef  string                 `json:"classification_profile_ref"`
	PolicyBaseline            PolicyBaseline         `json:"policy_baseline"`
	ApprovalBaseline          ApprovalBaseline       `json:"approval_baseline"`
	ClassificationBaseline    ClassificationBaseline `json:"classification_baseline"`
	CatalogProjection         []CatalogCapability    `json:"catalog_projection"`
	ConnectorProjection       []ConnectorProjection  `json:"connector_projection"`
	ContextSeed               map[string]any         `json:"context_seed,omitempty"`
	State                     State                  `json:"state"`
	BootstrapRecordRef        string                 `json:"tenant_bootstrap_record"`
	PolicyProfileAppliedRef   string                 `json:"tenant_policy_profile_applied"`
	ApprovalProfileAppliedRef string                 `json:"tenant_approval_profile_applied"`
	ClassificationAppliedRef  string                 `json:"tenant_classification_profile_applied"`
	CatalogProjectionRef      string                 `json:"tenant_catalog_projection"`
	ConnectorProjectionRef    string                 `json:"tenant_connector_projection"`
	OperabilityReport         OperabilityReport      `json:"tenant_operability_report"`
	CreatedAt                 time.Time              `json:"created_at"`
	UpdatedAt                 time.Time              `json:"updated_at"`
}

type Store interface {
	Save(record BootstrapRecord) error
	GetByTenantID(tenantID string) (BootstrapRecord, bool, error)
}
