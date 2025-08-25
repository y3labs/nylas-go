package models

// Reuse TrackingOptions from messages.go

type CustomHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Draft mirrors Message but sets object "draft"
type Draft struct {
	Message
	// Python sets object="draft". Optional to read; server may include it.
}

// --- Create/Update Draft requests ---

// CreateDraftRequest mirrors Python (all optional)
type CreateDraftRequest struct {
	Body             *string                   `json:"body,omitempty"`
	Subject          *string                   `json:"subject,omitempty"`
	To               []EmailName               `json:"to,omitempty"`
	Bcc              []EmailName               `json:"bcc,omitempty"`
	Cc               []EmailName               `json:"cc,omitempty"`
	ReplyTo          []EmailName               `json:"reply_to,omitempty"`
	Attachments      []CreateAttachmentRequest `json:"attachments,omitempty"`
	Starred          *bool                     `json:"starred,omitempty"`
	SendAt           *int64                    `json:"send_at,omitempty"`
	ReplyToMessageID *string                   `json:"reply_to_message_id,omitempty"`
	TrackingOptions  *TrackingOptions          `json:"tracking_options,omitempty"`
	CustomHeaders    []CustomHeader            `json:"custom_headers,omitempty"`
	Metadata         map[string]any            `json:"metadata,omitempty"`
}

type UpdateDraftRequest = CreateDraftRequest

// --- Query params ---

type ListDraftsQueryParams struct {
	// Base list
	Limit     *int    `url:"limit,omitempty"`
	PageToken *string `url:"page_token,omitempty"`

	Subject       *string  `url:"subject,omitempty"`
	AnyEmail      []string `url:"any_email,omitempty"`
	FromFilter    []string `url:"from,omitempty"`
	ToFilter      []string `url:"to,omitempty"`
	CcFilter      []string `url:"cc,omitempty"`
	BccFilter     []string `url:"bcc,omitempty"`
	In            []string `url:"in,omitempty"` // folder/label IDs
	Unread        *bool    `url:"unread,omitempty"`
	Starred       *bool    `url:"starred,omitempty"`
	ThreadID      *string  `url:"thread_id,omitempty"`
	HasAttachment *bool    `url:"has_attachment,omitempty"`
	Select        *string  `url:"select,omitempty"`
}

type FindDraftQueryParams struct {
	Select *string `url:"select,omitempty"`
}
