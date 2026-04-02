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
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
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
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
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

func TestGrantApprovalAndRevocationFlow(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "principal_ref": "user://alice", "principal_type": "person", "capability_id": "tenant.approval.release_blocked_execution", "allowed_actions": []string{"approve"}, "trace_ref": "trace-grant-sensitive-1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal created grant: %v", err)
	}
	grantID := created["grant_id"].(string)
	approveBody, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1", "decision_reason_codes": []string{"grant.approval.release"}})
	approveReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants/"+grantID+"/approve", bytes.NewReader(approveBody))
	approveW := httptest.NewRecorder()
	h.Routes().ServeHTTP(approveW, approveReq)
	if approveW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", approveW.Code, approveW.Body.String())
	}
	revokeBody, _ := json.Marshal(map[string]any{"revoked_by": "admin-1", "justification": "cleanup"})
	revokeReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants/"+grantID+"/revoke", bytes.NewReader(revokeBody))
	revokeW := httptest.NewRecorder()
	h.Routes().ServeHTTP(revokeW, revokeReq)
	if revokeW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", revokeW.Code, revokeW.Body.String())
	}
}

func TestDelegationApprovalAndRevocationFlow(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "source_principal": "role://tenant-admin", "target_principal": "user://bob", "authority_source": "tenant-admin", "scope_type": "capability", "scope_ref": "tenant.recovery.request_manual_compensation", "allowed_actions": []string{"use"}, "can_redelegate": true, "max_depth": 2, "trace_ref": "trace-delegation-sensitive-1"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal created delegation: %v", err)
	}
	grantID := created["grant_id"].(string)
	approveBody, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1", "decision_reason_codes": []string{"delegation.approval.release"}})
	approveReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations/"+grantID+"/approve", bytes.NewReader(approveBody))
	approveW := httptest.NewRecorder()
	h.Routes().ServeHTTP(approveW, approveReq)
	if approveW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", approveW.Code, approveW.Body.String())
	}
	revokeBody, _ := json.Marshal(map[string]any{"revoked_by": "admin-1", "justification": "expired"})
	revokeReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations/"+grantID+"/revoke", bytes.NewReader(revokeBody))
	revokeW := httptest.NewRecorder()
	h.Routes().ServeHTTP(revokeW, revokeReq)
	if revokeW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", revokeW.Code, revokeW.Body.String())
	}
}

func TestAccessWorkspaceReturnsUsableSummary(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	grantBody, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "principal_ref": "user://alice", "principal_type": "person", "capability_id": "tenant.approval.release_blocked_execution", "allowed_actions": []string{"approve"}, "trace_ref": "trace-grant-workspace"})
	grantReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants", bytes.NewReader(grantBody))
	grantW := httptest.NewRecorder()
	h.Routes().ServeHTTP(grantW, grantReq)
	if grantW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", grantW.Code, grantW.Body.String())
	}
	delBody, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "source_principal": "role://tenant-admin", "target_principal": "user://bob", "authority_source": "tenant-admin", "scope_type": "capability", "scope_ref": "tenant.recovery.request_manual_compensation", "allowed_actions": []string{"use"}, "can_redelegate": true, "max_depth": 2, "trace_ref": "trace-delegation-workspace"})
	delReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations", bytes.NewReader(delBody))
	delW := httptest.NewRecorder()
	h.Routes().ServeHTTP(delW, delReq)
	if delW.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", delW.Code, delW.Body.String())
	}
	wsReq := httptest.NewRequest(http.MethodGet, "/v1/tenant-access/workspace/tenant-1", nil)
	wsW := httptest.NewRecorder()
	h.Routes().ServeHTTP(wsW, wsReq)
	if wsW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", wsW.Code, wsW.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(wsW.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal workspace: %v", err)
	}
	if resp["boundary"] != "access_admin_surface_reads_grants_and_guides_governed_authority_changes" {
		t.Fatalf("unexpected boundary: %#v", resp["boundary"])
	}
	grants := resp["grants"].(map[string]any)
	if grants["blocked_grants"].(float64) < 1 {
		t.Fatalf("expected blocked grant summary, got %#v", grants)
	}
}

func TestApproveGrantFailsWhenAlreadyActive(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "principal_ref": "user://alice", "principal_type": "person", "capability_id": "tenant.execution.inspect_run", "allowed_actions": []string{"use"}, "trace_ref": "trace-grant-active"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	grantID := created["grant_id"].(string)
	approveBody, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1"})
	approveReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants/"+grantID+"/approve", bytes.NewReader(approveBody))
	approveW := httptest.NewRecorder()
	h.Routes().ServeHTTP(approveW, approveReq)
	if approveW.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", approveW.Code, approveW.Body.String())
	}
}

func TestRevokeGrantFailsWhenAlreadyRevoked(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "principal_ref": "user://alice", "principal_type": "person", "capability_id": "tenant.execution.inspect_run", "allowed_actions": []string{"use"}, "trace_ref": "trace-grant-revoke"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	grantID := created["grant_id"].(string)
	revokeBody, _ := json.Marshal(map[string]any{"revoked_by": "admin-1"})
	revokeReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants/"+grantID+"/revoke", bytes.NewReader(revokeBody))
	revokeW := httptest.NewRecorder()
	h.Routes().ServeHTTP(revokeW, revokeReq)
	if revokeW.Code != http.StatusOK {
		t.Fatalf("expected first revoke 200, got %d body=%s", revokeW.Code, revokeW.Body.String())
	}
	revokeAgainReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/grants/"+grantID+"/revoke", bytes.NewReader(revokeBody))
	revokeAgainW := httptest.NewRecorder()
	h.Routes().ServeHTTP(revokeAgainW, revokeAgainReq)
	if revokeAgainW.Code != http.StatusConflict {
		t.Fatalf("expected second revoke 409, got %d body=%s", revokeAgainW.Code, revokeAgainW.Body.String())
	}
}

func TestApproveDelegationFailsWhenAlreadyActive(t *testing.T) {
	store := memory.NewAccessStore()
	events := memory.NewEventLog()
	approvals := memory.NewApprovalStore()
	h := accessservice.NewHandler(store, events, approvals)
	body, _ := json.Marshal(map[string]any{"tenant_id": "tenant-1", "source_principal": "role://tenant-admin", "target_principal": "user://bob", "authority_source": "tenant-admin", "scope_type": "capability", "scope_ref": "tenant.execution.inspect_run", "allowed_actions": []string{"use"}, "can_redelegate": false, "max_depth": 1, "trace_ref": "trace-delegation-active"})
	req := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, req)
	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	grantID := created["grant_id"].(string)
	approveBody, _ := json.Marshal(map[string]any{"decided_by_subject_id": "approver-1"})
	approveReq := httptest.NewRequest(http.MethodPost, "/v1/tenant-access/delegations/"+grantID+"/approve", bytes.NewReader(approveBody))
	approveW := httptest.NewRecorder()
	h.Routes().ServeHTTP(approveW, approveReq)
	if approveW.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", approveW.Code, approveW.Body.String())
	}
}
