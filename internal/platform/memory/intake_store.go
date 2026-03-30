package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/intake"
)

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

func (s *IntakeStore) CreateTurn(turn intake.ConversationTurn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.turns[turn.ConversationTurnID]; exists {
		return errors.New("conversation turn already exists")
	}
	s.turns[turn.ConversationTurnID] = turn
	return nil
}

func (s *IntakeStore) CreateSession(session intake.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.IntakeSessionID] = session
	return nil
}

func (s *IntakeStore) SaveIntentCandidate(candidate intake.IntentCandidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.candidates[candidate.IntentCandidateID] = candidate
	return nil
}

func (s *IntakeStore) GetSession(intakeSessionID string) (intake.Session, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, found := s.sessions[intakeSessionID]
	return session, found, nil
}

func (s *IntakeStore) GetIntentCandidate(intentCandidateID string) (intake.IntentCandidate, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	candidate, found := s.candidates[intentCandidateID]
	return candidate, found, nil
}
