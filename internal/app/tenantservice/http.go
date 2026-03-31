package tenantservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/tenant"
)

type Store interface {
	Save(record tenant.BootstrapRecord) error
	GetByTenantID(tenantID string) (tenant.BootstrapRecord, bool, error)
}

type EventWriter interface {
	Append(record events.Record) error
}

type Handler struct {
	Store  Store
	Events EventWriter
}

type tenantAdminWorkspace struct {
	TenantID   string               `json:"tenant_id"`
	TenantName string               `json:"tenant_name"`
	State      tenant.State         `json:"state"`
	Summary    tenantSummaryCard    `json:"summary"`
	Governance tenantGovernanceCard `json:"governance"`
	Catalog    tenantCatalogCard    `json:"catalog"`
	Connectors tenantConnectorCard  `json:"connectors"`
	Impact     tenantImpactCard     `json:"impact"`
	Boundary   string               `json:"boundary"`
}

type tenantSummaryCard struct {
	Operable        bool     `json:"operable"`
	FailedChecks    []string `json:"failed_checks"`
	RecommendedNext []string `json:"recommended_next_actions"`
	Summary         string   `json:"summary"`
}

type tenantGovernanceCard struct {
	PolicyProfile         string   `json:"policy_profile"`
	ApprovalProfile       string   `json:"approval_profile"`
	ClassificationProfile string   `json:"classification_profile"`
	Guardrails            []string `json:"guardrails"`
	Summary               string   `json:"summary"`
}

type tenantCatalogCard struct {
	VisibleCapabilities    int      `json:"visible_capabilities"`
	ApprovalSensitive      int      `json:"approval_sensitive_capabilities"`
	RestrictedCapabilities int      `json:"restricted_capabilities"`
	HighRiskCapabilities   []string `json:"high_risk_capabilities"`
	Summary                string   `json:"summary"`
}

type tenantConnectorCard struct {
	EnabledConnectors    int      `json:"enabled_connectors"`
	RestrictedConnectors []string `json:"restricted_connectors"`
	ExecutionConnectors  []string `json:"execution_connectors"`
	Summary              string   `json:"summary"`
}

type tenantImpactCard struct {
	RequiresApprovalAreas []string `json:"requires_approval_areas"`
	SensitiveAreas        []string `json:"sensitive_areas"`
	PromotionAdvice       string   `json:"promotion_advice"`
}

type bootstrapRequest struct {
	TenantID                 string         `json:"tenant_id"`
	TenantName               string         `json:"tenant_name"`
	AdminSubjectID           string         `json:"admin_subject_id"`
	InitialCatalogRefs       []string       `json:"initial_catalog_refs"`
	InitialConnectorRefs     []string       `json:"initial_connector_refs"`
	PolicyProfileRef         string         `json:"policy_profile_ref"`
	ApprovalProfileRef       string         `json:"approval_profile_ref"`
	ClassificationProfileRef string         `json:"classification_profile_ref"`
	ContextSeed              map[string]any `json:"context_seed,omitempty"`
}

func NewHandler(store Store, eventWriter EventWriter) *Handler {
	return &Handler{Store: store, Events: eventWriter}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/tenants/bootstrap", h.handleBootstrap)
	mux.HandleFunc("GET /v1/tenants/", h.handleGetTenant)
	mux.HandleFunc("GET /v1/tenants-catalog/", h.handleGetTenantCatalog)
	mux.HandleFunc("GET /v1/tenants-connectors/", h.handleGetTenantConnectors)
	mux.HandleFunc("GET /v1/tenant-admin/workspace/", h.handleGetTenantAdminWorkspace)
	return mux
}

