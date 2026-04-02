package accessservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/access"
	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
)

type Store interface {
	SaveGrant(grant access.CapabilityGrant) error
	GetGrantByID(grantID string) (access.CapabilityGrant, bool, error)
	ListGrantsByTenant(tenantID string) ([]access.CapabilityGrant, error)
	SaveDelegation(grant access.DelegationGrant) error
	GetDelegationByID(grantID string) (access.DelegationGrant, bool, error)
	ListDelegationsByTenant(tenantID string) ([]access.DelegationGrant, error)
}

type EventWriter interface {
	Append(record events.Record) error
}

type ApprovalService interface {
	Create(request approvals.Request) error
	GetByID(approvalRequestID string) (approvals.Request, bool, error)
	Decide(approvalRequestID string, decision approvals.Decision) (approvals.Request, error)
}

type Handler struct {
	Store     Store
	Events    EventWriter
	Approvals ApprovalService
}

type accessAdminWorkspace struct {
	TenantID    string               `json:"tenant_id"`
	Summary     accessSummaryCard    `json:"summary"`
	Grants      accessGrantCard      `json:"grants"`
	Delegations accessDelegationCard `json:"delegations"`
	Governance  accessGovernanceCard `json:"governance"`
	Impact      accessImpactCard     `json:"impact"`
	Boundary    string               `json:"boundary"`
}

type accessSummaryCard struct {
	TotalGrants      int      `json:"total_grants"`
	TotalDelegations int      `json:"total_delegations"`
	BlockedItems     int      `json:"blocked_items"`
	RevokedItems     int      `json:"revoked_items"`
	ExpiredItems     int      `json:"expired_items"`
	RecommendedNext  []string `json:"recommended_next_actions"`
	Summary          string   `json:"summary"`
}

type accessGrantCard struct {
	ActiveGrants          int      `json:"active_grants"`
	BlockedGrants         int      `json:"blocked_grants"`
	ApprovalRequired      int      `json:"approval_required_grants"`
	BlockedGrantIDs       []string `json:"blocked_grant_ids"`
	RevokedGrantIDs       []string `json:"revoked_grant_ids"`
	ExpiredGrantIDs       []string `json:"expired_grant_ids"`
	SensitiveCapabilities []string `json:"sensitive_capabilities"`
	Summary               string   `json:"summary"`
}

type accessDelegationCard struct {
	ActiveDelegations      int      `json:"active_delegations"`
	BlockedDelegations     int      `json:"blocked_delegations"`
	RedelegableDelegations int      `json:"redelegable_delegations"`
	BlockedDelegationIDs   []string `json:"blocked_delegation_ids"`
	RevokedDelegationIDs   []string `json:"revoked_delegation_ids"`
	ExpiredDelegationIDs   []string `json:"expired_delegation_ids"`
	SensitiveDelegations   []string `json:"sensitive_delegations"`
	Summary                string   `json:"summary"`
}

type accessGovernanceCard struct {
	ApprovalRequestCount int      `json:"approval_request_count"`
	RevocationCount      int      `json:"revocation_count"`
	Guardrails           []string `json:"guardrails"`
	Summary              string   `json:"summary"`
}

type accessImpactCard struct {
	ApprovalSensitiveAreas   []string `json:"approval_sensitive_areas"`
	RevocationSensitiveAreas []string `json:"revocation_sensitive_areas"`
	PromotionAdvice          string   `json:"promotion_advice"`
}

type createGrantRequest struct {
	TenantID       string   `json:"tenant_id"`
	PrincipalRef   string   `json:"principal_ref"`
	PrincipalType  string   `json:"principal_type"`
	CapabilityID   string   `json:"capability_id"`
	ScopeRef       string   `json:"scope_ref,omitempty"`
	AllowedActions []string `json:"allowed_actions"`
	DeniedActions  []string `json:"denied_actions,omitempty"`
	Justification  string   `json:"justification,omitempty"`
	TraceRef       string   `json:"trace_ref"`
	ValidUntil     string   `json:"valid_until,omitempty"`
}

