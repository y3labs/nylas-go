package models

// Attachment mirrors nylas.models.attachments.Attachment
type Attachment struct {
	ID                 string  `json:"id"`
	GrantID            *string `json:"grant_id,omitempty"`
	Filename           *string `json:"filename,omitempty"`
	ContentType        *string `json:"content_type,omitempty"`
	Size               *int64  `json:"size,omitempty"`
	ContentID          *string `json:"content_id,omitempty"`
	ContentDisposition *string `json:"content_disposition,omitempty"`
	IsInline           *bool   `json:"is_inline,omitempty"`
}

// CreateAttachmentRequest mirrors the Python TypedDict.
// Note: Python allows Union[str, BinaryIO]. In Go, if you want to upload raw bytes,
// prefer the multipart Upload method in the resource. This struct is here for parity
// and for potential JSON/base64 workflows, using ContentBase64.
type CreateAttachmentRequest struct {
	Filename           string  `json:"filename"`
	ContentType        string  `json:"content_type"`
	ContentBase64      *string `json:"content,omitempty"` // Base64-encoded content
	Size               int64   `json:"size"`
	ContentID          *string `json:"content_id,omitempty"`
	ContentDisposition *string `json:"content_disposition,omitempty"`
	IsInline           *bool   `json:"is_inline,omitempty"`
}

// FindAttachmentQueryParams mirrors the Python TypedDict (required in Python).
// The query param is sent as ?message_id=...
type FindAttachmentQueryParams struct {
	MessageID string `json:"message_id" url:"message_id"`
}

// UploadAttachmentResponse is the JSON object Nylas returns after multipart upload.
/*
type UploadAttachmentResponse struct {
	RequestID  string     `json:"request_id"`
	Attachment Attachment `json:"attachment"`
}
*/
