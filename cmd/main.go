package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	daprclient "github.com/dapr/go-sdk/client"
	daprworkflow "github.com/dapr/go-sdk/workflow"
	"github.com/rynowak/workflow-recipe/pkg/activities"
	"github.com/rynowak/workflow-recipe/pkg/server"
	"github.com/rynowak/workflow-recipe/pkg/workflows"
)

func main() {
	ctx := context.Background()

	services := map[string]context.CancelFunc{}
	ctx, cancel := registerShutdown(ctx, services)
	defer cancel()
	go func() {

	}()

	err := start(ctx, services)
	if err != nil {
		slog.ErrorContext(ctx, "Error starting services", slog.Any("error", err))
		os.Exit(1)
		return // unreachable
	}

	slog.InfoContext(ctx, "Server started: Press CTRL+C to stop")
	<-ctx.Done()
	cancel()
}

func registerShutdown(ctx context.Context, services map[string]context.CancelFunc) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	return ctx, func() {
		cancel()
		for name, cancel := range services {
			slog.InfoContext(ctx, "Shutting down", slog.String("service", name))
			cancel()
		}
	}
}

func start(ctx context.Context, services map[string]context.CancelFunc) error {
	slog.InfoContext(ctx, "Connecting to Dapr")

	dapr, err := daprclient.NewClient()
	if err != nil {
		return fmt.Errorf("error creating Dapr client: %v", err)
	}
	services["daprclient"] = dapr.Close

	err = registerWorkflows(ctx, dapr)
	if err != nil {
		return fmt.Errorf("error initializing workflows: %v", err)
	}

	err = server.Start(ctx, services, dapr)
	if err != nil {
		return fmt.Errorf("error starting HTTP server: %v", err)
	}

	return nil
}

func registerWorkflows(ctx context.Context, dapr daprclient.Client) error {
	worker, err := daprworkflow.NewWorker(daprworkflow.WorkerWithDaprClient(dapr))
	if err != nil {
		return fmt.Errorf("error creating Dapr workflow worker: %w", err)
	}

	// TODO: register workflows and activities.

	err = worker.RegisterWorkflow(workflows.PostgresSQLDatabasesPut)
	if err != nil {
		return fmt.Errorf("error registering workflow: %w", err)
	}

	err = worker.RegisterWorkflow(workflows.PostgresSQLDatabasesDelete)
	if err != nil {
		return fmt.Errorf("error registering workflow: %w", err)
	}

	err = worker.RegisterActivity(activities.DeployKubernetesResources)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.RegisterActivity(activities.DeleteKubernetesResources)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.RegisterActivity(activities.CreatePostgresUser)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.RegisterActivity(activities.DeletePostgresUser)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.RegisterActivity(activities.CreatePostgresDatabase)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.RegisterActivity(activities.DeletePostgresDatabase)
	if err != nil {
		return fmt.Errorf("error registering activity: %w", err)
	}

	err = worker.Start()
	if err != nil {
		return fmt.Errorf("error starting Dapr workflow worker: %w", err)
	}

	slog.InfoContext(ctx, "Dapr workflow worker started")
	return nil
}
