# Tasks: OSF Consistency Improvement

## Phase 1: Foundation — Mechanical Extractions

### D2: Shared HTTP Utilities

- [x] T2.1 Create `internal/httputil/` package with `WriteJSON`, `WriteError`, `EncodeJSON`, `ParseJSON`
- [x] T2.2 Migrate all 9 app packages to use `httputil`
- [x] T2.3 Remove duplicated `writeJSON`/`writeError` from all 9 packages

### D3: Shared Test Utilities

- [x] T3.1 Create `internal/testutil/` package with `GetStringField`, `PostJSON`, `GetJSON`, `ErrorContains`
- [x] T3.2 Migrate test files to use `testutil` helpers
- [x] T3.3 Remove duplicated test helpers from app packages

### D9: Dead Code Cleanup

- [x] T9.1 Remove `internal/cache/store/` package
- [x] T9.2 Remove references to cache from main.go and any other consumers
- [x] T9.3 Verify build passes after removal

## Phase 2: Consistency — Interface & Style

- [x] T4.1 Add `var _ Interface = (*Impl)(nil)` checks to all platform implementations in `internal/platform/*/`
- [x] T4.2 Verify compile-time compliance
- [x] T5.1 Fix YAML formatting in `definitions/capabilities/*.yaml` (JSON → YAML)
- [x] T5.2 Verify round-trip equivalence
- [ ] T6.1 Review all handler packages for consistent error format
- [ ] T6.2 Align flat vs nested error shapes across packages

## Phase 3: Observability & Shutdown

- [x] D4.1 Refactor `internal/engine/simulation/service.go` — replace hardcoded results with real `Policy.Evaluate()` calls
- [x] D4.2 Use variant `Action`/`RiskLevel` per result family (policy, approval, classification, risk)
- [x] D4.3 Handle unknown/missing candidate input — return error instead of panic
- [x] T7.1 Replace `log.Printf`/`log.Fatalf`/`log.Fatal` with `log/slog` in `cmd/intent-service/main.go`
- [x] T7.2 Initialize `slog.NewJSONHandler` in main.go; add `Logger *slog.Logger` to `FoundationOrchestrator`
- [x] T8.1 Add graceful shutdown (SIGINT/SIGTERM via `signal.NotifyContext`), 5s timeout, DB cleanup, double-signal force quit
- [x] T8.2 Wire `context.Context` to memory store methods that use select/channel ops

## Quality Gates

- [x] G1: `go fmt ./...` — zero diffs
- [x] G2: `go vet ./...` — zero warnings
- [x] G3: `go build ./...` — zero errors
- [x] G4: `go test ./... -count=1` — all tests pass
- [ ] G5: `go test -race ./...` — blocked: requires C toolchain (gcc) for race detector on Windows
- [x] G6: `go test -cover ./...` — all packages at or above baseline
