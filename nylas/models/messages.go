package models

// ----- Enums -----

type MessageFields string

const (
	MessageFieldsStandard        MessageFields = "standard"
	MessageFieldsIncludeHeaders  MessageFields = "include_headers"
	MessageFieldsIncludeTracking MessageFields = "include_tracking_options"
	MessageFieldsRawMIME         MessageFields = "raw_mime"
)

// ----- Atomic types -----

type MessageHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TrackingOptions struct {
	Opens         *bool   `json:"opens,omitempty"`
	ThreadReplies *bool   `json:"thread_replies,omitempty"`
	Links         *bool   `json:"links,omitempty"`
	Label         *string `json:"label,omitempty"`
}

// EmailName already exists in your repo (models/events.go). Reuse that.

// ----- Message -----

type Message struct {
	GrantID     string           `json:"grant_id"`
	From        []EmailName      `json:"from,omitempty"`   // JSON key is "from" (not reserved in Go)
	Object      string           `json:"object,omitempty"` // "message"
	ID          string           `json:"id,omitempty"`
	Body        *string          `json:"body,omitempty"`
	ThreadID    *string          `json:"thread_id,omitempty"`
	Subject     *string          `json:"subject,omitempty"`
	Snippet     *string          `json:"snippet,omitempty"`
	To          []EmailName      `json:"to,omitempty"`
	Bcc         []EmailName      `json:"bcc,omitempty"`
	Cc          []EmailName      `json:"cc,omitempty"`
	ReplyTo     []EmailName      `json:"reply_to,omitempty"`
	Attachments []Attachment     `json:"attachments,omitempty"` // see models/attachments.go (add if missing)
	Folders     []string         `json:"folders,omitempty"`
	Headers     []MessageHeader  `json:"headers,omitempty"`
	Unread      *bool            `json:"unread,omitempty"`
	Starred     *bool            `json:"starred,omitempty"`
	CreatedAt   *int64           `json:"created_at,omitempty"`
	Date        *int64           `json:"date,omitempty"`
	ScheduleID  *string          `json:"schedule_id,omitempty"`
	SendAt      *int64           `json:"send_at,omitempty"`
	Metadata    map[string]any   `json:"metadata,omitempty"`
	Tracking    *TrackingOptions `json:"tracking_options,omitempty"`
	RawMIME     *string          `json:"raw_mime,omitempty"`
}

// ----- Query params -----

type ListMessagesQueryParams struct {
	// Base list
	Limit     *int    `query:"limit,omitempty"`
	PageToken *string `query:"page_token,omitempty"`

	Subject           *string        `query:"subject,omitempty"`
	AnyEmail          []string       `query:"any_email,omitempty"`
	FromFilter        []string       `query:"from,omitempty"`
	ToFilter          []string       `query:"to,omitempty"`
	CcFilter          []string       `query:"cc,omitempty"`
	BccFilter         []string       `query:"bcc,omitempty"`
	In                *string        `query:"in,omitempty"`
	Unread            *bool          `query:"unread,omitempty"`
	Starred           *bool          `query:"starred,omitempty"`
	ThreadID          *string        `query:"thread_id,omitempty"`
	ReceivedBefore    *int64         `query:"received_before,omitempty"`
	ReceivedAfter     *int64         `query:"received_after,omitempty"`
	HasAttachment     *bool          `query:"has_attachment,omitempty"`
	Fields            *MessageFields `query:"fields,omitempty"`
	SearchQueryNative *string        `query:"search_query_native,omitempty"`
	Select            *string        `query:"select,omitempty"`
}

type FindMessageQueryParams struct {
	Fields *MessageFields `query:"fields,omitempty"`
	Select *string        `query:"select,omitempty"`
}

// ----- Update -----

type UpdateMessageRequest struct {
	Unread   *bool          `json:"unread,omitempty"`
	Starred  *bool          `json:"starred,omitempty"`
	Folders  []string       `json:"folders,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ----- Schedules -----

type ScheduledMessageStatus struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ScheduledMessage struct {
	ScheduleID string                 `json:"schedule_id"`
	Status     ScheduledMessageStatus `json:"status"`
	CloseTime  *int64                 `json:"close_time,omitempty"`
}

type StopScheduledMessageResponse struct {
	Message string `json:"message"`
}

// ----- Clean Messages -----

type CleanMessagesRequest struct {
	MessageID               []string `json:"message_id"`
	IgnoreLinks             *bool    `json:"ignore_links,omitempty"`
	IgnoreImages            *bool    `json:"ignore_images,omitempty"`
	ImagesAsMarkdown        *bool    `json:"images_as_markdown,omitempty"`
	IgnoreTables            *bool    `json:"ignore_tables,omitempty"`
	RemoveConclusionPhrases *bool    `json:"remove_conclusion_phrases,omitempty"`
}

type CleanMessagesResponse struct {
	Message
	Conversation string `json:"conversation"`
}

type SendMessageTrackingOptions struct {
	Label         *string `json:"label,omitempty"`
	Links         *bool   `json:"links,omitempty"`
	Opens         *bool   `json:"opens,omitempty"`
	ThreadReplies *bool   `json:"thread_replies,omitempty"`
}

// SendMessageRequest mirrors the Python SDK "send" payload (JSON mode).
// For large files, use the multipart method provided below instead.
// Sending via /messages/send uses a body similar to CreateDraftRequest,
// plus From and UseDraft. We'll reuse for messages resource.
type SendMessageRequest struct {
	CreateDraftRequest

	From     []EmailName `json:"from,omitempty"`
	UseDraft *bool       `json:"use_draft,omitempty"`
}
