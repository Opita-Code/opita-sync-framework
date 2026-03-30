package policy

type Decision string

const (
	DecisionAllow             Decision = "allow"
	DecisionDenyBlock         Decision = "deny_block"
	DecisionRequireApproval   Decision = "require_approval"
	DecisionRequireEscalation Decision = "require_escalation"
	DecisionRestrictedView    Decision = "restricted_view"
)

type Input struct {
	TenantID              string
	ContractID            string
	ExecutionID           string
	ResourceKind          string
	Action                string
	ClassificationLevel   string
	ApprovalModeEffective string
	RiskLevel             string
}

type DecisionRecord struct {
	PolicyDecisionID string
	Decision         Decision
	ReasonCodes      []string
	PolicyVersion    string
	TraceID          string
}

type PolicyEngine interface {
	Evaluate(input Input) (DecisionRecord, error)
}