func (h *Handler) handleBootstrap(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "tenant.service_not_ready"})
		return
	}
	var req bootstrapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.invalid_json", "message": err.Error()})
		return
	}
	policyBaseline, policyOK := resolvePolicyBaseline(req.PolicyProfileRef)
	approvalBaseline, approvalOK := resolveApprovalBaseline(req.ApprovalProfileRef)
	classificationBaseline, classificationOK := resolveClassificationBaseline(req.ClassificationProfileRef)
	catalogProjection, catalogOK := resolveCatalogProjection(req.InitialCatalogRefs)
	connectorProjection, connectorOK := resolveConnectorProjection(req.InitialConnectorRefs)
	checks := validateBootstrap(req, policyOK, approvalOK, classificationOK, catalogOK, connectorOK)
	operable := allChecksPassed(checks)
	state := tenant.StateBlocked
	if operable {
		state = tenant.StateOperable
	}
	now := time.Now().UTC()
	record := tenant.BootstrapRecord{
		TenantID:                  strings.TrimSpace(req.TenantID),
		TenantName:                strings.TrimSpace(req.TenantName),
		AdminSubjectID:            strings.TrimSpace(req.AdminSubjectID),
		InitialCatalogRefs:        cleanStrings(req.InitialCatalogRefs),
		InitialConnectorRefs:      cleanStrings(req.InitialConnectorRefs),
		PolicyProfileRef:          strings.TrimSpace(req.PolicyProfileRef),
		ApprovalProfileRef:        strings.TrimSpace(req.ApprovalProfileRef),
		ClassificationProfileRef:  strings.TrimSpace(req.ClassificationProfileRef),
		PolicyBaseline:            policyBaseline,
		ApprovalBaseline:          approvalBaseline,
		ClassificationBaseline:    classificationBaseline,
		CatalogProjection:         catalogProjection,
		ConnectorProjection:       connectorProjection,
		ContextSeed:               req.ContextSeed,
		State:                     state,
		BootstrapRecordRef:        fmt.Sprintf("tenant-bootstrap-%s", strings.TrimSpace(req.TenantID)),
		PolicyProfileAppliedRef:   fmt.Sprintf("tenant-policy-%s", strings.TrimSpace(req.TenantID)),
		ApprovalProfileAppliedRef: fmt.Sprintf("tenant-approval-%s", strings.TrimSpace(req.TenantID)),
		ClassificationAppliedRef:  fmt.Sprintf("tenant-classification-%s", strings.TrimSpace(req.TenantID)),
		CatalogProjectionRef:      fmt.Sprintf("tenant-catalog-%s", strings.TrimSpace(req.TenantID)),
		ConnectorProjectionRef:    fmt.Sprintf("tenant-connectors-%s", strings.TrimSpace(req.TenantID)),
		OperabilityReport: tenant.OperabilityReport{
			TenantID: strings.TrimSpace(req.TenantID),
			Operable: operable,
			Checks:   checks,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.Store.Save(record); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.bootstrap_failed", "message": err.Error()})
		return
	}
	h.appendEvent(events.Record{
		EventID:    fmt.Sprintf("event-%d", now.UnixNano()),
		EventType:  "tenant.bootstrap_completed",
		TenantID:   record.TenantID,
		OccurredAt: now,
		Payload: map[string]any{
			"tenant_bootstrap_record":               record.BootstrapRecordRef,
			"tenant_policy_profile_applied":         record.PolicyProfileAppliedRef,
			"tenant_approval_profile_applied":       record.ApprovalProfileAppliedRef,
			"tenant_classification_profile_applied": record.ClassificationAppliedRef,
			"policy_baseline":                       record.PolicyBaseline,
			"approval_baseline":                     record.ApprovalBaseline,
			"classification_baseline":               record.ClassificationBaseline,
			"catalog_projection":                    record.CatalogProjection,
			"connector_projection":                  record.ConnectorProjection,
			"tenant_catalog_projection":             record.CatalogProjectionRef,
			"tenant_connector_projection":           record.ConnectorProjectionRef,
			"tenant_operability_report":             record.OperabilityReport,
			"state":                                 record.State,
		},
	})
	status := http.StatusCreated
	if !operable {
		status = http.StatusAccepted
	}
	writeJSON(w, status, record)
}

