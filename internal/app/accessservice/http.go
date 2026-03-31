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
	requiresApproval := strings.Contains(req.CapabilityID, "approval") || strings.Contains(req.CapabilityID, "recovery") || strings.Contains(req.CapabilityID, "restricted")
	state := access.StateActive
	if requiresApproval {
		state = access.StateBlocked
	}
	grant := access.CapabilityGrant{
		GrantID: fmt.Sprintf("grant-%d", now.UnixNano()), TenantID: strings.TrimSpace(req.TenantID), PrincipalRef: strings.TrimSpace(req.PrincipalRef), PrincipalType: strings.TrimSpace(req.PrincipalType), CapabilityID: strings.TrimSpace(req.CapabilityID), ScopeRef: strings.TrimSpace(req.ScopeRef), AllowedActions: cleanStrings(req.AllowedActions), DeniedActions: cleanStrings(req.DeniedActions), RequiresApproval: requiresApproval, Justification: strings.TrimSpace(req.Justification), TraceRef: strings.TrimSpace(req.TraceRef), State: state, ValidFrom: now, CreatedAt: now, UpdatedAt: now,
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
	requiresApproval := strings.Contains(req.ScopeRef, "approval") || strings.Contains(req.ScopeRef, "recovery") || req.MaxDepth > 1 || req.CanRedelegate
	state := access.StateActive
	if requiresApproval {
		state = access.StateBlocked
	}
	grant := access.DelegationGrant{
		GrantID: fmt.Sprintf("delegation-%d", now.UnixNano()), TenantID: strings.TrimSpace(req.TenantID), SourcePrincipal: strings.TrimSpace(req.SourcePrincipal), TargetPrincipal: strings.TrimSpace(req.TargetPrincipal), AuthoritySource: strings.TrimSpace(req.AuthoritySource), ScopeType: strings.TrimSpace(req.ScopeType), ScopeRef: strings.TrimSpace(req.ScopeRef), AllowedActions: cleanStrings(req.AllowedActions), DeniedActions: cleanStrings(req.DeniedActions), RequiresApproval: requiresApproval, CanRedelegate: req.CanRedelegate, MaxDepth: req.MaxDepth, Justification: strings.TrimSpace(req.Justification), TraceRef: strings.TrimSpace(req.TraceRef), State: state, ValidFrom: now, CreatedAt: now, UpdatedAt: now,
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
