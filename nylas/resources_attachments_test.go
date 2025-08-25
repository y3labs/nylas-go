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
