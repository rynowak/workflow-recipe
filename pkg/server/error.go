package server

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code           string                `json:"code"`
	Message        string                `json:"message"`
	Target         string                `json:"target,omitempty"`
	AdditionalInfo []ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
	Details        []ErrorDetails        `json:"details,omitempty"`
}

type ErrorAdditionalInfo struct {
	Type string         `json:"type"`
	Info map[string]any `json:"info"`
}
