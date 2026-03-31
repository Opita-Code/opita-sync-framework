package foundation

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/registry"
	"opita-sync-framework/internal/engine/runtime"
)

type FoundationOrchestrator struct {
	Compiler  intentCompiler
	Policy    policy.PolicyEngine
	Runtime   runtime.RuntimeService
	Events    events.EventLog
	Registry  registry.Resolver
	Runs      RunRepository
	Approvals approvals.Service
}

var uniqueIDCounter atomic.Uint64

type intentCompiler interface {
	Compile(ctx context.Context, input intent.IntentInput) (intent.CompiledContract, intent.CompilationReport, error)
}

type FoundationRunResult struct {
	Contract       intent.CompiledContract
	Report         intent.CompilationReport
	Execution      runtime.ExecutionRecord
	PolicyDecision policy.DecisionRecord
	Resolution     registry.ResolutionResult
	Approval       *approvals.Request
}

func (o *FoundationOrchestrator) Run(ctx context.Context, input intent.IntentInput) (FoundationRunResult, error) {
	contract, report, err := o.Compiler.Compile(ctx, input)
	if err != nil {
		return FoundationRunResult{}, err
	}

	resolution, err := o.Registry.Resolve(registry.ResolutionRequest{
		CapabilityID:          capabilityIDForResultType(contract.TipoResultadoEsperado),
		ContractSchemaVersion: contract.ContractVersion,
		SupportedResultType:   string(contract.TipoResultadoEsperado),
		Environment:           "dev",
	})
	if err != nil {
		return FoundationRunResult{}, fmt.Errorf("resolve capability: %w", err)
	}

	executionID := uniqueID("exec")
	traceID := contract.TraceID
	if traceID == "" {
		traceID = fmt.Sprintf("trace-%s", executionID)
	}

	decision, err := o.Policy.Evaluate(policy.Input{
		TenantID:              contract.TenantID,
		ContractID:            contract.ContractID,
		ExecutionID:           executionID,
		ResourceKind:          "execution_record",
		Action:                "execute",
		ClassificationLevel:   "internal",
		ApprovalModeEffective: deriveApprovalMode(contract.TipoResultadoEsperado),
		RiskLevel:             deriveRiskLevel(contract.TipoResultadoEsperado),
	})
	if err != nil {
		return FoundationRunResult{}, fmt.Errorf("evaluate policy: %w", err)
	}

	state := runtime.ExecutionStateExecutionReleased
	var approvalRecord *approvals.Request
	switch decision.Decision {
	case policy.DecisionRequireApproval:
		state = runtime.ExecutionStateAwaitingApproval
		record := approvals.Request{
			ApprovalRequestID:         uniqueID("approval"),
			ExecutionID:               executionID,
			ContractID:                contract.ContractID,
			TenantID:                  contract.TenantID,
			TraceID:                   traceID,
			State:                     approvals.StateAwaitingApproval,
			Mode:                      deriveApprovalMode(contract.TipoResultadoEsperado),
			ReasonCodes:               decision.ReasonCodes,
			SourceContractFingerprint: contract.Fingerprint,
			CreatedAt:                 time.Now().UTC(),
			UpdatedAt:                 time.Now().UTC(),
		}
		if err := o.Approvals.Create(record); err != nil {
			return FoundationRunResult{}, fmt.Errorf("create approval request: %w", err)
		}
		approvalRecord = &record
	case policy.DecisionRequireEscalation, policy.DecisionDenyBlock:
		state = runtime.ExecutionStateBlocked
	}

	execution := runtime.ExecutionRecord{
		ExecutionID:         executionID,
		TenantID:            contract.TenantID,
		ContractID:          contract.ContractID,
		ContractFingerprint: contract.Fingerprint,
		TraceID:             traceID,
		State:               state,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}

	if err := o.Runtime.CreateExecution(execution); err != nil {
		return FoundationRunResult{}, fmt.Errorf("create execution: %w", err)
	}

	eventsToAppend := []events.Record{
		{
			EventID:             uniqueID("event"),
			EventType:           "contract.compilation_completed",
			TenantID:            contract.TenantID,
			TraceID:             traceID,
			ContractID:          contract.ContractID,
			ContractFingerprint: contract.Fingerprint,
			ExecutionID:         executionID,
			PolicyDecisionID:    decision.PolicyDecisionID,
			ConversationTurnID:  contract.ConversationTurnID,
			IntakeSessionID:     contract.IntakeSessionID,
			IntentCandidateID:   contract.IntentCandidateID,
			ProposalDraftID:     contract.ProposalDraftID,
			PatchsetCandidateID: contract.PatchsetCandidateID,
			PreviewCandidateID:  contract.PreviewCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload:             map[string]any{"status": report.Status, "simulation_result_ids": contract.SimulationResultIDs},
		},
		{
			EventID:             uniqueID("event"),
			EventType:           "policy.decision_recorded",
			TenantID:            contract.TenantID,
			TraceID:             traceID,
			ContractID:          contract.ContractID,
			ContractFingerprint: contract.Fingerprint,
			ExecutionID:         executionID,
			PolicyDecisionID:    decision.PolicyDecisionID,
			ConversationTurnID:  contract.ConversationTurnID,
			IntakeSessionID:     contract.IntakeSessionID,
			IntentCandidateID:   contract.IntentCandidateID,
			ProposalDraftID:     contract.ProposalDraftID,
			PatchsetCandidateID: contract.PatchsetCandidateID,
			PreviewCandidateID:  contract.PreviewCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload:             map[string]any{"decision": decision.Decision, "reason_codes": decision.ReasonCodes, "simulation_result_ids": contract.SimulationResultIDs},
		},
		{
			EventID:             uniqueID("event"),
			EventType:           "execution.created",
			TenantID:            contract.TenantID,
			TraceID:             traceID,
			ContractID:          contract.ContractID,
			ContractFingerprint: contract.Fingerprint,
			ExecutionID:         executionID,
			PolicyDecisionID:    decision.PolicyDecisionID,
			ConversationTurnID:  contract.ConversationTurnID,
			IntakeSessionID:     contract.IntakeSessionID,
			IntentCandidateID:   contract.IntentCandidateID,
			ProposalDraftID:     contract.ProposalDraftID,
			PatchsetCandidateID: contract.PatchsetCandidateID,
			PreviewCandidateID:  contract.PreviewCandidateID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"runtime_state":         state,
				"capability_ref":        resolution.CapabilityManifestRef,
				"binding_id":            resolution.BindingID,
				"provider_ref":          resolution.ProviderRef,
				"simulation_result_ids": contract.SimulationResultIDs,
			},
		},
	}

	if approvalRecord != nil {
		eventsToAppend = append(eventsToAppend, events.Record{
			EventID:             uniqueID("event"),
			EventType:           "approval.awaiting",
			TenantID:            contract.TenantID,
			TraceID:             traceID,
			ContractID:          contract.ContractID,
			ContractFingerprint: contract.Fingerprint,
			ExecutionID:         executionID,
			PolicyDecisionID:    decision.PolicyDecisionID,
			ConversationTurnID:  contract.ConversationTurnID,
			IntakeSessionID:     contract.IntakeSessionID,
			IntentCandidateID:   contract.IntentCandidateID,
			ProposalDraftID:     contract.ProposalDraftID,
			PatchsetCandidateID: contract.PatchsetCandidateID,
			PreviewCandidateID:  contract.PreviewCandidateID,
			ApprovalRequestID:   approvalRecord.ApprovalRequestID,
			OccurredAt:          time.Now().UTC(),
			Payload: map[string]any{
				"approval_request_id":   approvalRecord.ApprovalRequestID,
				"mode":                  approvalRecord.Mode,
				"state":                 approvalRecord.State,
				"simulation_result_ids": contract.SimulationResultIDs,
			},
		})
	}

	for _, record := range eventsToAppend {
		if err := o.Events.Append(record); err != nil {
			return FoundationRunResult{}, fmt.Errorf("append event: %w", err)
		}
	}

	result := FoundationRunResult{
		Contract:       contract,
		Report:         report,
		Execution:      execution,
		PolicyDecision: decision,
		Resolution:     resolution,
		Approval:       approvalRecord,
	}

	if err := o.Runs.Save(result); err != nil {
		return FoundationRunResult{}, fmt.Errorf("store foundation run: %w", err)
	}

	return result, nil
}

func capabilityIDForResultType(resultType intent.ResultType) string {
	if resultType == "" {
		return "capability.plan.default"
	}
	return fmt.Sprintf("capability.%s.default", resultType)
}

func deriveApprovalMode(resultType intent.ResultType) string {
	switch resultType {
	case intent.ResultTypeExecution, intent.ResultTypeSystemUpdate, intent.ResultTypeChangeProposal:
		return "pre_execution"
	default:
		return "auto"
	}
}

func deriveRiskLevel(resultType intent.ResultType) string {
	switch resultType {
	case intent.ResultTypeSystemUpdate:
		return "high"
	case intent.ResultTypeExecution, intent.ResultTypeChangeProposal:
		return "medium"
	default:
		return "low"
	}
}

func uniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UTC().UnixNano(), uniqueIDCounter.Add(1))
}
