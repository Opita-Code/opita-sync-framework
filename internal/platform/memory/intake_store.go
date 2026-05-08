package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/intake"
)

var _ intake.Service = (*IntakeStore)(nil)

type IntakeStore struct {
	mu         sync.RWMutex
	turns      map[string]intake.ConversationTurn
	sessions   map[string]intake.Session
	candidates map[string]intake.IntentCandidate
}

func NewIntakeStore() *IntakeStore {
	return &IntakeStore{
		turns:      map[string]intake.ConversationTurn{},
		sessions:   map[string]intake.Session{},
		candidates: map[string]intake.IntentCandidate{},
	}
}

func (s *IntakeStore) CreateTurn(ctx context.Context, turn intake.ConversationTurn) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.turns[turn.ConversationTurnID]; exists {
		return errors.New("conversation turn already exists")
	}
	s.turns[turn.ConversationTurnID] = turn
	return nil
}

func (s *IntakeStore) CreateSession(ctx context.Context, session intake.Session) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.IntakeSessionID] = session
	return nil
}

func (s *IntakeStore) SaveIntentCandidate(ctx context.Context, candidate intake.IntentCandidate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.candidates[candidate.IntentCandidateID] = candidate
	return nil
}

func (s *IntakeStore) GetSession(ctx context.Context, intakeSessionID string) (intake.Session, bool, error) {
	select {
	case <-ctx.Done():
		return intake.Session{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, found := s.sessions[intakeSessionID]
	return session, found, nil
}

func (s *IntakeStore) GetIntentCandidate(ctx context.Context, intentCandidateID string) (intake.IntentCandidate, bool, error) {
	select {
	case <-ctx.Done():
		return intake.IntentCandidate{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[intentCandidateID]
	return candidate, found, nil
}