func (h *Handler) handleGetTenant(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "tenant.service_not_ready"})
		return
	}
	tenantID := strings.TrimPrefix(r.URL.Path, "/v1/tenants/")
	tenantID = strings.TrimSuffix(tenantID, "/")
	if tenantID == "" || tenantID == "bootstrap" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.missing_id"})
		return
	}
	record, found, err := h.Store.GetByTenantID(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "tenant.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "tenant.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (h *Handler) handleGetTenantCatalog(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimPrefix(r.URL.Path, "/v1/tenants-catalog/")
	tenantID = strings.TrimSuffix(tenantID, "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.missing_id"})
		return
	}
	record, found, err := h.Store.GetByTenantID(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "tenant.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "tenant.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant_id": tenantID, "catalog_projection": record.CatalogProjection})
}

func (h *Handler) handleGetTenantConnectors(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimPrefix(r.URL.Path, "/v1/tenants-connectors/")
	tenantID = strings.TrimSuffix(tenantID, "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.missing_id"})
		return
	}
	record, found, err := h.Store.GetByTenantID(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "tenant.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "tenant.not_found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant_id": tenantID, "connector_projection": record.ConnectorProjection})
}

func (h *Handler) handleGetTenantAdminWorkspace(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimPrefix(r.URL.Path, "/v1/tenant-admin/workspace/")
	tenantID = strings.TrimSuffix(tenantID, "/")
	if tenantID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "tenant.missing_id"})
		return
	}
	record, found, err := h.Store.GetByTenantID(tenantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "tenant.lookup_failed", "message": err.Error()})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "tenant.not_found"})
		return
	}
	workspace := buildTenantAdminWorkspace(record)
	writeJSON(w, http.StatusOK, workspace)
}

func validateBootstrap(req bootstrapRequest, policyOK, approvalOK, classificationOK, catalogOK, connectorOK bool) []tenant.OperabilityCheck {
	checks := []tenant.OperabilityCheck{
		makeCheck("tenant_identity", strings.TrimSpace(req.TenantID) != "" && strings.TrimSpace(req.TenantName) != "", "tenant_id and tenant_name are required"),
		makeCheck("admin_identity", strings.TrimSpace(req.AdminSubjectID) != "", "admin_subject_id is required"),
		makeCheck("policy_profile", strings.TrimSpace(req.PolicyProfileRef) != "" && policyOK, "policy_profile_ref is required and must be supported"),
		makeCheck("approval_profile", strings.TrimSpace(req.ApprovalProfileRef) != "" && approvalOK, "approval_profile_ref is required and must be supported"),
		makeCheck("classification_profile", strings.TrimSpace(req.ClassificationProfileRef) != "" && classificationOK, "classification_profile_ref is required and must be supported"),
		makeCheck("catalog_projection", len(cleanStrings(req.InitialCatalogRefs)) > 0 && catalogOK, "initial_catalog_refs must include only supported entries"),
		makeCheck("connector_projection", len(cleanStrings(req.InitialConnectorRefs)) > 0 && connectorOK, "initial_connector_refs must include only supported entries"),
		makeCheck("context_seed", req.ContextSeed != nil, "context_seed is required"),
	}
	functionalReady := true
	for _, check := range checks {
		if !check.Passed {
			functionalReady = false
			break
		}
	}
	checks = append(checks, tenant.OperabilityCheck{CheckName: "governed_corridor_readiness", Passed: functionalReady, Message: "minimum tenant baseline required for intake -> proposal -> preview -> compile -> approval -> inspection readiness"})
	return checks
}

