package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/approvals"
)

type ApprovalStore struct {
	store *Store
}

func NewApprovalStore(store *Store) *ApprovalStore {
	return &ApprovalStore{store: store}
}

func (s *ApprovalStore) Create(request approvals.Request) error {
	raw, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal approval request: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into approval_requests (approval_request_id, execution_id, contract_id, tenant_id, trace_id, state, mode, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, request.ApprovalRequestID, request.ExecutionID, request.ContractID, request.TenantID, request.TraceID, request.State, request.Mode, raw, request.CreatedAt, request.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert approval request: %w", err)
	}
	return nil
}

func (s *ApprovalStore) GetByID(approvalRequestID string) (approvals.Request, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from approval_requests where approval_request_id = $1`, approvalRequestID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return approvals.Request{}, false, nil
		}
		return approvals.Request{}, false, fmt.Errorf("select approval request: %w", err)
	}
	var request approvals.Request
	if err := json.Unmarshal(raw, &request); err != nil {
		return approvals.Request{}, false, fmt.Errorf("unmarshal approval request: %w", err)
	}
	return request, true, nil
}

func (s *ApprovalStore) Decide(approvalRequestID string, decision approvals.Decision) (approvals.Request, error) {
	request, found, err := s.GetByID(approvalRequestID)
	if err != nil {
		return approvals.Request{}, err
	}
	if !found {
		return approvals.Request{}, fmt.Errorf("approval request not found")
	}
	if decision.DecidedBySubjectID == "" {
		return approvals.Request{}, fmt.Errorf("decided_by_subject_id is required")
	}
	decidedAt := decision.DecidedAt
	if decidedAt.IsZero() {
		decidedAt = nowUTC()
	}
	switch decision.State {
	case approvals.StateReleased, approvals.StateRejected, approvals.StateEscalated:
	default:
		return approvals.Request{}, fmt.Errorf("invalid approval decision state: %s", decision.State)
	}
	request.State = decision.State
	request.DecidedBySubjectID = decision.DecidedBySubjectID
	request.DecisionComment = decision.DecisionComment
	request.DecisionReasonCodes = append([]string(nil), decision.DecisionReasonCodes...)
	request.UpdatedAt = decidedAt
	request.ReleasedAt = nil
	request.RejectedAt = nil
	request.EscalatedAt = nil
	switch decision.State {
	case approvals.StateReleased:
		request.ReleasedAt = &decidedAt
	case approvals.StateRejected:
		request.RejectedAt = &decidedAt
	case approvals.StateEscalated:
		request.EscalatedAt = &decidedAt
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return approvals.Request{}, fmt.Errorf("marshal approval request: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		update approval_requests set state = $2, payload = $3, updated_at = now() where approval_request_id = $1
	`, approvalRequestID, request.State, raw)
	if err != nil {
		return approvals.Request{}, fmt.Errorf("update approval request: %w", err)
	}
	return request, nil
}
