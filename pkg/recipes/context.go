package recipes

import (
	"log/slog"

	"github.com/go-openapi/jsonpointer"
)

// Context represents the context information which accesses portable resource properties. Recipe template authors
// can leverage the RecipeContext parameter to access portable resource properties to generate name and properties
// that are unique for the portable resource calling the recipe.
type Context struct {
	// Resource represents the resource information of the deploying recipe resource.
	Resource Resource `json:"resource,omitempty"`
	// Application represents environment resource information.
	Application ResourceInfo `json:"application,omitempty"`
	// Environment represents environment resource information.
	Environment ResourceInfo `json:"environment,omitempty"`
	// Runtime represents Kubernetes Runtime configuration.
	Runtime RuntimeConfiguration `json:"runtime,omitempty"`
	// Azure represents Azure provider scope.
	Azure *ProviderAzure `json:"azure,omitempty"`
	// AWS represents AWS provider scope.
	AWS *ProviderAWS `json:"aws,omitempty"`
}

func (c *Context) LogAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("resource.id", c.Resource.ID),
		slog.String("application.id", c.Application.ID),
		slog.String("environment.id", c.Application.Name),
	}
}

// Resource contains the information needed to deploy a recipe.
// In the case the resource is a portable resource, it represents the resource's id, name and type.
type Resource struct {
	// ResourceInfo represents name and id of the resource
	ResourceInfo
	// Type represents the resource type, this will be a namespace/type combo. Ex. Applications.Core/Environment
	Type string `json:"type"`

	// Properties are the properties of the resource.
	Properties map[string]any `json:"properties,omitempty"`
}

func (r *Resource) GetStringValue(key string) (string, bool) {
	ptr, err := jsonpointer.New(key)
	if err != nil {
		return "", false
	}

	value, _, err := ptr.Get(r.Properties)
	if err != nil {
		return "", false
	}

	return value.(string), true
}

// ResourceInfo represents name and id of the resource
type ResourceInfo struct {
	// Name represents the resource name.
	Name string `json:"name"`
	// ID represents fully qualified resource id.
	ID string `json:"id"`
}

// ProviderAzure contains Azure provider scope for recipe context.
type ProviderAzure struct {
	// ResourceGroup represents the resource group information.
	ResourceGroup AzureResourceGroup `json:"resourceGroup,omitempty"`
	// Subscription represents the subscription information.
	Subscription AzureSubscription `json:"subscription,omitempty"`
}

// AzureResourceGroup contains Azure Resource Group provider information.
type AzureResourceGroup struct {
	// Name represents the resource name.
	Name string `json:"name"`
	// ID represents fully qualified resource group name.
	ID string `json:"id"`
}

// AzureSubscription contains Azure Subscription provider information.
type AzureSubscription struct {
	// SubscriptionID represents the id of subscription.
	SubscriptionID string `json:"subscriptionId"`
	// ID represents fully qualified subscription id.
	ID string `json:"id"`
}

// ProviderAWS contains AWS Account provider scope for recipe context.
type ProviderAWS struct {
	// Region represents the region of the AWS account.
	Region string `json:"region"`
	// Account represents the account id of the AWS account.
	Account string `json:"account"`
}

// RuntimeConfiguration represents Kubernetes Runtime configuration for the environment.
type RuntimeConfiguration struct {
	Kubernetes *KubernetesRuntime `json:"kubernetes,omitempty"`
}

// KubernetesRuntime represents application and environment namespaces.
type KubernetesRuntime struct {
	// Namespace is set to the application namespace when the portable resource is application-scoped, and set to the environment namespace when it is environment scoped
	Namespace string `json:"namespace,omitempty"`
	// EnvironmentNamespace is set to environment namespace.
	EnvironmentNamespace string `json:"environmentNamespace"`
}
