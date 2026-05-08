package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"opita-sync-framework/internal/app/accessservice"
	"opita-sync-framework/internal/app/artifactservice"
	"opita-sync-framework/internal/app/devsurface"
	"opita-sync-framework/internal/app/intentservice"
	"opita-sync-framework/internal/app/operatorsurface"
	"opita-sync-framework/internal/app/pilotservice"
	"opita-sync-framework/internal/app/previewservice"
	"opita-sync-framework/internal/app/surfaceservice"
	"opita-sync-framework/internal/app/tenantservice"
	"opita-sync-framework/internal/engine/foundation"
	"opita-sync-framework/internal/engine/intent"
	"opita-sync-framework/internal/engine/policy"
	"opita-sync-framework/internal/engine/simulation"
	"opita-sync-framework/internal/platform/cerbos"
	"opita-sync-framework/internal/platform/filesystem"
	"opita-sync-framework/internal/platform/memory"
	pgplatform "opita-sync-framework/internal/platform/postgres"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	registryRoot := filepath.Join("definitions", "capabilities")
	registryResolver, err := filesystem.NewRegistryResolver(registryRoot)
	if err != nil {
		slog.Error("intent-service registry bootstrap failed", "error", err)
		os.Exit(1)
	}
	policyEngine := selectPolicyEngine()

	repo := memory.NewContractRepository()
	runtimeStore := memory.NewRuntimeService()
	eventLog := memory.NewEventLog()
	runStore := memory.NewFoundationRunStore()
	approvalStore := memory.NewApprovalStore()
	accessStore := memory.NewAccessStore()
	previewStore := memory.NewPreviewStore()
	intakeStore := memory.NewIntakeStore()
	proposalStore := memory.NewProposalStore()
	recoveryStore := memory.NewRecoveryStore()
	maintenanceStore := memory.NewMaintenanceStore()
	tenantStore := memory.NewTenantStore()
	artifactStore, err := filesystem.NewArtifactStore(filepath.Join("data", "artifacts"))
	if err != nil {
		slog.Error("artifact store bootstrap failed", "error", err)
		os.Exit(1)
	}
	retrievalStore := memory.NewRetrievalStore()

	var pgStore *pgplatform.Store

	if databaseURL := os.Getenv("OSF_DATABASE_URL"); databaseURL != "" {
		store, err := pgplatform.New(context.Background(), databaseURL)
		if err != nil {
			slog.Error("intent-service postgres bootstrap failed", "error", err)
			os.Exit(1)
		}
		pgStore = store
		repo = nil
		runtimeStore = nil
		eventLog = nil
		runStore = nil

		orchestrator, handler, previewHandler, surfaceHandler, operatorHandler, devHandler, artifactHandler, tenantHandler, accessHandler, pilotHandler := buildPostgresWiring(store, registryResolver)
		serveWithShutdown(orchestrator, handler, previewHandler, surfaceHandler, operatorHandler, devHandler, artifactHandler, tenantHandler, accessHandler, pilotHandler, pgStore)
		return
	}

	compiler := intent.NewCompiler(repo)
	orchestrator := &foundation.FoundationOrchestrator{
		Logger:    slog.Default(),
		Compiler:  compiler,
		Policy:    policyEngine,
		Runtime:   runtimeStore,
		Events:    eventLog,
		Registry:  registryResolver,
		Runs:      runStore,
		Approvals: approvalStore,
	}

	handler := intentservice.NewHandler(orchestrator, repo, runtimeStore, eventLog, runStore, registryResolver, approvalStore)
	previewHandler := previewservice.NewHandler(previewStore, simulation.NewService(policyEngine), eventLog)
	surfaceHandler := surfaceservice.NewHandler(intakeStore, proposalStore, eventLog)
	operatorHandler := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	devHandler := devsurface.NewHandler(runStore, maintenanceStore, eventLog)
	artifactHandler := artifactservice.NewHandler(artifactStore, retrievalStore)
	tenantHandler := tenantservice.NewHandler(tenantStore, eventLog)
	accessHandler := accessservice.NewHandler(accessStore, eventLog, approvalStore)
	pilotHandler := pilotservice.NewHandler(eventLog)
	serveWithShutdown(orchestrator, handler, previewHandler, surfaceHandler, operatorHandler, devHandler, artifactHandler, tenantHandler, accessHandler, pilotHandler, nil)
}

