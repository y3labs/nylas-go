package nylas

import (
	"io"
	"net/http"
)

// Response wraps a single object with metadata.
type Response[T any] struct {
	Data      T           `json:"data"`
	RequestID string      `json:"request_id"`
	Headers   http.Header `json:"-"`
}

// ListResponse wraps a list of objects with paging cursor.
type ListResponse[T any] struct {
	Data       []T         `json:"data"`
	RequestID  string      `json:"request_id"`
	NextCursor *string     `json:"next_cursor,omitempty"`
	Headers    http.Header `json:"-"`
}

// DeleteResponse Object returned from Nylas API
type DeleteResponse struct {
	RequestID string      `json:"request_id"`
	Headers   http.Header `json:"-"`
}

type FormField struct {
	Name        string // e.g., "message"
	Value       string // already-encoded bytes as string (e.g., JSON)
	ContentType string // optional; if empty, header omitted
}

// FormFile represents a file part with explicit filename and optional Content-Type.
type FormFile struct {
	Field       string    // e.g., "file0"
	Filename    string    // original filename for Content-Disposition
	ContentType string    // optional; if empty, header omitted (server may infer)
	Reader      io.Reader // file/content reader
}

// Helpers
func Ptr[T any](v T) *T          { return &v }
func StringPtr(s string) *string { return &s }
