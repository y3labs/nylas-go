package nylas

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestAttachmentDeserialization(t *testing.T) {
	js := []byte(`{
		"content_type": "image/png",
		"filename": "pic.png",
		"grant_id": "41009df5-bf11-4c97-aa18-b285b5f2e386",
		"id": "185e56cb50e12e82",
		"is_inline": true,
		"size": 13068,
		"content_id": "<ce9b9547-9eeb-43b2-ac4e-58768bdf04e4>"
	}`)

	var a models.Attachment
	if err := json.Unmarshal(js, &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	assertStrPtr(t, a.ContentType, "image/png", "ContentType")
	assertStrPtr(t, a.Filename, "pic.png", "Filename")
	assertStrPtr(t, a.GrantID, "41009df5-bf11-4c97-aa18-b285b5f2e386", "GrantID")

	if a.ID != "185e56cb50e12e82" {
		t.Fatalf("ID = %q, want %q", a.ID, "185e56cb50e12e82")
	}
	if a.IsInline == nil || *a.IsInline != true {
		t.Fatalf("IsInline = %v, want true", a.IsInline)
	}
	if a.Size == nil || *a.Size != 13068 {
		t.Fatalf("Size = %v, want 13068", a.Size)
	}
	if a.ContentID == nil || *a.ContentID != "<ce9b9547-9eeb-43b2-ac4e-58768bdf04e4>" {
		t.Fatalf("ContentID = %v, want %q", a.ContentID, "<ce9b9547-9eeb-43b2-ac4e-58768bdf04e4>")
	}
}

func TestAttachmentsFind(t *testing.T) {
	grantID := "abc-123"
	attID := "attachment-123"
	msgID := "message-123"

	// Mock API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		wantPath := "/v3/grants/" + grantID + "/attachments/" + attID
		if r.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
		}
		q := r.URL.Query()
		if got := q.Get("message_id"); got != msgID {
			t.Fatalf("query message_id = %q, want %q", got, msgID)
		}

		resp := map[string]any{
			"request_id": "abc-req",
			"data": map[string]any{
				"id":           attID,
				"grant_id":     grantID,
				"filename":     "pic.png",
				"content_type": "image/png",
				"is_inline":    true,
				"size":         13068,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	params := &models.FindAttachmentQueryParams{MessageID: msgID}

	out, err := c.Attachments().Get(context.Background(), grantID, attID, params)
	if err != nil {
		t.Fatalf("Attachments.Find error: %v", err)
	}
	if out == nil {
		t.Fatal("nil response")
	}
	if out.RequestID != "abc-req" {
		t.Fatalf("RequestID = %q, want %q", out.RequestID, "abc-req")
	}
	if out.Data.ID != attID {
		t.Fatalf("Attachment.ID = %q, want %q", out.Data.ID, attID)
	}
	assertStrPtr(t, out.Data.GrantID, grantID, "Attachment.GrantID")
	assertStrPtr(t, out.Data.ContentType, "image/png", "ContentType")
}

func TestAttachmentsDownload_Stream(t *testing.T) {
	grantID := "abc-123"
	attID := "attachment-123"
	msgID := "message-123"

	data := []byte("mock data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		wantPath := "/v3/grants/" + grantID + "/attachments/" + attID + "/download"
		if r.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if r.URL.Query().Get("message_id") != msgID {
			t.Fatalf("message_id = %q, want %q", r.URL.Query().Get("message_id"), msgID)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(200)
		w.Write(data)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	params := &models.FindAttachmentQueryParams{MessageID: msgID}

	resp, err := c.Attachments().Download(context.Background(), grantID, attID, params)
	if err != nil {
		t.Fatalf("Attachments.Download error: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if string(b) != "mock data" {
		t.Fatalf("downloaded body = %q, want %q", string(b), "mock data")
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/octet-stream" {
		t.Fatalf("Content-Type = %q, want application/octet-stream", ct)
	}
}

func TestAttachmentsDownloadBytes(t *testing.T) {
	grantID := "abc-123"
	attID := "attachment-123"
	msgID := "message-123"

	data := []byte("mock data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		wantPath := "/v3/grants/" + grantID + "/attachments/" + attID + "/download"
		if r.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if r.URL.Query().Get("message_id") != msgID {
			t.Fatalf("message_id = %q, want %q", r.URL.Query().Get("message_id"), msgID)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(200)
		w.Write(data)
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	params := &models.FindAttachmentQueryParams{MessageID: msgID}

	b, err := c.Attachments().DownloadBytes(context.Background(), grantID, attID, params)
	if err != nil {
		t.Fatalf("Attachments.DownloadBytes error: %v", err)
	}
	if string(b) != "mock data" {
		t.Fatalf("download bytes = %q, want %q", string(b), "mock data")
	}
}

func TestAttachments_Get_ErrorPath(t *testing.T) {
	// Server returns 500 with an error envelope → DoJSON returns error, and Get should forward it.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/grants/ten-1/attachments/att-1" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		if r.URL.Query().Get("message_id") != "msg-1" {
			t.Fatalf("missing/incorrect message_id: %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"type":    "server_error",
				"message": "boom",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c: c}
	_, err := r.Get(context.Background(), "ten-1", "att-1", &models.FindAttachmentQueryParams{MessageID: "msg-1"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestAttachments_Get_RequestIDFallbackFromHeader(t *testing.T) {
	// 200 OK, JSON has no top-level request_id, but header has X-Request-Id → Get should copy it to Response.RequestID.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/grants/ten-2/attachments/att-2" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		if r.URL.Query().Get("message_id") != "msg-2" {
			t.Fatalf("missing/incorrect message_id: %q", r.URL.RawQuery)
		}
		w.Header().Set("X-Request-Id", "rid-xyz")
		_ = json.NewEncoder(w).Encode(map[string]any{
			// No "request_id" on purpose
			"data": map[string]any{
				"id": "att-2",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c: c}
	res, err := r.Get(context.Background(), "ten-2", "att-2", &models.FindAttachmentQueryParams{MessageID: "msg-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-xyz" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
	if res.Headers.Get("X-Request-Id") != "rid-xyz" {
		t.Fatalf("headers not propagated")
	}
}

func TestAttachments_DownloadBytes_Success(t *testing.T) {
	// Server returns 200 with body; we should get the bytes back.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/grants/ten-3/attachments/att-3/download" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		if r.URL.Query().Get("message_id") != "msg-3" {
			t.Fatalf("missing/incorrect message_id: %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte("hello-bytes"))
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c: c}
	b, err := r.DownloadBytes(context.Background(), "ten-3", "att-3", &models.FindAttachmentQueryParams{MessageID: "msg-3"})
	if err != nil {
		t.Fatalf("download bytes err: %v", err)
	}
	if string(b) != "hello-bytes" {
		t.Fatalf("unexpected body: %q", string(b))
	}
}

func TestAttachments_DownloadBytes_ErrorPath(t *testing.T) {
	// Non-2xx → DownloadBytes should parse error and return it.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/grants/ten-4/attachments/att-4/download" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		if r.URL.Query().Get("message_id") != "msg-4" {
			t.Fatalf("missing/incorrect message_id: %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-deny",
			"error": map[string]any{
				"type":    "forbidden",
				"message": "nope",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c: c}
	_, err := r.DownloadBytes(context.Background(), "ten-4", "att-4", &models.FindAttachmentQueryParams{MessageID: "msg-4"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if e, ok := IsAPIError(err); !ok || e.StatusCode != http.StatusForbidden {
		t.Fatalf("expected APIError 403, got %#v", err)
	}
}

func TestAttachments_Get_ErrorPath2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure the request looks right
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/grants/g1/attachments/a1" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
		}
		// Force a non-2xx so DoJSON -> parseAPIError -> err
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-bad",
			"error": map[string]any{
				"type":    "invalid",
				"message": "nope",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c}
	_, err := r.Get(context.Background(), "g1", "a1",
		&models.FindAttachmentQueryParams{MessageID: "mid-123"},
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !IsStatus(err, 400) {
		t.Fatalf("want status 400, got %v", err)
	}
}

func TestAttachments_Get_HeaderFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "message_id=mid-123" {
			t.Fatalf("expected query message_id=mid-123, got %q", r.URL.RawQuery)
		}
		w.Header().Set("x-request-id", "rid-from-header")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{ // no "request_id" on purpose
				"id": "att-1",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &AttachmentsResource{c}
	out, err := r.Get(context.Background(), "g1", "att-1",
		&models.FindAttachmentQueryParams{MessageID: "mid-123"},
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil || out.RequestID != "rid-from-header" {
		t.Fatalf("request id fallback failed: %#v", out)
	}
}
