package intent_test

import (
	"context"
	"testing"
	"time"

	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/platform/memory"
)

func TestCompilerCompileSuccess(t *testing.T) {
	repo := memory.NewContractRepository()
	compiler := intent.NewCompiler(repo)

	contract, report, err := compiler.Compile(context.Background(), intent.IntentInput{
		RequestID:             "req-1",
		TenantID:              "tenant-1",
		WorkspaceID:           "workspace-1",
		UserID:                "user-1",
		SessionID:             "session-1",
		Objetivo:              "Crear un contrato compilado",
		Alcance:               "test-scope",
		TipoResultadoEsperado: intent.ResultTypePlan,
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contract.ContractID == "" {
		t.Fatalf("expected contract id")
	}
	if contract.Fingerprint == "" {
		t.Fatalf("expected fingerprint")
	}
	if report.Status != intent.CompilationStatusCompiled {
		t.Fatalf("expected status compiled, got %s", report.Status)
	}
}

func TestCompilerDetectsDuplicateFingerprint(t *testing.T) {
	repo := memory.NewContractRepository()
	compiler := intent.NewCompiler(repo)
	input := intent.IntentInput{
		RequestID:             "req-duplicate",
		TenantID:              "tenant-1",
		WorkspaceID:           "workspace-1",
		UserID:                "user-1",
		SessionID:             "session-1",
		Objetivo:              "Crear un contrato compilado",
		Alcance:               "test-scope",
		TipoResultadoEsperado: intent.ResultTypePlan,
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	}

	first, _, err := compiler.Compile(context.Background(), input)
	if err != nil {
		t.Fatalf("first compile failed: %v", err)
	}
	second, report, err := compiler.Compile(context.Background(), input)
	if err != nil {
		t.Fatalf("second compile failed: %v", err)
	}
	if report.Status != intent.CompilationStatusDuplicate {
		t.Fatalf("expected duplicate status, got %s", report.Status)
	}
	if first.ContractID != second.ContractID {
		t.Fatalf("expected same contract for duplicate fingerprint")
	}
}

func TestCompilerRejectsInvalidResultType(t *testing.T) {
	repo := memory.NewContractRepository()
	compiler := intent.NewCompiler(repo)

	_, report, err := compiler.Compile(context.Background(), intent.IntentInput{
		RequestID:             "req-invalid",
		TenantID:              "tenant-1",
		WorkspaceID:           "workspace-1",
		UserID:                "user-1",
		SessionID:             "session-1",
		Objetivo:              "Contrato inválido",
		Alcance:               "test-scope",
		TipoResultadoEsperado: intent.ResultType("invalid_type"),
		AutonomiaSolicitada:   intent.AutonomyAssisted,
		CreatedAt:             time.Now().UTC(),
	})
	if err == nil {
		t.Fatalf("expected error for invalid result type")
	}
	if report.Status != intent.CompilationStatusRejected {
		t.Fatalf("expected rejected status, got %s", report.Status)
	}
}