func resolveCatalogProjection(refs []string) ([]tenant.CatalogCapability, bool) {
	definitions := map[string]tenant.CatalogCapability{
		"tenant.intake.capture_intent":                {CapabilityID: "tenant.intake.capture_intent", Name: "Capture intent", ResultType: "plan", Sensitivity: "internal", RiskLevel: "low", RequiresApproval: false, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Assignable: true, Administrable: true}},
		"tenant.proposal.create_change_draft":         {CapabilityID: "tenant.proposal.create_change_draft", Name: "Create change draft", ResultType: "change_proposal", Sensitivity: "internal", RiskLevel: "medium", RequiresApproval: false, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Assignable: true, Administrable: true}},
		"tenant.preview.run_governance_preview":       {CapabilityID: "tenant.preview.run_governance_preview", Name: "Run governance preview", ResultType: "report", Sensitivity: "internal", RiskLevel: "medium", RequiresApproval: false, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Assignable: true, Administrable: true}},
		"tenant.execution.compile_governed_intent":    {CapabilityID: "tenant.execution.compile_governed_intent", Name: "Compile governed intent", ResultType: "execution", Sensitivity: "internal", RiskLevel: "medium", RequiresApproval: true, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Approvable: true, Assignable: true, Administrable: true}},
		"tenant.execution.inspect_run":                {CapabilityID: "tenant.execution.inspect_run", Name: "Inspect run", ResultType: "inspection", Sensitivity: "internal", RiskLevel: "low", RequiresApproval: false, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Assignable: true, Administrable: true}},
		"tenant.approval.release_blocked_execution":   {CapabilityID: "tenant.approval.release_blocked_execution", Name: "Release blocked execution", ResultType: "governance_decision", Sensitivity: "restricted", RiskLevel: "high", RequiresApproval: true, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Approvable: true, Assignable: true, Administrable: true}},
		"tenant.recovery.resume_after_approval":       {CapabilityID: "tenant.recovery.resume_after_approval", Name: "Resume after approval", ResultType: "execution", Sensitivity: "restricted", RiskLevel: "high", RequiresApproval: true, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Approvable: true, Assignable: true, Administrable: true}},
		"tenant.recovery.request_manual_compensation": {CapabilityID: "tenant.recovery.request_manual_compensation", Name: "Request manual compensation", ResultType: "governance_decision", Sensitivity: "restricted", RiskLevel: "high", RequiresApproval: true, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Approvable: true, Assignable: true, Administrable: true}},
		"tenant.maintenance.request_human_review":     {CapabilityID: "tenant.maintenance.request_human_review", Name: "Request human review", ResultType: "governance_decision", Sensitivity: "internal", RiskLevel: "medium", RequiresApproval: true, Status: "active", TenantVisibility: tenant.Visibility{Visible: true, Usable: true, Approvable: true, Assignable: true, Administrable: true}},
	}
	out := make([]tenant.CatalogCapability, 0, len(refs))
	for _, ref := range cleanStrings(refs) {
		capability, ok := definitions[ref]
		if !ok {
			return nil, false
		}
		out = append(out, capability)
	}
	return out, len(out) > 0
}

func resolveConnectorProjection(refs []string) ([]tenant.ConnectorProjection, bool) {
	definitions := map[string]tenant.ConnectorProjection{
		"connector.default":              {ConnectorRef: "connector.default", Enabled: true, Scope: "governed_core", SupportedActions: []string{"inspect", "dry_run", "execute", "normalize", "compensate"}},
		"connector.execution.default":    {ConnectorRef: "connector.execution.default", Enabled: true, Scope: "execution", SupportedActions: []string{"inspect", "dry_run", "execute", "normalize", "compensate"}},
		"connector.execution.restricted": {ConnectorRef: "connector.execution.restricted", Enabled: true, Scope: "execution_restricted", SupportedActions: []string{"inspect", "dry_run", "execute", "normalize", "compensate"}},
		"connector.plan.default":         {ConnectorRef: "connector.plan.default", Enabled: true, Scope: "planning", SupportedActions: []string{"inspect", "dry_run", "normalize"}},
	}
	out := make([]tenant.ConnectorProjection, 0, len(refs))
	for _, ref := range cleanStrings(refs) {
		connector, ok := definitions[ref]
		if !ok {
			return nil, false
		}
		out = append(out, connector)
	}
	return out, len(out) > 0
}

