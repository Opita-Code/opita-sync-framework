package postgres

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"opita-sync-framework/internal/engine/approvals"
	"opita-sync-framework/internal/engine/events"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/inspection"
	"opita-sync-framework/internal/engine/intake"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/maintenance"
	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/preview"
	"opita-sync-framework/internal/engine/proposal"
	"opita-sync-framework/internal/engine/registry"
	"opita-sync-framework/internal/engine/runtime"
)

func newMockStore(t *testing.T) (*Store, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	cleanup := func() { _ = db.Close() }
	return &Store{DB: db}, mock, cleanup
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return raw
}

func requireExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestContractRepositorySaveAndGetRoundTrip(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	repo := NewContractRepository(store)
	now := time.Date(2026, 3, 29, 10, 0, 0, 0, time.UTC)
	contract := intent.CompiledContract{
		ContractID:      "contract-1",
		ContractVersion: "v1",
		RequestID:       "request-1",
		TenantID:        "tenant-1",
		WorkspaceID:     "workspace-1",
		UserID:          "user-1",
		SessionID:       "session-1",
		Objetivo:        "actualizar catalogo",
		Alcance:         "tenant scope",
		Fingerprint:     "fp-1",
		State:           intent.ContractStateCompiled,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	raw := mustJSON(t, contract)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into compiled_contracts (contract_id, fingerprint, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		on conflict (contract_id) do update set
		  fingerprint = excluded.fingerprint,
		  payload = excluded.payload,
		  updated_at = excluded.updated_at
	`)).
		WithArgs(contract.ContractID, contract.Fingerprint, raw, contract.CreatedAt, contract.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Save(contextBackground(), contract); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from compiled_contracts where contract_id = $1`)).
		WithArgs(contract.ContractID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))

	got, found, err := repo.GetByID(contextBackground(), contract.ContractID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if !found {
		t.Fatal("expected contract to be found")
	}
	if got.ContractID != contract.ContractID || got.Fingerprint != contract.Fingerprint || got.TenantID != contract.TenantID {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, contract)
	}

	requireExpectations(t, mock)
}

