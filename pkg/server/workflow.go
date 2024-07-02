package server

import "encoding/json"

type WorkflowRequest struct {
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
	ID    string          `json:"id,omitempty"`
}