func resolvePolicyBaseline(profileRef string) (tenant.PolicyBaseline, bool) {
	switch strings.TrimSpace(profileRef) {
	case "policy.default":
		return tenant.PolicyBaseline{
			ProfileRef:          "policy.default",
			DecisionModel:       "rbac+abac+rebac+policy",
			DefaultAction:       "deny_by_default",
			SupportedDecisions:  []string{"allow", "deny_block", "require_approval", "require_escalation", "restricted_view"},
			GovernanceMode:      "governed",
			ClassificationGuard: true,
		}, true
	case "policy.restrictive-v1":
		return tenant.PolicyBaseline{
			ProfileRef:          "policy.restrictive-v1",
			DecisionModel:       "rbac+abac+rebac+policy",
			DefaultAction:       "deny_by_default",
			SupportedDecisions:  []string{"allow", "deny_block", "require_approval", "require_escalation", "restricted_view"},
			GovernanceMode:      "strict_governed",
			ClassificationGuard: true,
		}, true
	default:
		return tenant.PolicyBaseline{}, false
	}
}

func resolveApprovalBaseline(profileRef string) (tenant.ApprovalBaseline, bool) {
	switch strings.TrimSpace(profileRef) {
	case "approval.default":
		return tenant.ApprovalBaseline{
			ProfileRef:            "approval.default",
			DefaultMode:           "pre_execution",
			RequiredStates:        []string{"awaiting_approval", "released", "rejected", "escalated"},
			DecisionActions:       []string{"release", "reject", "escalate"},
			RequireDecisionActor:  true,
			FingerprintValidation: true,
		}, true
	case "approval.strict-v1":
		return tenant.ApprovalBaseline{
			ProfileRef:            "approval.strict-v1",
			DefaultMode:           "pre_execution",
			RequiredStates:        []string{"awaiting_approval", "released", "rejected", "escalated"},
			DecisionActions:       []string{"release", "reject", "escalate"},
			RequireDecisionActor:  true,
			FingerprintValidation: true,
		}, true
	default:
		return tenant.ApprovalBaseline{}, false
	}
}

func resolveClassificationBaseline(profileRef string) (tenant.ClassificationBaseline, bool) {
	switch strings.TrimSpace(profileRef) {
	case "classification.default":
		return tenant.ClassificationBaseline{
			ProfileRef:         "classification.default",
			DefaultLevel:       "internal",
			AllowedLevels:      []string{"public", "internal", "restricted"},
			RestrictedViewMode: "redacted",
			RedactionRequired:  true,
		}, true
	case "classification.internal-first":
		return tenant.ClassificationBaseline{
			ProfileRef:         "classification.internal-first",
			DefaultLevel:       "internal",
			AllowedLevels:      []string{"internal", "restricted"},
			RestrictedViewMode: "redacted",
			RedactionRequired:  true,
		}, true
	default:
		return tenant.ClassificationBaseline{}, false
	}
}

func makeCheck(name string, passed bool, failMessage string) tenant.OperabilityCheck {
	message := "ok"
	if !passed {
		message = failMessage
	}
	return tenant.OperabilityCheck{CheckName: name, Passed: passed, Message: message}
}

func allChecksPassed(checks []tenant.OperabilityCheck) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
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

