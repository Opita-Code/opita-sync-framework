package intent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var ErrInvalidIntentInput = errors.New("invalid intent input")

type ContractRepository interface {
	GetByFingerprint(ctx context.Context, fingerprint string) (CompiledContract, bool, error)
	GetByID(ctx context.Context, contractID string) (CompiledContract, bool, error)
	Save(ctx context.Context, contract CompiledContract) error
}

type Compiler struct {
	repo ContractRepository
}

func NewCompiler(repo ContractRepository) *Compiler {
	return &Compiler{repo: repo}
}

func (c *Compiler) Compile(ctx context.Context, input IntentInput) (CompiledContract, CompilationReport, error) {
	canonical := normalize(input)
	diagnostics := validate(canonical)
	if hasErrors(diagnostics) {
		report := CompilationReport{
			Status:      CompilationStatusRejected,
			Diagnostics: diagnostics,
			CompiledAt:  time.Now().UTC(),
		}
		return CompiledContract{}, report, ErrInvalidIntentInput
	}

	snapshots := SnapshotBundle{
		PolicyVersion:         "policy-v1",
		ClassificationVersion: "classification-v1",
		RiskVersion:           "risk-v1",
		PermissionVersion:     "permission-v1",
		CompiledAt:            time.Now().UTC(),
	}

	fingerprint, err := calculateFingerprint(canonical)
	if err != nil {
		return CompiledContract{}, CompilationReport{}, fmt.Errorf("calculate fingerprint: %w", err)
	}

	if existing, found, err := c.repo.GetByFingerprint(ctx, fingerprint); err != nil {
		return CompiledContract{}, CompilationReport{}, fmt.Errorf("lookup by fingerprint: %w", err)
	} else if found {
		report := CompilationReport{
			ContractID:   existing.ContractID,
			Fingerprint:  existing.Fingerprint,
			Status:       CompilationStatusDuplicate,
			Diagnostics:  []Diagnostic{{Code: "compiler.duplicate_fingerprint", Severity: SeverityInfo, Message: "same normalized input and snapshots resolved to an existing contract"}},
			DuplicatedOf: existing.ContractID,
			CompiledAt:   time.Now().UTC(),
		}
		return existing, report, nil
	}

	now := time.Now().UTC()
	contract := CompiledContract{
		ContractID:            buildContractID(canonical),
		ContractVersion:       "1.0",
		RequestID:             canonical.RequestID,
		TenantID:              canonical.TenantID,
		WorkspaceID:           canonical.WorkspaceID,
		UserID:                canonical.UserID,
		SessionID:             canonical.SessionID,
		TraceID:               canonical.TraceID,
		ConversationTurnID:    canonical.ConversationTurnID,
		IntakeSessionID:       canonical.IntakeSessionID,
		IntentCandidateID:     canonical.IntentCandidateID,
		ProposalDraftID:       canonical.ProposalDraftID,
		PatchsetCandidateID:   canonical.PatchsetCandidateID,
		PreviewCandidateID:    canonical.PreviewCandidateID,
		SimulationResultIDs:   canonical.SimulationResultIDs,
		Objetivo:              canonical.Objetivo,
		Alcance:               canonical.Alcance,
		TipoResultadoEsperado: canonical.TipoResultadoEsperado,
		Restricciones:         canonical.Restricciones,
		SistemasPosibles:      canonical.SistemasPosibles,
		DatosPermitidos:       canonical.DatosPermitidos,
		AutonomiaSolicitada:   canonical.AutonomiaSolicitada,
		AprobacionRequerida:   canonical.AprobacionRequerida,
		CriteriosDeExito:      canonical.CriteriosDeExito,
		Fingerprint:           fingerprint,
		State:                 ContractStateCompiled,
		Snapshots:             snapshots,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := c.repo.Save(ctx, contract); err != nil {
		return CompiledContract{}, CompilationReport{}, fmt.Errorf("persist contract: %w", err)
	}

	report := CompilationReport{
		ContractID:  contract.ContractID,
		Fingerprint: contract.Fingerprint,
		Status:      CompilationStatusCompiled,
		Diagnostics: append(diagnostics, Diagnostic{Code: "compiler.compiled", Severity: SeverityInfo, Message: "contract compiled successfully"}),
		CompiledAt:  snapshots.CompiledAt,
	}

	return contract, report, nil
}

func normalize(input IntentInput) CanonicalIntent {
	canonical := CanonicalIntent{
		RequestID:             strings.TrimSpace(input.RequestID),
		TenantID:              strings.TrimSpace(input.TenantID),
		WorkspaceID:           strings.TrimSpace(input.WorkspaceID),
		UserID:                strings.TrimSpace(input.UserID),
		SessionID:             strings.TrimSpace(input.SessionID),
		TraceID:               strings.TrimSpace(input.TraceID),
		ConversationTurnID:    strings.TrimSpace(input.ConversationTurnID),
		IntakeSessionID:       strings.TrimSpace(input.IntakeSessionID),
		IntentCandidateID:     strings.TrimSpace(input.IntentCandidateID),
		ProposalDraftID:       strings.TrimSpace(input.ProposalDraftID),
		PatchsetCandidateID:   strings.TrimSpace(input.PatchsetCandidateID),
		PreviewCandidateID:    strings.TrimSpace(input.PreviewCandidateID),
		SimulationResultIDs:   normalizeStringSlice(input.SimulationResultIDs),
		Objetivo:              strings.TrimSpace(input.Objetivo),
		Alcance:               strings.TrimSpace(input.Alcance),
		TipoResultadoEsperado: input.TipoResultadoEsperado,
		Restricciones:         normalizeStringSlice(input.Restricciones),
		SistemasPosibles:      normalizeStringSlice(input.SistemasPosibles),
		DatosPermitidos:       normalizeStringSlice(input.DatosPermitidos),
		AutonomiaSolicitada:   input.AutonomiaSolicitada,
		AprobacionRequerida:   strings.TrimSpace(input.AprobacionRequerida),
		CriteriosDeExito:      normalizeStringSlice(input.CriteriosDeExito),
		CreatedAt:             input.CreatedAt.UTC(),
	}
	if canonical.AutonomiaSolicitada == "" {
		canonical.AutonomiaSolicitada = AutonomyAssisted
	}
	return canonical
}

func validate(intent CanonicalIntent) []Diagnostic {
	var diagnostics []Diagnostic
	if intent.RequestID == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_request_id", Severity: SeverityError, Field: "request_id", Message: "request_id is required"})
	}
	if intent.TenantID == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_tenant_id", Severity: SeverityError, Field: "tenant_id", Message: "tenant_id is required"})
	}
	if intent.WorkspaceID == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_workspace_id", Severity: SeverityError, Field: "workspace_id", Message: "workspace_id is required"})
	}
	if intent.UserID == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_user_id", Severity: SeverityError, Field: "user_id", Message: "user_id is required"})
	}
	if intent.SessionID == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_session_id", Severity: SeverityError, Field: "session_id", Message: "session_id is required"})
	}
	if intent.Objetivo == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_objetivo", Severity: SeverityError, Field: "objetivo", Message: "objetivo is required"})
	}
	if intent.Alcance == "" {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.missing_alcance", Severity: SeverityError, Field: "alcance", Message: "alcance is required"})
	}
	if !isValidResultType(intent.TipoResultadoEsperado) {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.invalid_result_type", Severity: SeverityError, Field: "tipo_de_resultado_esperado", Message: "tipo_de_resultado_esperado is invalid"})
	}
	if intent.CreatedAt.IsZero() {
		diagnostics = append(diagnostics, Diagnostic{Code: "compiler.zero_created_at", Severity: SeverityWarning, Field: "created_at", Message: "created_at was zero and may reduce trace quality"})
	}
	return diagnostics
}

