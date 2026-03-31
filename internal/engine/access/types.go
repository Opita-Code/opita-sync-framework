package access

import "time"

type State string

const (
	StateDraft   State = "draft"
	StateActive  State = "active"
	StateBlocked State = "blocked"
	StateRevoked State = "revoked"
)

type CapabilityGrant struct {
	GrantID          string    `json:"grant_id"`
	TenantID         string    `json:"tenant_id"`
	PrincipalRef     string    `json:"principal_ref"`
	PrincipalType    string    `json:"principal_type"`
	CapabilityID     string    `json:"capability_id"`
	ScopeRef         string    `json:"scope_ref,omitempty"`
	AllowedActions   []string  `json:"allowed_actions"`
	DeniedActions    []string  `json:"denied_actions,omitempty"`
	RequiresApproval bool      `json:"requires_approval"`
	Justification    string    `json:"justification,omitempty"`
	TraceRef         string    `json:"trace_ref"`
	State            State     `json:"state"`
	ValidFrom        time.Time `json:"valid_from"`
	ValidUntil       time.Time `json:"valid_until,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type DelegationGrant struct {
	GrantID          string    `json:"grant_id"`
	TenantID         string    `json:"tenant_id"`
	SourcePrincipal  string    `json:"source_principal"`
	TargetPrincipal  string    `json:"target_principal"`
	AuthoritySource  string    `json:"authority_source"`
	ScopeType        string    `json:"scope_type"`
	ScopeRef         string    `json:"scope_ref"`
	AllowedActions   []string  `json:"allowed_actions"`
	DeniedActions    []string  `json:"denied_actions,omitempty"`
	RequiresApproval bool      `json:"requires_approval"`
	CanRedelegate    bool      `json:"can_redelegate"`
	MaxDepth         int       `json:"max_depth"`
	Justification    string    `json:"justification,omitempty"`
	TraceRef         string    `json:"trace_ref"`
	State            State     `json:"state"`
	ValidFrom        time.Time `json:"valid_from"`
	ValidUntil       time.Time `json:"valid_until,omitempty"`
	RevokedAt        time.Time `json:"revoked_at,omitempty"`
	RevokedBy        string    `json:"revoked_by,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Store interface {
	SaveGrant(grant CapabilityGrant) error
	ListGrantsByTenant(tenantID string) ([]CapabilityGrant, error)
	SaveDelegation(grant DelegationGrant) error
	ListDelegationsByTenant(tenantID string) ([]DelegationGrant, error)
}
