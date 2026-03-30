package intent

import "time"

type ResultType string

const (
	ResultTypePlan               ResultType = "plan"
	ResultTypeInspection         ResultType = "inspection"
	ResultTypeQuery              ResultType = "query"
	ResultTypeReport             ResultType = "report"
	ResultTypeChangeProposal     ResultType = "change_proposal"
	ResultTypeExecution          ResultType = "execution"
	ResultTypeSystemUpdate       ResultType = "system_update"
	ResultTypeGovernanceDecision ResultType = "governance_decision"
)

type Autonomy string

const (
	AutonomyManual     Autonomy = "manual"
	AutonomyAssisted   Autonomy = "assisted"
	AutonomyAutonomous Autonomy = "autonomous"
)

type ContractState string

const (
	ContractStateDraft    ContractState = "draft"
	ContractStateCompiled ContractState = "compiled"
)

type IntentInput struct {
	RequestID             string
	TenantID              string
	WorkspaceID           string
	UserID                string
	SessionID             string
	TraceID               string
	ConversationTurnID    string
	IntakeSessionID       string
	IntentCandidateID     string
	ProposalDraftID       string
	PatchsetCandidateID   string
	PreviewCandidateID    string
	SimulationResultIDs   []string
	Objetivo              string
	Alcance               string
	TipoResultadoEsperado ResultType
	Restricciones         []string
	SistemasPosibles      []string
	DatosPermitidos       []string
	AutonomiaSolicitada   Autonomy
	AprobacionRequerida   string
	CriteriosDeExito      []string
	CreatedAt             time.Time
}

type CanonicalIntent struct {
	RequestID             string
	TenantID              string
	WorkspaceID           string
	UserID                string
	SessionID             string
	TraceID               string
	ConversationTurnID    string
	IntakeSessionID       string
	IntentCandidateID     string
	ProposalDraftID       string
	PatchsetCandidateID   string
	PreviewCandidateID    string
	SimulationResultIDs   []string
	Objetivo              string
	Alcance               string
	TipoResultadoEsperado ResultType
	Restricciones         []string
	SistemasPosibles      []string
	DatosPermitidos       []string
	AutonomiaSolicitada   Autonomy
	AprobacionRequerida   string
	CriteriosDeExito      []string
	CreatedAt             time.Time
}

type SnapshotBundle struct {
	PolicyVersion         string
	ClassificationVersion string
	RiskVersion           string
	PermissionVersion     string
	CompiledAt            time.Time
}

type CompiledContract struct {
	ContractID            string
	ContractVersion       string
	RequestID             string
	TenantID              string
	WorkspaceID           string
	UserID                string
	SessionID             string
	TraceID               string
	ConversationTurnID    string
	IntakeSessionID       string
	IntentCandidateID     string
	ProposalDraftID       string
	PatchsetCandidateID   string
	PreviewCandidateID    string
	SimulationResultIDs   []string
	Objetivo              string
	Alcance               string
	TipoResultadoEsperado ResultType
	Restricciones         []string
	SistemasPosibles      []string
	DatosPermitidos       []string
	AutonomiaSolicitada   Autonomy
	AprobacionRequerida   string
	CriteriosDeExito      []string
	Fingerprint           string
	State                 ContractState
	Snapshots             SnapshotBundle
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type DiagnosticSeverity string

const (
	SeverityError   DiagnosticSeverity = "error"
	SeverityWarning DiagnosticSeverity = "warning"
	SeverityInfo    DiagnosticSeverity = "info"
)

type Diagnostic struct {
	Code     string
	Severity DiagnosticSeverity
	Field    string
	Message  string
}

type CompilationStatus string

const (
	CompilationStatusCompiled  CompilationStatus = "compiled"
	CompilationStatusRejected  CompilationStatus = "rejected"
	CompilationStatusDuplicate CompilationStatus = "duplicate"
)

type CompilationReport struct {
	ContractID   string
	Fingerprint  string
	Status       CompilationStatus
	Diagnostics  []Diagnostic
	DuplicatedOf string
	CompiledAt   time.Time
}
