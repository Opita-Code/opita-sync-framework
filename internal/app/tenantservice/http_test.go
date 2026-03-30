package tenantservice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opita-sync-framework/internal/app/tenantservice"
	"opita-sync-framework/internal/platform/memory"
)

func TestBootstrapTenantOperableSuccess(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":                  "tenant-1",
		"tenant_name":                "Tenant Uno",
		"admin_subject_id":           "admin-1",
		"initial_catalog_refs":       []string{"tenant.intake.capture_intent", "tenant.execution.inspect_run"},
		"initial_connector_refs":     []string{"connector.default"},
		"policy_profile_ref":         "policy.default",
		"approval_profile_ref":       "approval.default",
		"classification_profile_ref": "classification.default",
		"context_seed":               map[string]any{"company": "ACME"},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var record map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &record); err != nil {
		t.Fatalf("unmarshal bootstrap record: %v", err)
	}
	if _, ok := record["policy_baseline"]; !ok {
		t.Fatalf("expected policy_baseline in response")
	}
	if _, ok := record["catalog_projection"]; !ok {
		t.Fatalf("expected catalog_projection in response")
	}
	if _, ok := record["connector_projection"]; !ok {
		t.Fatalf("expected connector_projection in response")
	}
}

func TestBootstrapTenantBlockedWhenMinimumHardMissing(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":            "tenant-2",
		"tenant_name":          "Tenant Dos",
		"admin_subject_id":     "admin-2",
		"initial_catalog_refs": []string{"tenant.intake.capture_intent"},
		"context_seed":         map[string]any{"company": "ACME"},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetBootstrappedTenant(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":                  "tenant-3",
		"tenant_name":                "Tenant Tres",
		"admin_subject_id":           "admin-3",
		"initial_catalog_refs":       []string{"tenant.intake.capture_intent"},
		"initial_connector_refs":     []string{"connector.default"},
		"policy_profile_ref":         "policy.default",
		"approval_profile_ref":       "approval.default",
		"classification_profile_ref": "classification.default",
		"context_seed":               map[string]any{"company": "ACME"},
	})
	postReq := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	postW := httptest.NewRecorder()
	h.Routes().ServeHTTP(postW, postReq)
	if postW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", postW.Code)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/v1/tenants/tenant-3", nil)
	getW := httptest.NewRecorder()
	h.Routes().ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", getW.Code, getW.Body.String())
	}
}

func TestBootstrapTenantBlockedWhenProfileUnsupported(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":                  "tenant-4",
		"tenant_name":                "Tenant Cuatro",
		"admin_subject_id":           "admin-4",
		"initial_catalog_refs":       []string{"tenant.intake.capture_intent"},
		"initial_connector_refs":     []string{"connector.default"},
		"policy_profile_ref":         "policy.unknown",
		"approval_profile_ref":       "approval.default",
		"classification_profile_ref": "classification.default",
		"context_seed":               map[string]any{"company": "ACME"},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetTenantCatalogAndConnectors(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":                  "tenant-5",
		"tenant_name":                "Tenant Cinco",
		"admin_subject_id":           "admin-5",
		"initial_catalog_refs":       []string{"tenant.intake.capture_intent", "tenant.execution.inspect_run"},
		"initial_connector_refs":     []string{"connector.default"},
		"policy_profile_ref":         "policy.default",
		"approval_profile_ref":       "approval.default",
		"classification_profile_ref": "classification.default",
		"context_seed":               map[string]any{"company": "ACME"},
	})
	postReq := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	postW := httptest.NewRecorder()
	h.Routes().ServeHTTP(postW, postReq)
	if postW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", postW.Code, postW.Body.String())
	}
	catReq := httptest.NewRequest(http.MethodGet, "/v1/tenants-catalog/tenant-5", nil)
	catW := httptest.NewRecorder()
	h.Routes().ServeHTTP(catW, catReq)
	if catW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", catW.Code, catW.Body.String())
	}
	connReq := httptest.NewRequest(http.MethodGet, "/v1/tenants-connectors/tenant-5", nil)
	connW := httptest.NewRecorder()
	h.Routes().ServeHTTP(connW, connReq)
	if connW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", connW.Code, connW.Body.String())
	}
}

func TestBootstrapTenantBlockedWhenCatalogOrConnectorUnsupported(t *testing.T) {
	store := memory.NewTenantStore()
	events := memory.NewEventLog()
	h := tenantservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{
		"tenant_id":                  "tenant-6",
		"tenant_name":                "Tenant Seis",
		"admin_subject_id":           "admin-6",
		"initial_catalog_refs":       []string{"tenant.unknown.capability"},
		"initial_connector_refs":     []string{"connector.unknown"},
		"policy_profile_ref":         "policy.default",
		"approval_profile_ref":       "approval.default",
		"classification_profile_ref": "classification.default",
		"context_seed":               map[string]any{"company": "ACME"},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bootstrap", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d body=%s", w.Code, w.Body.String())
	}
}
