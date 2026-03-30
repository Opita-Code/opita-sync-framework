package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"opita-sync-framework/internal/engine/intake"
)

type IntakeStore struct {
	store *Store
}

func NewIntakeStore(store *Store) *IntakeStore {
	return &IntakeStore{store: store}
}

func (s *IntakeStore) CreateTurn(turn intake.ConversationTurn) error {
	raw, err := json.Marshal(turn)
	if err != nil {
		return fmt.Errorf("marshal intake turn: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into intake_turns (conversation_turn_id, session_id, tenant_id, subject_id, trace_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6, $7)
	`, turn.ConversationTurnID, turn.SessionID, turn.TenantID, turn.SubjectID, turn.TraceID, raw, turn.Timestamp)
	if err != nil {
		return fmt.Errorf("insert intake turn: %w", err)
	}
	return nil
}

func (s *IntakeStore) CreateSession(session intake.Session) error {
	raw, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal intake session: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into intake_sessions (intake_session_id, session_id, tenant_id, subject_id, trace_id, payload, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)
		on conflict (intake_session_id) do update set payload = excluded.payload, updated_at = excluded.updated_at
	`, session.IntakeSessionID, session.SessionID, session.TenantID, session.SubjectID, session.TraceID, raw, session.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert intake session: %w", err)
	}
	return nil
}

func (s *IntakeStore) SaveIntentCandidate(candidate intake.IntentCandidate) error {
	raw, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("marshal intent candidate: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into intent_candidates (intent_candidate_id, payload, created_at)
		values ($1, $2, $3)
		on conflict (intent_candidate_id) do update set payload = excluded.payload
	`, candidate.IntentCandidateID, raw, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("upsert intent candidate: %w", err)
	}
	return nil
}

func (s *IntakeStore) GetSession(intakeSessionID string) (intake.Session, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from intake_sessions where intake_session_id = $1`, intakeSessionID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return intake.Session{}, false, nil
		}
		return intake.Session{}, false, fmt.Errorf("select intake session: %w", err)
	}
	var session intake.Session
	if err := json.Unmarshal(raw, &session); err != nil {
		return intake.Session{}, false, fmt.Errorf("unmarshal intake session: %w", err)
	}
	return session, true, nil
}

func (s *IntakeStore) GetIntentCandidate(intentCandidateID string) (intake.IntentCandidate, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from intent_candidates where intent_candidate_id = $1`, intentCandidateID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return intake.IntentCandidate{}, false, nil
		}
		return intake.IntentCandidate{}, false, fmt.Errorf("select intent candidate: %w", err)
	}
	var candidate intake.IntentCandidate
	if err := json.Unmarshal(raw, &candidate); err != nil {
		return intake.IntentCandidate{}, false, fmt.Errorf("unmarshal intent candidate: %w", err)
	}
	return candidate, true, nil
}

var _ intake.Service = (*IntakeStore)(nil)

func (s *IntakeStore) GetByID(_ context.Context, _ string) (intake.IntentCandidate, bool, error) {
	return intake.IntentCandidate{}, false, fmt.Errorf("not supported")
}