func TestRuntimeServiceCreateGetAndUpdateRoundTrip(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	svc := NewRuntimeService(store)
	now := time.Date(2026, 3, 29, 11, 0, 0, 0, time.UTC)
	record := runtime.ExecutionRecord{
		ExecutionID:         "exec-1",
		TenantID:            "tenant-1",
		ContractID:          "contract-1",
		ContractFingerprint: "fp-1",
		TraceID:             "trace-1",
		State:               runtime.ExecutionStateAwaitingApproval,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	raw := mustJSON(t, record)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into execution_records (execution_id, contract_id, tenant_id, trace_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).
		WithArgs(record.ExecutionID, record.ContractID, record.TenantID, record.TraceID, record.State, raw, record.CreatedAt, record.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := svc.CreateExecution(record); err != nil {
		t.Fatalf("CreateExecution returned error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from execution_records where execution_id = $1`)).
		WithArgs(record.ExecutionID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))

	got, found, err := svc.GetExecution(record.ExecutionID)
	if err != nil {
		t.Fatalf("GetExecution returned error: %v", err)
	}
	if !found {
		t.Fatal("expected execution to be found")
	}
	if got.ExecutionID != record.ExecutionID || got.ContractID != record.ContractID || got.TraceID != record.TraceID {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, record)
	}

	updated := record
	updated.State = runtime.ExecutionStateExecutionReleased
	updatedRaw := mustJSON(t, updated)

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from execution_records where execution_id = $1`)).
		WithArgs(record.ExecutionID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))
	mock.ExpectExec(regexp.QuoteMeta(`
		update execution_records set state = $2, payload = $3, updated_at = now() where execution_id = $1
	`)).
		WithArgs(record.ExecutionID, updated.State, updatedRaw).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := svc.UpdateExecutionState(record.ExecutionID, updated.State)
	if err != nil {
		t.Fatalf("UpdateExecutionState returned error: %v", err)
	}
	if result.ExecutionID != record.ExecutionID || result.State != updated.State {
		t.Fatalf("update result mismatch: %+v", result)
	}

	requireExpectations(t, mock)
}

func TestApprovalStoreCreateGetAndDecideRoundTrip(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	svc := NewApprovalStore(store)
	now := time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC)
	request := approvals.Request{
		ApprovalRequestID:         "approval-1",
		ExecutionID:               "exec-1",
		ContractID:                "contract-1",
		TenantID:                  "tenant-1",
		TraceID:                   "trace-1",
		State:                     approvals.StateAwaitingApproval,
		Mode:                      "pre_execution",
		ReasonCodes:               []string{"policy.needs_approval"},
		SourceContractFingerprint: "fp-1",
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	raw := mustJSON(t, request)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into approval_requests (approval_request_id, execution_id, contract_id, tenant_id, trace_id, state, mode, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)).
		WithArgs(request.ApprovalRequestID, request.ExecutionID, request.ContractID, request.TenantID, request.TraceID, request.State, request.Mode, raw, request.CreatedAt, request.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := svc.Create(request); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from approval_requests where approval_request_id = $1`)).
		WithArgs(request.ApprovalRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))

	got, found, err := svc.GetByID(request.ApprovalRequestID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if !found {
		t.Fatal("expected approval to be found")
	}
	if got.ExecutionID != request.ExecutionID || got.ContractID != request.ContractID || got.TraceID != request.TraceID {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, request)
	}

	updated := request
	updated.State = approvals.StateReleased
	updated.DecidedBySubjectID = "approver-1"
	updated.DecisionComment = "looks safe"
	updated.DecisionReasonCodes = []string{"approval.release.manual"}
	updated.UpdatedAt = time.Now().UTC()
	updated.ReleasedAt = &updated.UpdatedAt
	updatedRaw := mustJSON(t, updated)

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from approval_requests where approval_request_id = $1`)).
		WithArgs(request.ApprovalRequestID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))
	mock.ExpectExec(regexp.QuoteMeta(`
		update approval_requests set state = $2, payload = $3, updated_at = now() where approval_request_id = $1
	`)).
		WithArgs(request.ApprovalRequestID, updated.State, updatedRaw).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := svc.Decide(request.ApprovalRequestID, approvals.Decision{
		State:               updated.State,
		DecidedBySubjectID:  updated.DecidedBySubjectID,
		DecisionComment:     updated.DecisionComment,
		DecisionReasonCodes: updated.DecisionReasonCodes,
		DecidedAt:           updated.UpdatedAt,
	})
	if err != nil {
		t.Fatalf("Decide returned error: %v", err)
	}
	if result.State != approvals.StateReleased || result.DecidedBySubjectID != updated.DecidedBySubjectID || result.SourceContractFingerprint != request.SourceContractFingerprint {
		t.Fatalf("expected enriched released state, got %+v", result)
	}

	requireExpectations(t, mock)
}

func TestEventLogAppendAndReadByExecution(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	log := NewEventLog(store)
	now := time.Date(2026, 3, 29, 13, 0, 0, 0, time.UTC)
	record := events.Record{
		EventID:             "event-1",
		EventType:           "execution.created",
		TenantID:            "tenant-1",
		TraceID:             "trace-1",
		ContractID:          "contract-1",
		ContractFingerprint: "fp-1",
		ExecutionID:         "exec-1",
		PolicyDecisionID:    "policy-1",
		OccurredAt:          now,
		Payload:             map[string]any{"state": "created"},
	}
	raw := mustJSON(t, record)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into event_records (event_id, execution_id, event_type, trace_id, tenant_id, occurred_at, payload)
		values ($1, $2, $3, $4, $5, $6, $7)
	`)).
		WithArgs(record.EventID, record.ExecutionID, record.EventType, record.TraceID, record.TenantID, record.OccurredAt, raw).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := log.Append(record); err != nil {
		t.Fatalf("Append returned error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from event_records where execution_id = $1 order by occurred_at asc`)).
		WithArgs(record.ExecutionID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))

	got := log.RecordsByExecution(record.ExecutionID)
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
	if got[0].ExecutionID != record.ExecutionID || got[0].ContractID != record.ContractID || got[0].TraceID != record.TraceID {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got[0], record)
	}

	requireExpectations(t, mock)
}

