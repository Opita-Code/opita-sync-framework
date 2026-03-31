create table if not exists compiled_contracts (
  contract_id text primary key,
  fingerprint text not null unique,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table if not exists execution_records (
  execution_id text primary key,
  contract_id text not null,
  tenant_id text not null,
  trace_id text not null,
  state text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_execution_records_contract_id on execution_records (contract_id);
create index if not exists idx_execution_records_tenant_id on execution_records (tenant_id);
create index if not exists idx_execution_records_trace_id on execution_records (trace_id);

create table if not exists event_records (
  event_id text primary key,
  execution_id text not null,
  event_type text not null,
  trace_id text not null,
  tenant_id text not null,
  occurred_at timestamptz not null,
  payload jsonb not null
);

create index if not exists idx_event_records_execution_id on event_records (execution_id);
create index if not exists idx_event_records_trace_id on event_records (trace_id);
create index if not exists idx_event_records_tenant_id on event_records (tenant_id);
create index if not exists idx_event_records_contract_id on event_records ((payload->>'ContractID'));
create index if not exists idx_event_records_occurred_at on event_records (occurred_at);

create table if not exists foundation_runs (
  execution_id text primary key,
  contract_id text not null,
  trace_id text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_foundation_runs_contract_id on foundation_runs (contract_id);
create index if not exists idx_foundation_runs_trace_id on foundation_runs (trace_id);

create table if not exists approval_requests (
  approval_request_id text primary key,
  execution_id text not null,
  contract_id text not null,
  tenant_id text not null,
  trace_id text not null,
  state text not null,
  mode text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_approval_requests_execution_id on approval_requests (execution_id);
create index if not exists idx_approval_requests_contract_id on approval_requests (contract_id);
create index if not exists idx_approval_requests_tenant_id on approval_requests (tenant_id);
create index if not exists idx_approval_requests_trace_id on approval_requests (trace_id);
create index if not exists idx_approval_requests_state on approval_requests (state);

create table if not exists intake_turns (
  conversation_turn_id text primary key,
  session_id text not null,
  tenant_id text not null,
  subject_id text not null,
  trace_id text,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_intake_turns_session_id on intake_turns (session_id);
create index if not exists idx_intake_turns_tenant_id on intake_turns (tenant_id);
create index if not exists idx_intake_turns_trace_id on intake_turns (trace_id);

create table if not exists intake_sessions (
  intake_session_id text primary key,
  session_id text not null,
  tenant_id text not null,
  subject_id text not null,
  trace_id text,
  payload jsonb not null,
  updated_at timestamptz not null
);

create index if not exists idx_intake_sessions_session_id on intake_sessions (session_id);
create index if not exists idx_intake_sessions_tenant_id on intake_sessions (tenant_id);
create index if not exists idx_intake_sessions_trace_id on intake_sessions (trace_id);

create table if not exists intent_candidates (
  intent_candidate_id text primary key,
  payload jsonb not null,
  created_at timestamptz not null
);

create table if not exists proposal_drafts (
  proposal_draft_id text primary key,
  tenant_id text not null,
  session_id text not null,
  subject_id text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_proposal_drafts_tenant_id on proposal_drafts (tenant_id);
create index if not exists idx_proposal_drafts_session_id on proposal_drafts (session_id);

create table if not exists patchset_candidates (
  patchset_candidate_id text primary key,
  proposal_draft_id text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_patchset_candidates_proposal_draft_id on patchset_candidates (proposal_draft_id);

create table if not exists preview_candidates (
  preview_candidate_id text primary key,
  tenant_id text not null,
  session_id text not null,
  subject_id text not null,
  contract_id text not null,
  execution_id text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_preview_candidates_tenant_id on preview_candidates (tenant_id);
create index if not exists idx_preview_candidates_session_id on preview_candidates (session_id);
create index if not exists idx_preview_candidates_contract_id on preview_candidates (contract_id);
create index if not exists idx_preview_candidates_execution_id on preview_candidates (execution_id);

create table if not exists simulation_results (
  simulation_result_id text primary key,
  preview_candidate_id text not null,
  family text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_simulation_results_preview_candidate_id on simulation_results (preview_candidate_id);
create index if not exists idx_simulation_results_created_at on simulation_results (created_at);

create table if not exists recovery_action_candidates (
  recovery_action_candidate_id text primary key,
  execution_id text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_recovery_action_candidates_execution_id on recovery_action_candidates (execution_id);

create table if not exists maintenance_action_candidates (
  maintenance_action_candidate_id text primary key,
  tenant_id text not null,
  payload jsonb not null,
  created_at timestamptz not null
);

create index if not exists idx_maintenance_action_candidates_tenant_id on maintenance_action_candidates (tenant_id);

create table if not exists tenant_bootstrap_records (
  tenant_id text primary key,
  state text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_tenant_bootstrap_records_state on tenant_bootstrap_records (state);

create table if not exists tenant_capability_grants (
  grant_id text primary key,
  tenant_id text not null,
  state text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_tenant_capability_grants_tenant_id on tenant_capability_grants (tenant_id);
create index if not exists idx_tenant_capability_grants_state on tenant_capability_grants (state);

create table if not exists tenant_delegation_grants (
  grant_id text primary key,
  tenant_id text not null,
  state text not null,
  payload jsonb not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists idx_tenant_delegation_grants_tenant_id on tenant_delegation_grants (tenant_id);
create index if not exists idx_tenant_delegation_grants_state on tenant_delegation_grants (state);
