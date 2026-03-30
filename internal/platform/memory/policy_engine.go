package memory

import (
	"fmt"
	"time"

	"opita-sync-framework/internal/engine/policy"
)

type PolicyEngine struct{}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

func (e *PolicyEngine) Evaluate(input policy.Input) (policy.DecisionRecord, error) {
	decision := policy.DecisionAllow
	reasonCodes := []string{"policy.allow.default"}

	if input.Action == "execute" && input.ApprovalModeEffective != "" && input.ApprovalModeEffective != "auto" {
		decision = policy.DecisionRequireApproval
		reasonCodes = []string{"policy.require_approval.non_auto_mode"}
	}

	if input.ClassificationLevel == "restricted" && input.Action == "view_restricted" {
		decision = policy.DecisionRestrictedView
		reasonCodes = []string{"policy.restricted_view.classification_guard"}
	}

	if input.ResourceKind == "execution_record" && input.RiskLevel == "critical" && input.Action == "execute" {
		decision = policy.DecisionRequireEscalation
		reasonCodes = []string{"policy.require_escalation.critical_risk"}
	}

	return policy.DecisionRecord{
		PolicyDecisionID: fmt.Sprintf("policy-%d", time.Now().UTC().UnixNano()),
		Decision:         decision,
		ReasonCodes:      reasonCodes,
		PolicyVersion:    "cerbos-policy-v1",
		TraceID:          input.ExecutionID,
	}, nil
}