func TestFoundationRunStoreRoundTrip(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	svc := NewFoundationRunStore(store)
	now := time.Date(2026, 3, 29, 14, 0, 0, 0, time.UTC)
	run := foundation.FoundationRunResult{
		Contract:       intent.CompiledContract{ContractID: "contract-1", TenantID: "tenant-1", Fingerprint: "fp-1"},
		Execution:      runtime.ExecutionRecord{ExecutionID: "exec-1", TenantID: "tenant-1", ContractID: "contract-1", TraceID: "trace-1", State: runtime.ExecutionStateExecutionReleased, CreatedAt: now, UpdatedAt: now},
		PolicyDecision: policy.DecisionRecord{PolicyDecisionID: "policy-1"},
		Resolution:     registry.ResolutionResult{CapabilityManifestRef: "capability.execution.default", BindingID: "binding-1", ProviderRef: "provider-1"},
	}
	raw := mustJSON(t, run)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into foundation_runs (execution_id, contract_id, trace_id, payload, created_at)
		values ($1, $2, $3, $4, $5)
		on conflict (execution_id) do update set
		  payload = excluded.payload
	`)).
		WithArgs(run.Execution.ExecutionID, run.Contract.ContractID, run.Execution.TraceID, raw, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := svc.Save(run); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from foundation_runs where execution_id = $1`)).
		WithArgs(run.Execution.ExecutionID).
		WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(raw))

	got, found, err := svc.GetByExecutionID(run.Execution.ExecutionID)
	if err != nil {
		t.Fatalf("GetByExecutionID returned error: %v", err)
	}
	if !found {
		t.Fatal("expected foundation run to be found")
	}
	if got.Execution.ExecutionID != run.Execution.ExecutionID || got.Contract.ContractID != run.Contract.ContractID || got.Execution.TraceID != run.Execution.TraceID {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, run)
	}

	requireExpectations(t, mock)
}