func buildPostgresWiring(store *pgplatform.Store, registryResolver *filesystem.RegistryResolver) (*foundation.FoundationOrchestrator, *intentservice.Handler, *previewservice.Handler, *surfaceservice.Handler, *operatorsurface.Handler, *devsurface.Handler, *artifactservice.Handler, *tenantservice.Handler, *accessservice.Handler, *pilotservice.Handler) {
	contractRepo := pgplatform.NewContractRepository(store)
	compiler := intent.NewCompiler(contractRepo)
	runtimeStore := pgplatform.NewRuntimeService(store)
	eventLog := pgplatform.NewEventLog(store)
	runStore := pgplatform.NewFoundationRunStore(store)
	approvalStore := pgplatform.NewApprovalStore(store)
	accessStore := pgplatform.NewAccessStore(store)
	previewStore := pgplatform.NewPreviewStore(store)
	intakeStore := pgplatform.NewIntakeStore(store)
	proposalStore := pgplatform.NewProposalStore(store)
	recoveryStore := pgplatform.NewRecoveryStore(store)
	maintenanceStore := pgplatform.NewMaintenanceStore(store)
	tenantStore := pgplatform.NewTenantStore(store)
	artifactStore, err := filesystem.NewArtifactStore(filepath.Join("data", "artifacts"))
	if err != nil {
		slog.Error("artifact store bootstrap failed", "error", err)
		os.Exit(1)
	}
	retrievalStore := memory.NewRetrievalStore()
	orchestrator := &foundation.FoundationOrchestrator{
		Logger:    slog.Default(),
		Compiler:  compiler,
		Policy:    selectPolicyEngine(),
		Runtime:   runtimeStore,
		Events:    eventLog,
		Registry:  registryResolver,
		Runs:      runStore,
		Approvals: approvalStore,
	}
	handler := intentservice.NewHandler(orchestrator, contractRepo, runtimeStore, eventLog, runStore, registryResolver, approvalStore)
	previewHandler := previewservice.NewHandler(previewStore, simulation.NewService(selectPolicyEngine()), eventLog)
	surfaceHandler := surfaceservice.NewHandler(intakeStore, proposalStore, eventLog)
	operatorHandler := operatorsurface.NewHandler(runtimeStore, eventLog, runStore, approvalStore, recoveryStore)
	devHandler := devsurface.NewHandler(runStore, maintenanceStore, eventLog)
	artifactHandler := artifactservice.NewHandler(artifactStore, retrievalStore)
	tenantHandler := tenantservice.NewHandler(tenantStore, eventLog)
	accessHandler := accessservice.NewHandler(accessStore, eventLog, approvalStore)
	pilotHandler := pilotservice.NewHandler(eventLog)
	return orchestrator, handler, previewHandler, surfaceHandler, operatorHandler, devHandler, artifactHandler, tenantHandler, accessHandler, pilotHandler
}

func serveWithShutdown(orchestrator *foundation.FoundationOrchestrator, handler *intentservice.Handler, previewHandler *previewservice.Handler, surfaceHandler *surfaceservice.Handler, operatorHandler *operatorsurface.Handler, devHandler *devsurface.Handler, artifactHandler *artifactservice.Handler, tenantHandler *tenantservice.Handler, accessHandler *accessservice.Handler, pilotHandler *pilotservice.Handler, dbStore *pgplatform.Store) {
	if err := orchestrator.Validate(); err != nil {
		slog.Error("intent-service wiring invalid", "error", err)
		os.Exit(1)
	}

	if err := intentservice.Warmup(context.Background(), orchestrator); err != nil {
		slog.Error("intent-service bootstrap failed", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler.Routes())
	mux.Handle("/v1/intake/", surfaceHandler.Routes())
	mux.Handle("/v1/proposals", surfaceHandler.Routes())
	mux.Handle("/v1/proposals/", surfaceHandler.Routes())
	mux.Handle("/v1/patchsets", surfaceHandler.Routes())
	mux.Handle("/v1/patchsets/", surfaceHandler.Routes())
	mux.Handle("/v1/workspaces/", surfaceHandler.Routes())
	mux.Handle("/v1/previews", previewHandler.Routes())
	mux.Handle("/v1/previews/", previewHandler.Routes())
	mux.Handle("/v1/simulations", previewHandler.Routes())
	mux.Handle("/v1/readable-previews/", previewHandler.Routes())
	mux.Handle("/v1/inspection/", operatorHandler.Routes())
	mux.Handle("/v1/operator/", operatorHandler.Routes())
	mux.Handle("/v1/recovery-actions", operatorHandler.Routes())
	mux.Handle("/v1/recovery-actions/", operatorHandler.Routes())
	mux.Handle("/v1/debug/", devHandler.Routes())
	mux.Handle("/v1/maintenance-actions", devHandler.Routes())
	mux.Handle("/v1/maintenance-actions/", devHandler.Routes())
	mux.Handle("/v1/artifacts", artifactHandler.Routes())
	mux.Handle("/v1/artifacts/", artifactHandler.Routes())
	mux.Handle("/v1/retrieval/search", artifactHandler.Routes())
	mux.Handle("/v1/tenants/", tenantHandler.Routes())
	mux.Handle("/v1/tenants/bootstrap", tenantHandler.Routes())
	mux.Handle("/v1/tenants-catalog/", tenantHandler.Routes())
	mux.Handle("/v1/tenants-connectors/", tenantHandler.Routes())
	mux.Handle("/v1/tenant-admin/", tenantHandler.Routes())
	mux.Handle("/v1/tenant-access/", accessHandler.Routes())
	mux.Handle("/v1/pilot/", pilotHandler.Routes())

	addr := ":8080"
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		addr = port
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("intent-service ready",
		"addr", server.Addr,
		"mode", detectMode(handler),
	)

	// Graceful shutdown on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()

		slog.Info("shutting down gracefully...")

		// Arm a second-signal handler for force-quit.
		go func() {
			secondCtx, secondStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer secondStop()
			<-secondCtx.Done()
			slog.Error("second signal received, force quitting")
			if dbStore != nil {
				dbStore.Close()
			}
			os.Exit(1)
		}()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shutdown error", "error", err)
		}

		if dbStore != nil {
			dbStore.Close()
			slog.Info("database connection closed")
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func detectMode(handler *intentservice.Handler) string {
	switch fmt.Sprintf("%T", handler.Runtime) {
	case "*postgres.RuntimeService":
		return "postgres"
	default:
		return "memory"
	}
}

func selectPolicyEngine() policy.PolicyEngine {
	if baseURL := strings.TrimSpace(os.Getenv("OSF_CERBOS_URL")); baseURL != "" {
		return cerbos.NewClient(baseURL)
	}
	return memory.NewPolicyEngine()
}
