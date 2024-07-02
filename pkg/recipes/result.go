package recipes

// Result represents the result of a recipe execution.
type Result struct {
	// Values represents the output values of the recipe.
	Values map[string]any `json:"values,omitempty"`
	// Secrets represents the output secrets of the recipe.
	Secrets map[string]any `json:"secrets,omitempty"`
	// Resources represents the output resources of the recipe.
	Resources []string `json:"resources,omitempty"`
}