type createDelegationRequest struct {
	TenantID        string   `json:"tenant_id"`
	SourcePrincipal string   `json:"source_principal"`
	TargetPrincipal string   `json:"target_principal"`
	AuthoritySource string   `json:"authority_source"`
	ScopeType       string   `json:"scope_type"`
	ScopeRef        string   `json:"scope_ref"`
	AllowedActions  []string `json:"allowed_actions"`
	DeniedActions   []string `json:"denied_actions,omitempty"`
	CanRedelegate   bool     `json:"can_redelegate"`
	MaxDepth        int      `json:"max_depth"`
	Justification   string   `json:"justification,omitempty"`
	TraceRef        string   `json:"trace_ref"`
	ValidUntil      string   `json:"valid_until,omitempty"`
}

type approvalActionRequest struct {
	DecidedBySubjectID  string   `json:"decided_by_subject_id"`
	DecisionComment     string   `json:"decision_comment,omitempty"`
	DecisionReasonCodes []string `json:"decision_reason_codes,omitempty"`
}

type revokeRequest struct {
	RevokedBy     string `json:"revoked_by"`
	Justification string `json:"justification,omitempty"`
}

func NewHandler(store Store, eventWriter EventWriter, approvalService ApprovalService) *Handler {
	return &Handler{Store: store, Events: eventWriter, Approvals: approvalService}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/tenant-access/grants", h.handleCreateGrant)
	mux.HandleFunc("GET /v1/tenant-access/grants/", h.handleListGrants)
	mux.HandleFunc("POST /v1/tenant-access/grants/", h.handleGrantAction)
	mux.HandleFunc("POST /v1/tenant-access/delegations", h.handleCreateDelegation)
	mux.HandleFunc("GET /v1/tenant-access/delegations/", h.handleListDelegations)
	mux.HandleFunc("POST /v1/tenant-access/delegations/", h.handleDelegationAction)
	mux.HandleFunc("GET /v1/tenant-access/workspace/", h.handleWorkspace)
	return mux
}

func (h *Handler) handleCreateGrant(w http.ResponseWriter, r *http.Request) {
	var req createGrantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.invalid_json", "message": err.Error()})
		return
	}
	if strings.TrimSpace(req.TenantID) == "" || strings.TrimSpace(req.PrincipalRef) == "" || strings.TrimSpace(req.CapabilityID) == "" || strings.TrimSpace(req.TraceRef) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.missing_required_fields"})
		return
	}
	now := time.Now().UTC()
	validUntil, err := parseOptionalRFC3339(req.ValidUntil)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.invalid_valid_until", "message": err.Error()})
		return
	}
	if !validUntil.IsZero() && validUntil.Before(now) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.valid_until_in_past"})
		return
	}
	requiresApproval := strings.Contains(req.CapabilityID, "approval") || strings.Contains(req.CapabilityID, "recovery") || strings.Contains(req.CapabilityID, "restricted")
	state := access.StateActive
	if requiresApproval {
		state = access.StateBlocked
	}
	grant := access.CapabilityGrant{
		GrantID: fmt.Sprintf("grant-%d", now.UnixNano()), TenantID: strings.TrimSpace(req.TenantID), PrincipalRef: strings.TrimSpace(req.PrincipalRef), PrincipalType: strings.TrimSpace(req.PrincipalType), CapabilityID: strings.TrimSpace(req.CapabilityID), ScopeRef: strings.TrimSpace(req.ScopeRef), AllowedActions: cleanStrings(req.AllowedActions), DeniedActions: cleanStrings(req.DeniedActions), RequiresApproval: requiresApproval, Justification: strings.TrimSpace(req.Justification), TraceRef: strings.TrimSpace(req.TraceRef), State: state, ValidFrom: now, ValidUntil: validUntil, CreatedAt: now, UpdatedAt: now,
	}
	if requiresApproval && h.Approvals != nil {
		approval := approvals.Request{ApprovalRequestID: fmt.Sprintf("approval-access-grant-%d", now.UnixNano()), ExecutionID: grant.GrantID, ContractID: grant.GrantID, TenantID: grant.TenantID, TraceID: grant.TraceRef, State: approvals.StateAwaitingApproval, Mode: "pre_execution", ReasonCodes: []string{"grant.requires_approval"}, CreatedAt: now, UpdatedAt: now}
		if err := h.Approvals.Create(approval); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.approval_create_failed", "message": err.Error()})
			return
		}
		grant.ApprovalRequestID = approval.ApprovalRequestID
	}
	if err := h.Store.SaveGrant(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()), EventType: "tenant_access.grant_created", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "principal_ref": grant.PrincipalRef, "capability_id": grant.CapabilityID, "state": grant.State}})
	if grant.ApprovalRequestID != "" {
		h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()+1), EventType: "tenant_access.grant_awaiting_approval", TenantID: grant.TenantID, TraceID: grant.TraceRef, ApprovalRequestID: grant.ApprovalRequestID, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "approval_request_id": grant.ApprovalRequestID}})
	}
	writeJSON(w, http.StatusCreated, grant)
}

