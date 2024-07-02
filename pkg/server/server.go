package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	daprclient "github.com/dapr/go-sdk/client"
	daprworkflow "github.com/dapr/go-sdk/workflow"
	"github.com/microsoft/durabletask-go/api"
)

const (
	Address = ":7999"
)

func Start(ctx context.Context, services map[string]context.CancelFunc, dapr daprclient.Client) error {
	workflowClient, err := daprworkflow.NewClient(daprworkflow.WithDaprClient(dapr))
	if err != nil {
		return fmt.Errorf("error creating Dapr workflow client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		mustWriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /workflows/{id}", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(ctx, "Fetching workflow metadata", slog.String("id", r.PathValue("id")))

		id := r.PathValue("id")
		metadata, err := workflowClient.FetchWorkflowMetadata(r.Context(), id, daprworkflow.WithFetchPayloads(true))
		if err != nil {
			mustWriteError(w, http.StatusInternalServerError, "Internal", err)
			return
		}

		mustWriteJSON(w, http.StatusOK, metadata)
	})

	mux.HandleFunc("PUT /workflows", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		request := WorkflowRequest{}
		err := decoder.Decode(&request)
		if err != nil {
			mustWriteError(w, http.StatusBadRequest, "Invalid", err)
			return
		}

		if request.Name == "" || len(request.Input) == 0 {
			mustWriteError(w, http.StatusBadRequest, "Invalid", errors.New("name and input are required"))
			return
		}

		slog.InfoContext(ctx, "Starting new workflow", slog.String("id", r.PathValue("id")), slog.String("name", request.Name))

		opts := []api.NewOrchestrationOptions{}
		opts = append(opts, daprworkflow.WithRawInput(string(request.Input)))
		if request.ID != "" {
			opts = append(opts, daprworkflow.WithInstanceID(request.ID))
		}

		result, err := workflowClient.ScheduleNewWorkflow(r.Context(), request.Name, opts...)
		if err != nil {
			mustWriteError(w, http.StatusInternalServerError, "Internal", err)
			return
		}

		slog.InfoContext(ctx, "Workflow started", slog.String("id", result), slog.String("name", request.Name))
		mustWriteJSON(w, http.StatusCreated, map[string]any{"id": result})
	})

	server := &http.Server{
		Addr:    Address,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	listener, err := net.Listen("tcp", Address)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	slog.InfoContext(ctx, "Server is listening", slog.String("address", listener.Addr().String()))

	go func() {
		err := server.Serve(listener)
		if errors.Is(http.ErrServerClosed, err) {
			slog.InfoContext(ctx, "server shutdown gracefully")
		} else {
			slog.ErrorContext(ctx, "server error", slog.Any("error", err))
		}
	}()

	services["http"] = func() {
		err := server.Shutdown(ctx)
		slog.ErrorContext(ctx, "server shutdown error", slog.Any("error", err))
	}

	return nil
}

func mustWriteJSON(w http.ResponseWriter, code int, v any) {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		mustWriteError(w, http.StatusInternalServerError, "Internal", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(bs)
}

func mustWriteError(w http.ResponseWriter, statusCode int, errorCode string, err error) {
	e := ErrorResponse{
		Error: ErrorDetails{
			Code:    errorCode,
			Message: err.Error(),
		},
	}

	// This should never fail.
	bs, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(bs)
}