func calculateFingerprint(intent CanonicalIntent) (string, error) {
	payload := map[string]any{
		"request_id":                 intent.RequestID,
		"tenant_id":                  intent.TenantID,
		"workspace_id":               intent.WorkspaceID,
		"user_id":                    intent.UserID,
		"session_id":                 intent.SessionID,
		"objetivo":                   intent.Objetivo,
		"alcance":                    intent.Alcance,
		"tipo_de_resultado_esperado": intent.TipoResultadoEsperado,
		"restricciones":              intent.Restricciones,
		"sistemas_posibles":          intent.SistemasPosibles,
		"datos_permitidos":           intent.DatosPermitidos,
		"autonomia_solicitada":       intent.AutonomiaSolicitada,
		"aprobacion_requerida":       intent.AprobacionRequerida,
		"criterios_de_exito":         intent.CriteriosDeExito,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(raw)
	return hex.EncodeToString(hash[:]), nil
}

func normalizeStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		set[normalized] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func buildContractID(intent CanonicalIntent) string {
	base := strings.ReplaceAll(strings.ToLower(intent.TenantID), " ", "-")
	if base == "" {
		base = "tenant"
	}
	return fmt.Sprintf("contract-%s-%d", base, time.Now().UTC().UnixNano())
}

func hasErrors(diagnostics []Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}
	return false
}

func isValidResultType(resultType ResultType) bool {
	switch resultType {
	case ResultTypePlan,
		ResultTypeInspection,
		ResultTypeQuery,
		ResultTypeReport,
		ResultTypeChangeProposal,
		ResultTypeExecution,
		ResultTypeSystemUpdate,
		ResultTypeGovernanceDecision:
		return true
	default:
		return false
	}
}
