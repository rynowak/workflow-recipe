package activities

import (
	"fmt"
	"log/slog"
	"time"

	daprworkflow "github.com/dapr/go-sdk/workflow"
)

func CallDeployKubernetesResources(ctx *daprworkflow.WorkflowContext, input DeployKubernetesResourcesInput) (DeployKubernetesResourcesOutput, error) {
	task := ctx.CallActivity(DeployKubernetesResources, daprworkflow.ActivityInput(input))

	output := DeployKubernetesResourcesOutput{}
	err := task.Await(&output)
	if err != nil {
		return DeployKubernetesResourcesOutput{}, err
	}

	return output, nil
}

type DeployKubernetesResourcesInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeployKubernetesResourcesOutput struct {
	Resources []string `json:"resources"`
	Host      string   `json:"host"`
	Port      int      `json:"port"`
}

func DeployKubernetesResources(ctx daprworkflow.ActivityContext) (any, error) {
	input := DeployKubernetesResourcesInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	// Pretend we are deploying resources...
	logger := slog.Default()
	logger.Info("Deploying Kubernetes Deployment")
	logger.Info("Deploying Kubernetes Service")
	logger.Info("Waiting for pods to be ready")
	time.Sleep(2 * time.Second)
	logger.Info("Pods are ready")

	return DeployKubernetesResourcesOutput{
		Host: fmt.Sprintf("%s.%s.svc.cluster.local", input.Namespace, input.Name),
		Port: 5432,
		Resources: []string{
			"/planes/kubernetes/local/namespaces/" + input.Namespace + "/providers/core/Service/" + input.Name,
			"/planes/kubernetes/local/namespaces/" + input.Namespace + "/providers/apps/Deployment/" + input.Name,
		},
	}, nil
}

func CallDeleteKubernetesResources(ctx *daprworkflow.WorkflowContext, input DeleteKubernetesResourcesInput) (DeleteKubernetesResourcesOutput, error) {
	task := ctx.CallActivity(DeleteKubernetesResources, daprworkflow.ActivityInput(input))

	output := DeleteKubernetesResourcesOutput{}
	err := task.Await(&output)
	if err != nil {
		return DeleteKubernetesResourcesOutput{}, err
	}

	return output, nil
}

type DeleteKubernetesResourcesInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeleteKubernetesResourcesOutput struct {
}

func DeleteKubernetesResources(ctx daprworkflow.ActivityContext) (any, error) {
	input := DeleteKubernetesResourcesInput{}
	err := ctx.GetInput(&input)
	if err != nil {
		return nil, err
	}

	// Pretend we are deleting resources...
	logger := slog.Default()
	logger.Info("Deleting Kubernetes Deployment")
	logger.Info("Deleting Kubernetes Service")

	return DeleteKubernetesResourcesOutput{}, nil
}
