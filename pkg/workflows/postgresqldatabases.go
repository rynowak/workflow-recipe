package workflows

import (
	"fmt"
	"log/slog"

	daprworkflow "github.com/dapr/go-sdk/workflow"
	"github.com/rynowak/workflow-recipe/pkg/activities"
	"github.com/rynowak/workflow-recipe/pkg/recipes"
)

func PostgresSQLDatabasesPut(ctx *daprworkflow.WorkflowContext) (any, error) {
	request := recipes.Context{}
	err := ctx.GetInput(&request)
	if err != nil {
		return nil, err
	}

	logger := slog.Default()
	if ctx.IsReplaying() {
		logger.Info("Resuming PostgresSQL database creation/update")
	} else {
		logger.Info("Creating/Updating PostgresSQL database")
	}

	deployed, err := activities.CallDeployKubernetesResources(ctx, activities.DeployKubernetesResourcesInput{
		Namespace: request.Runtime.Kubernetes.Namespace,
		Name:      request.Resource.Name,
	})
	if err != nil {
		return nil, err
	}

	credentials, err := activities.CallCreatePostgresUser(ctx, activities.CreatePostgresUserInput{})
	if err != nil {
		return nil, err
	}

	database, err := activities.CallCreatePostgresDatabase(ctx, activities.CreatePostgresDatabaseInput{
		Username:       credentials.Username,
		Password:       credentials.Password,
		DatabasePrefix: request.Resource.Name,
	})
	if err != nil {
		return nil, err
	}

	// Return data to Radius
	result := recipes.Result{
		Values: map[string]any{
			"host":     deployed.Host,
			"port":     deployed.Port,
			"username": credentials.Username,
			"database": database.Database,
		},
		Secrets: map[string]any{
			"password": credentials.Password,
			"uri":      fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", credentials.Username, credentials.Password, deployed.Host, deployed.Port, database.Database),
		},
		Resources: deployed.Resources,
	}

	logger.Info("Done creating/updating PostgresSQL database")
	return result, nil
}

func PostgresSQLDatabasesDelete(ctx *daprworkflow.WorkflowContext) (any, error) {
	request := recipes.Context{}
	err := ctx.GetInput(&request)
	if err != nil {
		return nil, err
	}

	logger := slog.Default()
	if ctx.IsReplaying() {
		logger.Info("Resuming PostgresSQL database deletion")
	} else {
		logger.Info("Deleting PostgresSQL database")
	}

	database, ok := request.Resource.GetStringValue("/status/binding/database")
	if !ok {
		_, err = activities.CallDeletePostgresDatabase(ctx, activities.DeletePostgresDatabaseInput{
			Database:     database,
			CreateBackup: true,
		})
		if err != nil {
			return nil, err
		}
	}

	username, ok := request.Resource.GetStringValue("/status/binding/username")
	if !ok {
		_, err = activities.CallDeletePostgresUser(ctx, activities.DeletePostgresUserInput{
			Username: username,
		})
		if err != nil {
			return nil, err
		}
	}

	_, err = activities.CallDeleteKubernetesResources(ctx, activities.DeleteKubernetesResourcesInput{
		Namespace: request.Runtime.Kubernetes.Namespace,
		Name:      request.Resource.Name,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("Done deleting PostgresSQL database")
	return struct{}{}, nil
}
