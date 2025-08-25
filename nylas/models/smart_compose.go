package models

// ComposeMessageRequest represents the Smart Compose prompt payload.
type ComposeMessageRequest struct {
	Prompt string `json:"prompt"`
}

// ComposeMessageResponse is the generated suggestion from Smart Compose.
type ComposeMessageResponse struct {
	Suggestion string `json:"suggestion"`
}
