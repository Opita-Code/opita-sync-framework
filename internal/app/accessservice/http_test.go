package accessservice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opita-sync-framework/internal/app/accessservice"
	"opita-sync-framework/internal/platform/memory"
)

func TestCreateGrantAndListByTenant(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	h := accessservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "principal_ref": "user://alice", "principal_type": "person", "capability_id": "tenant.execution.inspect_run", "allowed_actions": []string{"use"}, "trace_ref": "trace-grant-1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	listReq := httptest.NewRequest(http.MethodGet, "/v1/tenant-access/grants/tenant-1", nil)
	listW := httptest.NewRecorder()
	h.Routes().ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", listW.Code, listW.Body.String())
	}
}

func TestCreateDelegationAndListByTenant(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	h := accessservice.NewHandler(store, events)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "source_principal": "role://tenant-admin", "target_principal": "user://bob", "authority_source": "tenant-admin", "scope_type": "capability", "scope_ref": "tenant.approval.release_blocked_execution", "allowed_actions": []string{"approve"}, "can_redelegate": true, "max_depth": 2, "trace_ref": "trace-delegation-1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	listReq := httptest.NewRequest(http.MethodGet, "/v1/tenant-access/delegations/tenant-1", nil)
	listW := httptest.NewRecorder()
	h.Routes().ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", listW.Code, listW.Body.String())
	}
}