func (h *Handler) handleListGrants(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/tenant-access/grants/"), "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.missing_tenant_id"})
		return
	}
	grants, err := h.Store.ListGrantsByTenant(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.list_failed", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant_id": tenantID, "grants": grants})
}

func (h *Handler) handleCreateDelegation(w http.ResponseWriter, r *http.Request) {
	var req createDelegationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.invalid_json", "message": err.Error()})
		return
	}
	if strings.TrimSpace(req.TenantID) == "" || strings.TrimSpace(req.SourcePrincipal) == "" || strings.TrimSpace(req.TargetPrincipal) == "" || strings.TrimSpace(req.ScopeRef) == "" || strings.TrimSpace(req.TraceRef) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.missing_required_fields"})
		return
	}
	now := time.Now().UTC()
	validUntil, err := parseOptionalRFC3339(req.ValidUntil)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.invalid_valid_until", "message": err.Error()})
		return
	}
	if !validUntil.IsZero() && validUntil.Before(now) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.valid_until_in_past"})
		return
	}
	requiresApproval := strings.Contains(req.ScopeRef, "approval") || strings.Contains(req.ScopeRef, "recovery") || req.MaxDepth > 1 || req.CanRedelegate
	state := access.StateActive
	if requiresApproval {
		state = access.StateBlocked
	}
	grant := access.DelegationGrant{
		GrantID: fmt.Sprintf("delegation-%d", now.UnixNano()), TenantID: strings.TrimSpace(req.TenantID), SourcePrincipal: strings.TrimSpace(req.SourcePrincipal), TargetPrincipal: strings.TrimSpace(req.TargetPrincipal), AuthoritySource: strings.TrimSpace(req.AuthoritySource), ScopeType: strings.TrimSpace(req.ScopeType), ScopeRef: strings.TrimSpace(req.ScopeRef), AllowedActions: cleanStrings(req.AllowedActions), DeniedActions: cleanStrings(req.DeniedActions), RequiresApproval: requiresApproval, CanRedelegate: req.CanRedelegate, MaxDepth: req.MaxDepth, Justification: strings.TrimSpace(req.Justification), TraceRef: strings.TrimSpace(req.TraceRef), State: state, ValidFrom: now, ValidUntil: validUntil, CreatedAt: now, UpdatedAt: now,
	}
	if requiresApproval && h.Approvals != nil {
		approval := approvals.Request{ApprovalRequestID: fmt.Sprintf("approval-access-delegation-%d", now.UnixNano()), ExecutionID: grant.GrantID, ContractID: grant.GrantID, TenantID: grant.TenantID, TraceID: grant.TraceRef, State: approvals.StateAwaitingApproval, Mode: "pre_execution", ReasonCodes: []string{"delegation.requires_approval"}, CreatedAt: now, UpdatedAt: now}
		if err := h.Approvals.Create(approval); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.approval_create_failed", "message": err.Error()})
			return
		}
		grant.ApprovalRequestID = approval.ApprovalRequestID
	}
	if err := h.Store.SaveDelegation(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()+1), EventType: "tenant_access.delegation_created", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "source_principal": grant.SourcePrincipal, "target_principal": grant.TargetPrincipal, "scope_ref": grant.ScopeRef, "state": grant.State}})
	if grant.ApprovalRequestID != "" {
		h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()+2), EventType: "tenant_access.delegation_awaiting_approval", TenantID: grant.TenantID, TraceID: grant.TraceRef, ApprovalRequestID: grant.ApprovalRequestID, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "approval_request_id": grant.ApprovalRequestID}})
	}
	writeJSON(w, http.StatusCreated, grant)
}

func (h *Handler) handleListDelegations(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/tenant-access/delegations/"), "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.missing_tenant_id"})
		return
	}
	delegations, err := h.Store.ListDelegationsByTenant(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.list_failed", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant_id": tenantID, "delegations": delegations})
}