func TestSurfaceStoresRoundTrip(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	intakeStore := NewIntakeStore(store)
	proposalStore := NewProposalStore(store)
	previewStore := NewPreviewStore(store)
	recoveryStore := NewRecoveryStore(store)
	maintenanceStore := NewMaintenanceStore(store)

	now := time.Date(2026, 3, 29, 15, 0, 0, 0, time.UTC)
	turn := intake.ConversationTurn{ConversationTurnID: "turn-1", SessionID: "session-1", TenantID: "tenant-1", SubjectID: "subject-1", RawText: "hola", Timestamp: now, TraceID: "trace-1"}
	session := intake.Session{IntakeSessionID: "is-1", SessionID: "session-1", TenantID: "tenant-1", SubjectID: "subject-1", CurrentState: intake.SessionStateShaping, TraceID: "trace-1", UpdatedAt: now}
	candidate := intake.IntentCandidate{IntentCandidateID: "candidate-1", SourceTurnIDs: []string{"turn-1"}, ObjetivoCandidate: "actualizar", ReadyForProposalDraft: true}
	draft := proposal.Draft{ProposalDraftID: "draft-1", TenantID: "tenant-1", SessionID: "session-1", SubjectID: "subject-1", Title: "Cambio", CurrentState: proposal.DraftStateOpen, CreatedAt: now}
	patchset := proposal.PatchsetCandidate{PatchsetCandidateID: "patch-1", ProposalDraftID: draft.ProposalDraftID, ReadyForPreview: true, CreatedAt: now}
	previewCandidate := preview.Candidate{PreviewCandidateID: "preview-1", TenantID: "tenant-1", SessionID: "session-1", SubjectID: "subject-1", ContractID: "contract-1", ExecutionID: "exec-1", State: preview.StatusPreviewOK, CreatedAt: now}
	simulation := preview.Result{SimulationResultID: "simulation-1", PreviewCandidateID: previewCandidate.PreviewCandidateID, Family: preview.SimulationPolicy, Status: preview.StatusPreviewOK, CreatedAt: now}
	recovery := inspection.RecoveryActionCandidate{RecoveryActionCandidateID: "recovery-1", ExecutionID: "exec-1", RequestedAction: inspection.RecoveryResumeAfterApproval, RequestedBySubjectID: "subject-1", CurrentRuntimeState: "awaiting_approval", ReadyForExecution: true, State: inspection.RecoveryCandidatePending, CreatedAt: now, UpdatedAt: now}
	maint := maintenance.ActionCandidate{MaintenanceActionCandidateID: "maintenance-1", TenantID: "tenant-1", RequestedBySubjectID: "subject-1", CreatedAt: now}

	turnRaw := mustJSON(t, turn)
	sessionRaw := mustJSON(t, session)
	candidateRaw := mustJSON(t, candidate)
	draftRaw := mustJSON(t, draft)
	patchsetRaw := mustJSON(t, patchset)
	previewRaw := mustJSON(t, previewCandidate)
	simulationRaw := mustJSON(t, simulation)
	recoveryRaw := mustJSON(t, recovery)
	maintRaw := mustJSON(t, maint)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into intake_turns (conversation_turn_id, session_id, tenant_id, subject_id, trace_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6, $7)
	`)).WithArgs(turn.ConversationTurnID, turn.SessionID, turn.TenantID, turn.SubjectID, turn.TraceID, turnRaw, turn.Timestamp).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := intakeStore.CreateTurn(turn); err != nil {
		t.Fatalf("CreateTurn: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into intake_sessions (intake_session_id, session_id, tenant_id, subject_id, trace_id, payload, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)
		on conflict (intake_session_id) do update set payload = excluded.payload, updated_at = excluded.updated_at
	`)).WithArgs(session.IntakeSessionID, session.SessionID, session.TenantID, session.SubjectID, session.TraceID, sessionRaw, session.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := intakeStore.CreateSession(session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into intent_candidates (intent_candidate_id, payload, created_at)
		values ($1, $2, $3)
		on conflict (intent_candidate_id) do update set payload = excluded.payload
	`)).WithArgs(candidate.IntentCandidateID, candidateRaw, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := intakeStore.SaveIntentCandidate(candidate); err != nil {
		t.Fatalf("SaveIntentCandidate: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from intake_sessions where intake_session_id = $1`)).WithArgs(session.IntakeSessionID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(sessionRaw))
	gotSession, found, err := intakeStore.GetSession(session.IntakeSessionID)
	if err != nil || !found || gotSession.TraceID != session.TraceID {
		t.Fatalf("GetSession mismatch: found=%v err=%v got=%+v", found, err, gotSession)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from intent_candidates where intent_candidate_id = $1`)).WithArgs(candidate.IntentCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(candidateRaw))
	gotCandidate, found, err := intakeStore.GetIntentCandidate(candidate.IntentCandidateID)
	if err != nil || !found || gotCandidate.IntentCandidateID != candidate.IntentCandidateID {
		t.Fatalf("GetIntentCandidate mismatch: found=%v err=%v got=%+v", found, err, gotCandidate)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into proposal_drafts (proposal_draft_id, tenant_id, session_id, subject_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`)).WithArgs(draft.ProposalDraftID, draft.TenantID, draft.SessionID, draft.SubjectID, draftRaw, draft.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := proposalStore.CreateDraft(draft); err != nil {
		t.Fatalf("CreateDraft: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into patchset_candidates (patchset_candidate_id, proposal_draft_id, payload, created_at)
		values ($1, $2, $3, $4)
		on conflict (patchset_candidate_id) do update set payload = excluded.payload
	`)).WithArgs(patchset.PatchsetCandidateID, patchset.ProposalDraftID, patchsetRaw, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := proposalStore.SavePatchset(patchset); err != nil {
		t.Fatalf("SavePatchset: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from proposal_drafts where proposal_draft_id = $1`)).WithArgs(draft.ProposalDraftID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(draftRaw))
	gotDraft, found, err := proposalStore.GetDraft(draft.ProposalDraftID)
	if err != nil || !found || gotDraft.ProposalDraftID != draft.ProposalDraftID {
		t.Fatalf("GetDraft mismatch: found=%v err=%v got=%+v", found, err, gotDraft)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from patchset_candidates where patchset_candidate_id = $1`)).WithArgs(patchset.PatchsetCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(patchsetRaw))
	gotPatchset, found, err := proposalStore.GetPatchset(patchset.PatchsetCandidateID)
	if err != nil || !found || gotPatchset.PatchsetCandidateID != patchset.PatchsetCandidateID {
		t.Fatalf("GetPatchset mismatch: found=%v err=%v got=%+v", found, err, gotPatchset)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into preview_candidates (preview_candidate_id, tenant_id, session_id, subject_id, contract_id, execution_id, payload, created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).WithArgs(previewCandidate.PreviewCandidateID, previewCandidate.TenantID, previewCandidate.SessionID, previewCandidate.SubjectID, previewCandidate.ContractID, previewCandidate.ExecutionID, previewRaw, previewCandidate.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := previewStore.CreateCandidate(previewCandidate); err != nil {
		t.Fatalf("CreateCandidate: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into simulation_results (simulation_result_id, preview_candidate_id, family, payload, created_at)
		values ($1, $2, $3, $4, $5)
		on conflict (simulation_result_id) do update set payload = excluded.payload
	`)).WithArgs(simulation.SimulationResultID, simulation.PreviewCandidateID, simulation.Family, simulationRaw, simulation.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := previewStore.SaveResult(simulation); err != nil {
		t.Fatalf("SaveResult: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from preview_candidates where preview_candidate_id = $1`)).WithArgs(previewCandidate.PreviewCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(previewRaw))
	gotPreview, found, err := previewStore.GetCandidate(previewCandidate.PreviewCandidateID)
	if err != nil || !found || gotPreview.ExecutionID != previewCandidate.ExecutionID || gotPreview.ContractID != previewCandidate.ContractID {
		t.Fatalf("GetCandidate mismatch: found=%v err=%v got=%+v", found, err, gotPreview)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from simulation_results where preview_candidate_id = $1 order by created_at asc`)).WithArgs(previewCandidate.PreviewCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(simulationRaw))
	results, err := previewStore.ListResults(previewCandidate.PreviewCandidateID)
	if err != nil || len(results) != 1 || results[0].PreviewCandidateID != previewCandidate.PreviewCandidateID {
		t.Fatalf("ListResults mismatch: err=%v results=%+v", err, results)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into recovery_action_candidates (recovery_action_candidate_id, execution_id, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
	`)).WithArgs(recovery.RecoveryActionCandidateID, recovery.ExecutionID, recoveryRaw, recovery.CreatedAt, recovery.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := recoveryStore.Create(recovery); err != nil {
		t.Fatalf("Recovery Create: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from recovery_action_candidates where recovery_action_candidate_id = $1`)).WithArgs(recovery.RecoveryActionCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(recoveryRaw))
	gotRecovery, found, err := recoveryStore.GetByID(recovery.RecoveryActionCandidateID)
	if err != nil || !found || gotRecovery.ExecutionID != recovery.ExecutionID {
		t.Fatalf("Recovery GetByID mismatch: found=%v err=%v got=%+v", found, err, gotRecovery)
	}

	recovery.State = inspection.RecoveryCandidateExecuted
	recovery.UpdatedAt = now.Add(time.Minute)
	recoveryUpdatedRaw := mustJSON(t, recovery)
	mock.ExpectExec(regexp.QuoteMeta(`
		update recovery_action_candidates set payload = $2, updated_at = $3 where recovery_action_candidate_id = $1
	`)).WithArgs(recovery.RecoveryActionCandidateID, recoveryUpdatedRaw, recovery.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := recoveryStore.Update(recovery); err != nil {
		t.Fatalf("Recovery Update: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into maintenance_action_candidates (maintenance_action_candidate_id, tenant_id, payload, created_at)
		values ($1, $2, $3, $4)
	`)).WithArgs(maint.MaintenanceActionCandidateID, maint.TenantID, maintRaw, maint.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := maintenanceStore.Create(maint); err != nil {
		t.Fatalf("Maintenance Create: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select payload from maintenance_action_candidates where maintenance_action_candidate_id = $1`)).WithArgs(maint.MaintenanceActionCandidateID).WillReturnRows(sqlmock.NewRows([]string{"payload"}).AddRow(maintRaw))
	gotMaint, found, err := maintenanceStore.GetByID(maint.MaintenanceActionCandidateID)
	if err != nil || !found || gotMaint.MaintenanceActionCandidateID != maint.MaintenanceActionCandidateID {
		t.Fatalf("Maintenance GetByID mismatch: found=%v err=%v got=%+v", found, err, gotMaint)
	}

	requireExpectations(t, mock)
}

func TestCriticalStoresReturnTraceablePersistenceErrors(t *testing.T) {
	store, mock, cleanup := newMockStore(t)
	defer cleanup()

	contractRepo := NewContractRepository(store)
	runtimeSvc := NewRuntimeService(store)
	eventLog := NewEventLog(store)
	approvalStore := NewApprovalStore(store)

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into compiled_contracts (contract_id, fingerprint, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		on conflict (contract_id) do update set
		  fingerprint = excluded.fingerprint,
		  payload = excluded.payload,
		  updated_at = excluded.updated_at
	`)).WillReturnError(sql.ErrConnDone)
	if err := contractRepo.Save(contextBackground(), intent.CompiledContract{ContractID: "c", Fingerprint: "fp", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err == nil || !regexp.MustCompile(`upsert compiled contract`).MatchString(err.Error()) {
		t.Fatalf("expected traceable contract error, got %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into execution_records (execution_id, contract_id, tenant_id, trace_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).WillReturnError(sql.ErrConnDone)
	if err := runtimeSvc.CreateExecution(runtime.ExecutionRecord{ExecutionID: "e", ContractID: "c", TenantID: "t", TraceID: "tr", State: runtime.ExecutionStateCreated, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err == nil || !regexp.MustCompile(`insert execution record`).MatchString(err.Error()) {
		t.Fatalf("expected traceable runtime error, got %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into event_records (event_id, execution_id, event_type, trace_id, tenant_id, occurred_at, payload)
		values ($1, $2, $3, $4, $5, $6, $7)
	`)).WillReturnError(sql.ErrConnDone)
	if err := eventLog.Append(events.Record{EventID: "ev", ExecutionID: "e", EventType: "event", TraceID: "tr", TenantID: "t", OccurredAt: time.Now()}); err == nil || !regexp.MustCompile(`insert event record`).MatchString(err.Error()) {
		t.Fatalf("expected traceable event error, got %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		insert into approval_requests (approval_request_id, execution_id, contract_id, tenant_id, trace_id, state, mode, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)).WillReturnError(sql.ErrConnDone)
	if err := approvalStore.Create(approvals.Request{ApprovalRequestID: "a", ExecutionID: "e", ContractID: "c", TenantID: "t", TraceID: "tr", State: approvals.StateAwaitingApproval, Mode: "pre_execution", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err == nil || !regexp.MustCompile(`insert approval request`).MatchString(err.Error()) {
		t.Fatalf("expected traceable approval error, got %v", err)
	}

	requireExpectations(t, mock)
}
