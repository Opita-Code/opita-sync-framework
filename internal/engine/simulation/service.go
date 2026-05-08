package simulation

import (
	"context"
	"fmt"
	"time"

	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/preview"
)

type Service struct {
	Policy policy.PolicyEngine
}

func NewService(policyEngine policy.PolicyEngine) *Service {
	return &Service{Policy: policyEngine}
}

func (s *Service) RunAll(ctx context.Context, candidate preview.Candidate) ([]preview.Result, error) {
	if candidate.TenantID == "" {
		return nil, fmt.Errorf("simulation: tenant_id is required")
	}
	if candidate.ContractID == "" {
		return nil, fmt.Errorf("simulation: contract_id is required")
	}

	now := time.Now().UTC()

	// Define family-specific evaluation variants.
	// Each family calls Policy.Evaluate with distinct Action/RiskLevel to
	// produce realistic, differentiated results instead of hardcoded values.
	type familyEval struct {
		Family     preview.SimulationFamily
		Action     string
		RiskLevel  string
		OutputFmt  string
		Confidence string
	}

	families := []familyEval{
		{preview.SimulationPolicy, "execute", "medium", "policy decision: %s", "high"},
		{preview.SimulationApproval, "approve", "medium", "approval preview: pre_execution approval %s", "medium"},
		{preview.SimulationClassification, "classify", "low", "classification preview: internal, %s", "medium"},
		{preview.SimulationRisk, "assess_risk", "medium", "risk preview: %s impact expected", "medium"},
	}

	results := make([]preview.Result, 0, len(families))
	for i, fam := range families {
		input := policy.Input{
			TenantID:              candidate.TenantID,
			ContractID:            candidate.ContractID,
			ExecutionID:           candidate.ExecutionID,
			ResourceKind:          "execution_record",
			Action:                fam.Action,
			ClassificationLevel:   "internal",
			ApprovalModeEffective: "pre_execution",
			RiskLevel:             fam.RiskLevel,
		}

		decision, err := s.Policy.Evaluate(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("simulation %s evaluation: %w", fam.Family, err)
		}

		status, reasonCodes := decisionToStatus(decision.Decision)
		results = append(results, preview.Result{
			SimulationResultID: fmt.Sprintf("sim-%d", now.UnixNano()+int64(i)),
			PreviewCandidateID: candidate.PreviewCandidateID,
			Family:             fam.Family,
			Status:             status,
			ReasonCodes:        reasonCodes,
			InputsRefs:         []string{candidate.PreviewCandidateID, candidate.ContractID},
			OutputsSummary:     fmt.Sprintf(fam.OutputFmt, decision.Decision),
			ConfidenceLevel:    fam.Confidence,
			CreatedAt:          now,
		})
	}

	return results, nil
}

// decisionToStatus maps a policy decision to a preview status and reason codes.
func decisionToStatus(d policy.Decision) (preview.Status, []string) {
	switch d {
	case policy.DecisionAllow:
		return preview.StatusPreviewOK, []string{"preview.ready"}
	case policy.DecisionRequireApproval:
		return preview.StatusPreviewWarning, []string{"preview.warning.policy_sensitive"}
	case policy.DecisionDenyBlock, policy.DecisionRequireEscalation:
		return preview.StatusPreviewBlocked, []string{"preview.warning.policy_blocked"}
	case policy.DecisionRestrictedView:
		return preview.StatusPreviewWarning, []string{"preview.warning.restricted_view"}
	default:
		return preview.StatusPreviewWarning, []string{"preview.warning.unknown_decision"}
	}
}
