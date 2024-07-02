package activities

import (
	"fmt"
	"log/slog"
	"time"

	daprworkflow "github.com/dapr/go-sdk/workflow"
	"github.com/google/uuid"
)

func CallCreatePostgresUser(ctx *daprworkflow.WorkflowContext, input CreatePostgresUserInput) (CreatePostgresUserOutput, error) {
	task := ctx.CallActivity(CreatePostgresUser, daprworkflow.ActivityInput(input))

	output := CreatePostgresUserOutput{}
	err := task.Await(&output)
	if err != nil {
		return CreatePostgresUserOutput{}, err
	}

	return output, nil
}

type CreatePostgresUserInput struct {
}

type CreatePostgresUserOutput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreatePostgresUser(ctx daprworkflow.ActivityContext) (any, error) {
	input := CreatePostgresUserInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	// Pretend we are very secure...
	username := "pguser"
	password := uuid.NewString()

	logger := slog.Default()
	logger.Info("Generating new postgres user", slog.String("username", "pguser"))
	logger.Info("Generating really really secure password", slog.String("password", "********")) // Haha, just kidding

	return CreatePostgresUserOutput{
		Username: username,
		Password: password,
	}, nil
}

func CallDeletePostgresUser(ctx *daprworkflow.WorkflowContext, input DeletePostgresUserInput) (DeletePostgresUserOutput, error) {
	task := ctx.CallActivity(DeletePostgresUser, daprworkflow.ActivityInput(input))

	output := DeletePostgresUserOutput{}
	err := task.Await(&output)
	if err != nil {
		return DeletePostgresUserOutput{}, err
	}

	return output, nil
}

type DeletePostgresUserInput struct {
	Username string `json:"username"`
}

type DeletePostgresUserOutput struct {
}

func DeletePostgresUser(ctx daprworkflow.ActivityContext) (any, error) {
	input := DeletePostgresUserInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	logger := slog.Default()
	logger.Info("Deleting postgres user", slog.String("username", input.Username))

	return DeletePostgresUserOutput{}, nil
}

func CallCreatePostgresDatabase(ctx *daprworkflow.WorkflowContext, input CreatePostgresDatabaseInput) (CreatePostgresDatabaseOutput, error) {
	task := ctx.CallActivity(CreatePostgresDatabase, daprworkflow.ActivityInput(input))

	output := CreatePostgresDatabaseOutput{}
	err := task.Await(&output)
	if err != nil {
		return CreatePostgresDatabaseOutput{}, err
	}

	return output, nil
}

type CreatePostgresDatabaseInput struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	DatabasePrefix string `json:"databasePrefix"`
}

type CreatePostgresDatabaseOutput struct {
	Database string `json:"database"`
}

func CreatePostgresDatabase(ctx daprworkflow.ActivityContext) (any, error) {
	input := CreatePostgresDatabaseInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	// Pretend we are using a real database server...
	database := fmt.Sprintf("%s_%s", input.DatabasePrefix, uuid.NewString())

	logger := slog.Default()
	logger.Info("Creating new database", slog.String("database", database))
	logger.Info("Granting user permission", slog.String("username", input.Username))

	return CreatePostgresDatabaseOutput{
		Database: database,
	}, nil
}

func CallDeletePostgresDatabase(ctx *daprworkflow.WorkflowContext, input DeletePostgresDatabaseInput) (DeletePostgresDatabaseOutput, error) {
	task := ctx.CallActivity(DeletePostgresDatabase, daprworkflow.ActivityInput(input))

	output := DeletePostgresDatabaseOutput{}
	err := task.Await(&output)
	if err != nil {
		return DeletePostgresDatabaseOutput{}, err
	}

	return output, nil
}

type DeletePostgresDatabaseInput struct {
	Database     string `json:"database"`
	CreateBackup bool   `json:"createBackup"`
}

type DeletePostgresDatabaseOutput struct {
}

func DeletePostgresDatabase(ctx daprworkflow.ActivityContext) (any, error) {
	input := DeletePostgresDatabaseInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	// Pretend we are using a real database server...

	logger := slog.Default()
	logger.Info("Creating a backup", slog.String("database", input.Database))
	time.Sleep(5 * time.Second)
	logger.Info("Deleting database", slog.String("database", input.Database))

	return DeletePostgresDatabaseOutput{}, nil
}