func (h *Handler) handleWorkspace(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/tenant-access/workspace/"), "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "workspace.missing_tenant_id"})
		return
	}
	grants, err := h.Store.ListGrantsByTenant(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.grants_failed", "message": err.Error()})
		return
	}
	delegations, err := h.Store.ListDelegationsByTenant(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "workspace.delegations_failed", "message": err.Error()})
		return
	}
	workspace := buildAccessAdminWorkspace(tenantID, grants, delegations)
	writeJSON(w, http.StatusOK, workspace)
}

func (h *Handler) handleGrantAction(w http.ResponseWriter, r *http.Request) {
	grantID, action := splitActionPath(strings.TrimPrefix(r.URL.Path, "/v1/tenant-access/grants/"))
	if grantID == "" || action == "" {
		return
	}
	grant, found, err := h.Store.GetGrantByID(grantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "grant.not_found"})
		return
	}
	switch action {
	case "approve":
		h.handleApproveGrant(w, r, grant)
	case "revoke":
		h.handleRevokeGrant(w, r, grant)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.invalid_action"})
	}
}

func (h *Handler) handleDelegationAction(w http.ResponseWriter, r *http.Request) {
	grantID, action := splitActionPath(strings.TrimPrefix(r.URL.Path, "/v1/tenant-access/delegations/"))
	if grantID == "" || action == "" {
		return
	}
	grant, found, err := h.Store.GetDelegationByID(grantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "delegation.not_found"})
		return
	}
	switch action {
	case "approve":
		h.handleApproveDelegation(w, r, grant)
	case "revoke":
		h.handleRevokeDelegation(w, r, grant)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.invalid_action"})
	}
}

func (h *Handler) handleApproveGrant(w http.ResponseWriter, r *http.Request, grant access.CapabilityGrant) {
	var req approvalActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.invalid_json", "message": err.Error()})
		return
	}
	if grant.State != access.StateBlocked {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "grant.invalid_state_for_approve", "state": grant.State})
		return
	}
	if h.Approvals == nil || grant.ApprovalRequestID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.approval_not_available"})
		return
	}
	_, err := h.Approvals.Decide(grant.ApprovalRequestID, approvals.Decision{State: approvals.StateReleased, DecidedBySubjectID: req.DecidedBySubjectID, DecisionComment: req.DecisionComment, DecisionReasonCodes: req.DecisionReasonCodes, DecidedAt: time.Now().UTC()})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.approval_failed", "message": err.Error()})
		return
	}
	grant.State = access.StateActive
	grant.UpdatedAt = time.Now().UTC()
	if err := h.Store.SaveGrant(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()), EventType: "tenant_access.grant_released", TenantID: grant.TenantID, TraceID: grant.TraceRef, ApprovalRequestID: grant.ApprovalRequestID, OccurredAt: time.Now().UTC(), Payload: map[string]any{"grant_id": grant.GrantID, "state": grant.State}})
	writeJSON(w, http.StatusOK, grant)
}

func (h *Handler) handleApproveDelegation(w http.ResponseWriter, r *http.Request, grant access.DelegationGrant) {
	var req approvalActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.invalid_json", "message": err.Error()})
		return
	}
	if grant.State != access.StateBlocked {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "delegation.invalid_state_for_approve", "state": grant.State})
		return
	}
	if h.Approvals == nil || grant.ApprovalRequestID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.approval_not_available"})
		return
	}
	_, err := h.Approvals.Decide(grant.ApprovalRequestID, approvals.Decision{State: approvals.StateReleased, DecidedBySubjectID: req.DecidedBySubjectID, DecisionComment: req.DecisionComment, DecisionReasonCodes: req.DecisionReasonCodes, DecidedAt: time.Now().UTC()})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.approval_failed", "message": err.Error()})
		return
	}
	grant.State = access.StateActive
	grant.UpdatedAt = time.Now().UTC()
	if err := h.Store.SaveDelegation(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", time.Now().UTC().UnixNano()), EventType: "tenant_access.delegation_released", TenantID: grant.TenantID, TraceID: grant.TraceRef, ApprovalRequestID: grant.ApprovalRequestID, OccurredAt: time.Now().UTC(), Payload: map[string]any{"grant_id": grant.GrantID, "state": grant.State}})
	writeJSON(w, http.StatusOK, grant)
}

