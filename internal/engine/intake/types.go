package intake

import "time"

type SessionState string

const (
	SessionStateFreeChat              SessionState = "free_chat"
	SessionStateShaping               SessionState = "shaping"
	SessionStateAwaitingClarification SessionState = "awaiting_clarification"
	SessionStateIntentReady           SessionState = "intent_ready"
	SessionStateProposalReady         SessionState = "proposal_ready"
	SessionStateOutOfScope            SessionState = "out_of_scope"
	SessionStateClosed                SessionState = "closed"
)

type Decision string

const (
	DecisionContinueFreeChat  Decision = "continue_free_chat"
	DecisionAskClarification  Decision = "ask_clarification"
	DecisionEmitIntentInput   Decision = "emit_intent_input"
	DecisionEmitProposalDraft Decision = "emit_proposal_draft"
	DecisionStopOutOfScope    Decision = "stop_out_of_scope"
)

type ConversationTurn struct {
	ConversationTurnID string
	SessionID          string
	TenantID           string
	SubjectID          string
	MessageRole        string
	RawText            string
	AttachmentsRefs    []string
	Timestamp          time.Time
	TurnClassification string
	TraceID            string
}

type Session struct {
	IntakeSessionID string
	SessionID       string
	TenantID        string
	SubjectID       string
	CurrentState    SessionState
	OpenQuestions   []string
	ResolvedFacts   []string
	AmbiguityLevel  string
	LastDecision    Decision
	TraceID         string
	UpdatedAt       time.Time
}

type IntentCandidate struct {
	IntentCandidateID     string
	SourceTurnIDs         []string
	ObjetivoCandidate     string
	AlcanceCandidate      string
	ArtifactsCandidate    []string
	ConstraintsCandidate  []string
	Assumptions           []string
	OpenQuestions         []string
	ConfidenceLevel       string
	ReadyForIntentInput   bool
	ReadyForProposalDraft bool
}

type Service interface {
	CreateTurn(turn ConversationTurn) error
	CreateSession(session Session) error
	SaveIntentCandidate(candidate IntentCandidate) error
	GetSession(intakeSessionID string) (Session, bool, error)
	GetIntentCandidate(intentCandidateID string) (IntentCandidate, bool, error)
}
