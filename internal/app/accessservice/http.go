package accessservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/access"
	"opita-sync-framework/internal/engine/events"
)

type Store interface {
	SaveGrant(grant access.CapabilityGrant) error
	ListGrantsByTenant(tenantID string) ([]access.CapabilityGrant, error)
	SaveDelegation(grant access.DelegationGrant) error
	ListDelegationsByTenant(tenantID string) ([]access.DelegationGrant, error)
}

type EventWriter interface {
	Append(record events.Record) error
}

type Handler struct {
	Store  Store
	Events EventWriter
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

func NewHandler(store Store, eventWriter EventWriter) *Handler {
	return &Handler{Store: store, Events: eventWriter}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/tenant-access/grants", h.handleCreateGrant)
	mux.HandleFunc("GET /v1/tenant-access/grants/", h.handleListGrants)
	mux.HandleFunc("POST /v1/tenant-access/delegations", h.handleCreateDelegation)
	mux.HandleFunc("GET /v1/tenant-access/delegations/", h.handleListDelegations)
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
	if err := h.Store.SaveGrant(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "grant.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()), EventType: "tenant_access.grant_created", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "principal_ref": grant.PrincipalRef, "capability_id": grant.CapabilityID, "state": grant.State}})
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
	if err := h.Store.SaveDelegation(grant); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "delegation.save_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{EventID: fmt.Sprintf("event-%d", now.UnixNano()+1), EventType: "tenant_access.delegation_created", TenantID: grant.TenantID, TraceID: grant.TraceRef, OccurredAt: now, Payload: map[string]any{"grant_id": grant.GrantID, "source_principal": grant.SourcePrincipal, "target_principal": grant.TargetPrincipal, "scope_ref": grant.ScopeRef, "state": grant.State}})
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