func (h *Handler) handleRevokeGrant(w http.ResponseWriter, r *http.Request, grant access.CapabilityGrant) {
	var req revokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.invalid_json", "message": err.Error()})
		return
	}
	if strings.TrimSpace(req.RevokedBy) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "grant.missing_revoked_by"})
		return
	}
	if grant.State == access.StateRevoked {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "grant.already_revoked", "state": grant.State})
		return
	}
	now := time.Now().UTC()
	grant.State = access.StateRevoked
	grant.RevokedBy = strings.TrimSpace(req.RevokedBy)
	grant.RevokedAt = now
	grant.UpdatedAt = now
	if req.Justification != "" {
		grant.Justification = strings.TrimSpace(req.Justification)
	}
	if err := h.Store.SaveGrant(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()), EventType: "tenant_access.grant_revoked", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "revoked_by": grant.RevokedBy, "state": grant.State}})
	writeJSON(w, http.StatusOK, grant)
}

func (h *Handler) handleRevokeDelegation(w http.ResponseWriter, r *http.Request, grant access.DelegationGrant) {
	var req revokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.invalid_json", "message": err.Error()})
		return
	}
	if strings.TrimSpace(req.RevokedBy) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "delegation.missing_revoked_by"})
		return
	}
	if grant.State == access.StateRevoked {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "delegation.already_revoked", "state": grant.State})
		return
	}
	now := time.Now().UTC()
	grant.State = access.StateRevoked
	grant.RevokedBy = strings.TrimSpace(req.RevokedBy)
	grant.RevokedAt = now
	grant.UpdatedAt = now
	if req.Justification != "" {
		grant.Justification = strings.TrimSpace(req.Justification)
	}
	if err := h.Store.SaveDelegation(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()), EventType: "tenant_access.delegation_revoked", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "revoked_by": grant.RevokedBy, "state": grant.State}})
	writeJSON(w, http.StatusOK, grant)
}

