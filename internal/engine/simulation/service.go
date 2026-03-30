package simulation

import (
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

func (s *Service) RunAll(candidate preview.Candidate) ([]preview.Result, error) {
	now := time.Now().UTC()
	policyDecision, err := s.Policy.Evaluate(policy.Input{
		TenantID:              candidate.TenantID,
		ContractID:            candidate.ContractID,
		ExecutionID:           candidate.ExecutionID,
		ResourceKind:          "execution_record",
		Action:                "execute",
		ClassificationLevel:   "internal",
		ApprovalModeEffective: "pre_execution",
		RiskLevel:             "medium",
	})
	if err != nil {
		return nil, err
	}

	policyStatus := preview.StatusPreviewOK
	policyReason := []string{"preview.ready"}
	if policyDecision.Decision == policy.DecisionRequireApproval {
		policyStatus = preview.StatusPreviewWarning
		policyReason = []string{"preview.warning.policy_sensitive"}
	}

	results := []preview.Result{
		{
			SimulationResultID: fmt.Sprintf("sim-%d", now.UnixNano()),
			PreviewCandidateID: candidate.PreviewCandidateID,
			Family:             preview.SimulationPolicy,
			Status:             policyStatus,
			ReasonCodes:        policyReason,
			InputsRefs:         []string{candidate.PreviewCandidateID, candidate.ContractID},
			OutputsSummary:     string(policyDecision.Decision),
			ConfidenceLevel:    "high",
			CreatedAt:          now,
		},
		{
			SimulationResultID: fmt.Sprintf("sim-%d", now.UnixNano()+1),
			PreviewCandidateID: candidate.PreviewCandidateID,
			Family:             preview.SimulationApproval,
			Status:             preview.StatusPreviewWarning,
			ReasonCodes:        []string{"preview.warning.policy_sensitive"},
			InputsRefs:         []string{candidate.PreviewCandidateID},
			OutputsSummary:     "pre_execution approval likely required",
			ConfidenceLevel:    "medium",
			CreatedAt:          now,
		},
		{
			SimulationResultID: fmt.Sprintf("sim-%d", now.UnixNano()+2),
			PreviewCandidateID: candidate.PreviewCandidateID,
			Family:             preview.SimulationClassification,
			Status:             preview.StatusPreviewOK,
			ReasonCodes:        []string{"preview.ready"},
			InputsRefs:         []string{candidate.PreviewCandidateID},
			OutputsSummary:     "classification preview: internal, no redaction expected",
			ConfidenceLevel:    "medium",
			CreatedAt:          now,
		},
		{
			SimulationResultID: fmt.Sprintf("sim-%d", now.UnixNano()+3),
			PreviewCandidateID: candidate.PreviewCandidateID,
			Family:             preview.SimulationRisk,
			Status:             preview.StatusPreviewWarning,
			ReasonCodes:        []string{"preview.warning.risk_high"},
			InputsRefs:         []string{candidate.PreviewCandidateID},
			OutputsSummary:     "risk preview: medium impact expected",
			ConfidenceLevel:    "medium",
			CreatedAt:          now,
		},
	}

	return results, nil
}
