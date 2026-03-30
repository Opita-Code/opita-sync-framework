package memory

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"opita-sync-framework/internal/engine/approvals"
)

type ApprovalStore struct {
	mu       sync.RWMutex
	requests map[string]approvals.Request
}

func NewApprovalStore() *ApprovalStore {
	return &ApprovalStore{requests: map[string]approvals.Request{}}
}

func (s *ApprovalStore) Create(request approvals.Request) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.requests[request.ApprovalRequestID]; exists {
		return errors.New("approval request already exists")
	}
	s.requests[request.ApprovalRequestID] = request
	return nil
}

func (s *ApprovalStore) GetByID(approvalRequestID string) (approvals.Request, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req, found := s.requests[approvalRequestID]
	return req, found, nil
}

func (s *ApprovalStore) Decide(approvalRequestID string, decision approvals.Decision) (approvals.Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	req, found := s.requests[approvalRequestID]
	if !found {
		return approvals.Request{}, errors.New("approval request not found")
	}
	if decision.DecidedBySubjectID == "" {
		return approvals.Request{}, errors.New("decided_by_subject_id is required")
	}
	switch decision.State {
	case approvals.StateReleased, approvals.StateRejected, approvals.StateEscalated:
	default:
		return approvals.Request{}, fmt.Errorf("invalid approval decision state: %s", decision.State)
	}
	decidedAt := decision.DecidedAt
	if decidedAt.IsZero() {
		decidedAt = time.Now().UTC()
	}
	req.State = decision.State
	req.DecidedBySubjectID = decision.DecidedBySubjectID
	req.DecisionComment = decision.DecisionComment
	req.DecisionReasonCodes = append([]string(nil), decision.DecisionReasonCodes...)
	req.UpdatedAt = decidedAt
	req.ReleasedAt = nil
	req.RejectedAt = nil
	req.EscalatedAt = nil
	switch decision.State {
	case approvals.StateReleased:
		req.ReleasedAt = &decidedAt
	case approvals.StateRejected:
		req.RejectedAt = &decidedAt
	case approvals.StateEscalated:
		req.EscalatedAt = &decidedAt
	}
	s.requests[approvalRequestID] = req
	return req, nil
}