func buildTenantAdminWorkspace(record tenant.BootstrapRecord) tenantAdminWorkspace {
	failedChecks := make([]string, 0)
	for _, check := range record.OperabilityReport.Checks {
		if !check.Passed {
			failedChecks = append(failedChecks, check.CheckName)
		}
	}
	highRiskCapabilities := make([]string, 0)
	restrictedCapabilities := 0
	approvalSensitive := 0
	for _, capability := range record.CatalogProjection {
		if capability.RequiresApproval {
			approvalSensitive++
		}
		if capability.Sensitivity == "restricted" {
			restrictedCapabilities++
		}
		if capability.RiskLevel == "high" {
			highRiskCapabilities = append(highRiskCapabilities, capability.CapabilityID)
		}
	}
	restrictedConnectors := make([]string, 0)
	executionConnectors := make([]string, 0)
	for _, connector := range record.ConnectorProjection {
		if strings.Contains(connector.Scope, "restricted") {
			restrictedConnectors = append(restrictedConnectors, connector.ConnectorRef)
		}
		if strings.Contains(connector.Scope, "execution") {
			executionConnectors = append(executionConnectors, connector.ConnectorRef)
		}
	}
	requiresApprovalAreas := make([]string, 0)
	if approvalSensitive > 0 {
		requiresApprovalAreas = append(requiresApprovalAreas, "approval-sensitive capabilities")
	}
	if len(restrictedConnectors) > 0 {
		requiresApprovalAreas = append(requiresApprovalAreas, "restricted connectors")
	}
	sensitiveAreas := make([]string, 0)
	if record.ClassificationBaseline.RedactionRequired {
		sensitiveAreas = append(sensitiveAreas, "classification baseline requires redaction")
	}
	sensitiveAreas = append(sensitiveAreas, highRiskCapabilities...)
	promotionAdvice := "ready_for_governed_changes"
	if !record.OperabilityReport.Operable {
		promotionAdvice = "fix_failed_checks_before_promoting_changes"
	} else if len(requiresApprovalAreas) > 0 {
		promotionAdvice = "use_preview_and_explicit_approval_for_sensitive_changes"
	}
	recommendedNext := []string{"review catalog projection", "review connector projection", "review governance baselines"}
	if !record.OperabilityReport.Operable {
		recommendedNext = append(recommendedNext, "resolve failed operability checks")
	}
	return tenantAdminWorkspace{
		TenantID:   record.TenantID,
		TenantName: record.TenantName,
		State:      record.State,
		Summary: tenantSummaryCard{
			Operable:        record.OperabilityReport.Operable,
			FailedChecks:    failedChecks,
			RecommendedNext: recommendedNext,
			Summary:         fmt.Sprintf("tenant %s is %s with %d visible capabilities and %d connectors", record.TenantID, record.State, len(record.CatalogProjection), len(record.ConnectorProjection)),
		},
		Governance: tenantGovernanceCard{
			PolicyProfile:         record.PolicyBaseline.ProfileRef,
			ApprovalProfile:       record.ApprovalBaseline.ProfileRef,
			ClassificationProfile: record.ClassificationBaseline.ProfileRef,
			Guardrails: []string{
				record.PolicyBaseline.DefaultAction,
				record.ApprovalBaseline.DefaultMode,
				record.ClassificationBaseline.RestrictedViewMode,
			},
			Summary: "tenant governance reflects baseline profiles and should remain governed through preview and approval",
		},
		Catalog: tenantCatalogCard{
			VisibleCapabilities:    len(record.CatalogProjection),
			ApprovalSensitive:      approvalSensitive,
			RestrictedCapabilities: restrictedCapabilities,
			HighRiskCapabilities:   highRiskCapabilities,
			Summary:                fmt.Sprintf("catalog exposes %d capabilities, %d of them approval-sensitive", len(record.CatalogProjection), approvalSensitive),
		},
		Connectors: tenantConnectorCard{
			EnabledConnectors:    len(record.ConnectorProjection),
			RestrictedConnectors: restrictedConnectors,
			ExecutionConnectors:  executionConnectors,
			Summary:              fmt.Sprintf("tenant has %d enabled connectors, %d in execution scope", len(record.ConnectorProjection), len(executionConnectors)),
		},
		Impact: tenantImpactCard{
			RequiresApprovalAreas: requiresApprovalAreas,
			SensitiveAreas:        sensitiveAreas,
			PromotionAdvice:       promotionAdvice,
		},
		Boundary: "tenant_admin_surface_reads_bootstrap_state_and_guides_governed_changes",
	}
}

func (h *Handler) appendEvent(record events.Record) {
	if h.Events == nil {
		return
	}
	_ = h.Events.Append(record)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
