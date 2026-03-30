package foundation_test

import (
	"context"
	"testing"
	"time"

	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/platform/filesystem"
	"opita-sync-framework/internal/platform/memory"
)

func TestFoundationOrchestratorRunPlanFlow(t *testing.T) {
	resolver, err := filesystem.NewRegistryResolver("../../../definitions/capabilities")
	if err != nil {
		t.Fatalf("expected resolver, got %v", err)
	}
	orchestrator := &foundation.FoundationOrchestrator{
		Compiler:  intent.NewCompiler(memory.NewContractRepository()),
		Policy:    memory.NewPolicyEngine(),
		Runtime:   memory.NewRuntimeService(),
		Events:    memory.NewEventLog(),
		Registry:  resolver,
		Runs:      memory.NewFoundationRunStore(),
		Approvals: memory.NewApprovalStore(),
	}
	if err := orchestrator.Validate(); err != nil {
		t.Fatalf("expected valid orchestrator, got %v", err)
	}
	result, err := orchestrator.Run(context.Background(), intent.IntentInput{
		RequestID:             "req-foundation-plan",
		TenantID:              "tenant-1",
		WorkspaceID:           "workspace-1",
		UserID:                "user-1",
		SessionID:             "session-1",
		Objetivo:              "Compilar y crear ejecución plan",
		Alcance:               "foundation-scope",
		TipoResultadoEsperado: intent.ResultTypePlan,
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Execution.State != "execution_released" {
		t.Fatalf("expected execution_released, got %s", result.Execution.State)
	}
	if result.PolicyDecision.Decision != "allow" {
		t.Fatalf("expected allow decision, got %s", result.PolicyDecision.Decision)
	}
}

func TestFoundationOrchestratorRunExecutionRequiresApproval(t *testing.T) {
	resolver, err := filesystem.NewRegistryResolver("../../../definitions/capabilities")
	if err != nil {
		t.Fatalf("expected resolver, got %v", err)
	}
	orchestrator := &foundation.FoundationOrchestrator{
		Compiler:  intent.NewCompiler(memory.NewContractRepository()),
		Policy:    memory.NewPolicyEngine(),
		Runtime:   memory.NewRuntimeService(),
		Events:    memory.NewEventLog(),
		Registry:  resolver,
		Runs:      memory.NewFoundationRunStore(),
		Approvals: memory.NewApprovalStore(),
	}
	result, err := orchestrator.Run(context.Background(), intent.IntentInput{
		RequestID:             "req-foundation-execution",
		TenantID:              "tenant-1",
		WorkspaceID:           "workspace-1",
		UserID:                "user-1",
		SessionID:             "session-1",
		Objetivo:              "Compilar y crear ejecución governed",
		Alcance:               "foundation-scope",
		TipoResultadoEsperado: intent.ResultTypeExecution,
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Execution.State != "awaiting_approval" {
		t.Fatalf("expected awaiting_approval, got %s", result.Execution.State)
	}
	if result.PolicyDecision.Decision != "require_approval" {
		t.Fatalf("expected require_approval decision, got %s", result.PolicyDecision.Decision)
	}
	if result.Approval == nil {
		t.Fatalf("expected approval request to be created")
	}
	if result.Approval.State != "awaiting_approval" {
		t.Fatalf("expected approval state awaiting_approval, got %s", result.Approval.State)
	}
}