func splitActionPath(path string) (string, string) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func buildAccessAdminWorkspace(tenantID string, grants []access.CapabilityGrant, delegations []access.DelegationGrant) accessAdminWorkspace {
	blockedItems := 0
	revokedItems := 0
	expiredItems := 0
	activeGrants := 0
	blockedGrants := 0
	blockedGrantIDs := make([]string, 0)
	revokedGrantIDs := make([]string, 0)
	expiredGrantIDs := make([]string, 0)
	approvalRequiredGrants := 0
	sensitiveCapabilities := make([]string, 0)
	activeDelegations := 0
	blockedDelegations := 0
	blockedDelegationIDs := make([]string, 0)
	revokedDelegationIDs := make([]string, 0)
	expiredDelegationIDs := make([]string, 0)
	redelegableDelegations := 0
	sensitiveDelegations := make([]string, 0)
	approvalRequests := 0
	revocations := 0
	now := time.Now().UTC()
	for _, grant := range grants {
		switch grant.State {
		case access.StateActive:
			activeGrants++
		case access.StateBlocked:
			blockedGrants++
			blockedItems++
			blockedGrantIDs = append(blockedGrantIDs, grant.GrantID)
		case access.StateRevoked:
			revokedItems++
			revocations++
			revokedGrantIDs = append(revokedGrantIDs, grant.GrantID)
		}
		if grant.RequiresApproval {
			approvalRequiredGrants++
		}
		if grant.ApprovalRequestID != "" {
			approvalRequests++
		}
		if !grant.ValidUntil.IsZero() && grant.ValidUntil.Before(now) {
			expiredItems++
			expiredGrantIDs = append(expiredGrantIDs, grant.GrantID)
		}
		if strings.Contains(grant.CapabilityID, "approval") || strings.Contains(grant.CapabilityID, "recovery") || strings.Contains(grant.CapabilityID, "restricted") {
			sensitiveCapabilities = append(sensitiveCapabilities, grant.CapabilityID)
		}
	}
	for _, grant := range delegations {
		switch grant.State {
		case access.StateActive:
			activeDelegations++
		case access.StateBlocked:
			blockedDelegations++
			blockedItems++
			blockedDelegationIDs = append(blockedDelegationIDs, grant.GrantID)
		case access.StateRevoked:
			revokedItems++
			revocations++
			revokedDelegationIDs = append(revokedDelegationIDs, grant.GrantID)
		}
		if grant.ApprovalRequestID != "" {
			approvalRequests++
		}
		if !grant.ValidUntil.IsZero() && grant.ValidUntil.Before(now) {
			expiredItems++
			expiredDelegationIDs = append(expiredDelegationIDs, grant.GrantID)
		}
		if grant.CanRedelegate {
			redelegableDelegations++
		}
		if grant.RequiresApproval || grant.MaxDepth > 1 || grant.CanRedelegate {
			sensitiveDelegations = append(sensitiveDelegations, grant.GrantID)
		}
	}
	recommendedNext := []string{"review blocked items", "review approval-sensitive grants", "review delegation depth and redelegation"}
	promotionAdvice := "ready_for_governed_access_changes"
	if blockedItems > 0 {
		promotionAdvice = "resolve_blocked_grants_or_delegations_before_promoting_sensitive_changes"
		recommendedNext = []string{"resolve blocked grants or delegations", "review pending approvals", "confirm sensitive access changes before promotion"}
	} else if expiredItems > 0 {
		promotionAdvice = "review_expired_grants_or_delegations_before_trusting_current_access_state"
		recommendedNext = []string{"review expired grants or delegations", "confirm whether they should be revoked or renewed", "check sensitive access paths after expiration"}
	} else if revocations > 0 {
		promotionAdvice = "review_recent_revocations_before_promoting_new_access_changes"
		recommendedNext = []string{"review recent revocations", "confirm remaining active access is still correct", "check sensitive delegation paths"}
	} else if approvalRequests > 0 {
		promotionAdvice = "review_approval_sensitive_changes_before_promoting_broader_access"
		recommendedNext = []string{"review approval-sensitive grants", "review delegation depth", "confirm no accidental authority escalation"}
	}
	summaryText := fmt.Sprintf("tenant %s currently has %d active grants, %d active delegations, %d blocked items, %d revoked items and %d expired items", tenantID, activeGrants, activeDelegations, blockedItems, revokedItems, expiredItems)
	grantSummary := fmt.Sprintf("%d grants active, %d blocked, %d approval-sensitive, %d expired", activeGrants, blockedGrants, approvalRequiredGrants, len(expiredGrantIDs))
	delegationSummary := fmt.Sprintf("%d delegations active, %d blocked, %d redelegable, %d expired", activeDelegations, blockedDelegations, redelegableDelegations, len(expiredDelegationIDs))
	governanceSummary := fmt.Sprintf("%d approval requests and %d revocations shape the current access governance state", approvalRequests, revocations)
	return accessAdminWorkspace{
		TenantID: tenantID,
		Summary: accessSummaryCard{
			TotalGrants:      len(grants),
			TotalDelegations: len(delegations),
			BlockedItems:     blockedItems,
			RevokedItems:     revokedItems,
			ExpiredItems:     expiredItems,
			RecommendedNext:  recommendedNext,
			Summary:          summaryText,
		},
		Grants: accessGrantCard{
			ActiveGrants:          activeGrants,
			BlockedGrants:         blockedGrants,
			ApprovalRequired:      approvalRequiredGrants,
			BlockedGrantIDs:       blockedGrantIDs,
			RevokedGrantIDs:       revokedGrantIDs,
			ExpiredGrantIDs:       expiredGrantIDs,
			SensitiveCapabilities: uniqueStrings(sensitiveCapabilities),
			Summary:               grantSummary,
		},
		Delegations: accessDelegationCard{
			ActiveDelegations:      activeDelegations,
			BlockedDelegations:     blockedDelegations,
			RedelegableDelegations: redelegableDelegations,
			BlockedDelegationIDs:   blockedDelegationIDs,
			RevokedDelegationIDs:   revokedDelegationIDs,
			ExpiredDelegationIDs:   expiredDelegationIDs,
			SensitiveDelegations:   uniqueStrings(sensitiveDelegations),
			Summary:                delegationSummary,
		},
		Governance: accessGovernanceCard{
			ApprovalRequestCount: approvalRequests,
			RevocationCount:      revocations,
			Guardrails:           []string{"approval-sensitive access stays governed", "revocation remains explicit and auditable", "delegation does not bypass tenant policies"},
			Summary:              governanceSummary,
		},
		Impact: accessImpactCard{
			ApprovalSensitiveAreas:   []string{"restricted capabilities", "redelegable delegations", "high-sensitivity access changes"},
			RevocationSensitiveAreas: []string{"active grants", "active delegations", "delegated approval paths"},
			PromotionAdvice:          promotionAdvice,
		},
		Boundary: "access_admin_surface_reads_grants_and_guides_governed_authority_changes",
	}
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func parseOptionalRFC3339(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func cleanStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func (h *Handler) appendEvent(record events.Record) {
	if h.Events != nil {
		_ = h.Events.Append(record)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
