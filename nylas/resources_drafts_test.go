package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/y3labs/nylas-go/nylas/models"
)

func newTestClient(baseURL, apiKey string) *Client {
	return NewClient(
		apiKey,
		WithServerURL(baseURL),
		WithTimeout(5*time.Second),
		WithUserAgent("test-agent/0.0.1"),
	)
}

// --- Deserialization ---

func TestDraftDeserialization(t *testing.T) {
	var d models.Draft
	if err := json.Unmarshal([]byte(`{
		"body": "Hello, I just sent a message using Nylas!",
		"cc": [{"email": "arya.stark@example.com"}],
		"attachments": [
			{
				"content_type": "text/calendar",
				"id": "4kj2jrcoj9ve5j9yxqz5cuv98",
				"size": 1708
			}
		],
		"folders": ["8l6c4d11y1p4dm4fxj52whyr9", "d9zkcr2tljpu3m4qpj7l2hbr0"],
		"from": [{"name": "Daenerys Targaryen", "email": "daenerys.t@example.com"}],
		"grant_id": "41009df5-bf11-4c97-aa18-b285b5f2e386",
		"id": "5d3qmne77v32r8l4phyuksl2x",
		"object": "draft",
		"reply_to": [{"name": "Daenerys Targaryen", "email": "daenerys.t@example.com"}],
		"snippet": "Hello, I just sent a message using Nylas!",
		"starred": true,
		"subject": "Hello from Nylas!",
		"thread_id": "1t8tv3890q4vgmwq6pmdwm8qgsaer",
		"to": [{"name": "Jon Snow", "email": "j.snow@example.com"}],
		"date": 1705084742,
		"created_at": 1705084926
	}`), &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if d.Body == nil || *d.Body != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("Body = %q", strval(d.Body))
	}
	if len(d.Cc) != 1 || d.Cc[0].Email != "arya.stark@example.com" {
		t.Fatalf("Cc = %#v", d.Cc)
	}
	if len(d.Attachments) != 1 {
		t.Fatalf("Attachments len = %d", len(d.Attachments))
	}
	if strval(d.Attachments[0].ContentType) != "text/calendar" {
		t.Fatalf("Attachment.ContentType = %q", strval(d.Attachments[0].ContentType))
	}
	if d.Attachments[0].ID != "4kj2jrcoj9ve5j9yxqz5cuv98" {
		t.Fatalf("Attachment.ID = %q", d.Attachments[0].ID)
	}
	if d.Attachments[0].Size == nil || *d.Attachments[0].Size != 1708 {
		t.Fatalf("Attachment.Size = %#v", d.Attachments[0].Size)
	}
	if len(d.Folders) != 2 {
		t.Fatalf("Folders = %#v", d.Folders)
	}
	if len(d.From) != 1 || strval(d.From[0].Name) != "Daenerys Targaryen" || d.From[0].Email != "daenerys.t@example.com" {
		t.Fatalf("From = %#v", d.From)
	}
	if d.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("GrantID = %q", d.GrantID)
	}
	if d.ID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("ID = %q", d.ID)
	}
	if d.Object != "draft" {
		t.Fatalf("Object = %q", d.Object)
	}
	if d.ReplyTo == nil || len(d.ReplyTo) != 1 || d.ReplyTo[0].Email != "daenerys.t@example.com" {
		t.Fatalf("ReplyTo = %#v", d.ReplyTo)
	}
	if d.Snippet == nil || *d.Snippet != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("Snippet = %q", strval(d.Snippet))
	}
	if d.Starred == nil || !*d.Starred {
		t.Fatalf("Starred = %#v", d.Starred)
	}
	if d.Subject == nil || *d.Subject != "Hello from Nylas!" {
		t.Fatalf("Subject = %q", strval(d.Subject))
	}
	if d.ThreadID == nil || *d.ThreadID != "1t8tv3890q4vgmwq6pmdwm8qgsaer" {
		t.Fatalf("ThreadID = %q", strval(d.ThreadID))
	}
	if len(d.To) != 1 || d.To[0].Email != "j.snow@example.com" {
		t.Fatalf("To = %#v", d.To)
	}
	if d.Date == nil || *d.Date != 1705084742 {
		t.Fatalf("Date = %#v", d.Date)
	}
	if d.CreatedAt == nil || *d.CreatedAt != 1705084926 {
		t.Fatalf("CreatedAt = %#v", d.CreatedAt)
	}
}

// --- List ---

func TestListDrafts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/drafts")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Drafts().List(context.Background(), "abc-123", nil)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if out == nil || out.RequestID == "" {
		t.Fatalf("unexpected list response: %#v", out)
	}
}

func TestListDraftsWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/drafts")
		if got := r.URL.Query().Get("subject"); got != "Hello from Nylas!" {
			t.Fatalf("subject query = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	q := &models.ListDraftsQueryParams{Subject: strptr("Hello from Nylas!")}
	if _, err := c.Drafts().List(context.Background(), "abc-123", q); err != nil {
		t.Fatalf("List with query error: %v", err)
	}
}

// --- Find ---

func TestFindDraft(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/drafts/draft-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "draft-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Drafts().Find(context.Background(), "abc-123", "draft-123", nil)
	if err != nil {
		t.Fatalf("Find error: %v", err)
	}
	if out.Data.ID != "draft-123" {
		t.Fatalf("Find ID = %q", out.Data.ID)
	}
}

func TestFindDraftEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/drafts/"+encoded)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": raw},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Drafts().Find(context.Background(), "abc-123", raw, nil); err != nil {
		t.Fatalf("Find encoded error: %v", err)
	}
}

// --- Create ---

func TestCreateDraft(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/drafts")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["subject"] != "Hello from Nylas!" {
			t.Fatalf("subject = %v", body["subject"])
		}
		if body["body"] != "This is the body of my draft message." {
			t.Fatalf("body = %v", body["body"])
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "draft-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.CreateDraftRequest{
		Subject: strptr("Hello from Nylas!"),
		To:      []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
		Cc:      []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
		Body:    strptr("This is the body of my draft message."),
	}
	if _, err := c.Drafts().Create(context.Background(), "abc-123", req); err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func TestCreateDraftWithMetadata(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/drafts")

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		md, ok := body["metadata"].(map[string]any)
		if !ok || md["custom_field"] != "value" || md["another_field"] != float64(123) {
			t.Fatalf("metadata = %#v", body["metadata"])
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "draft-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.CreateDraftRequest{
		Subject:  strptr("Hello from Nylas!"),
		To:       []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
		Cc:       []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
		Body:     strptr("This is the body of my draft message."),
		Metadata: map[string]any{"custom_field": "value", "another_field": 123},
	}
	if _, err := c.Drafts().Create(context.Background(), "abc-123", req); err != nil {
		t.Fatalf("Create (metadata) error: %v", err)
	}
}

func TestCreateDraftWithAttachmentSmallAndLargeJSONOnly(t *testing.T) {
	// Our Go SDK currently always JSON-encodes attachments; both sizes go through the same path.
	tests := []int64{3, 3 * 1024 * 1024}

	for _, size := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/drafts")

			var body struct {
				Attachments []map[string]any `json:"attachments"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body.Attachments) != 1 {
				t.Fatalf("attachments len = %d", len(body.Attachments))
			}
			gotSize, _ := body.Attachments[0]["size"].(float64)
			if int64(gotSize) != size {
				t.Fatalf("attachment size = %v, want %d", gotSize, size)
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-123", "data": map[string]any{"id": "draft-123"}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		req := models.CreateDraftRequest{
			Subject: strptr("Hello from Nylas!"),
			To:      []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
			Cc:      []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
			Body:    strptr("This is the body of my draft message."),
			Attachments: []models.CreateAttachmentRequest{
				{
					Filename:      "file1.txt",
					ContentType:   "text/plain",
					ContentBase64: strptr("this is a file"),
					Size:          size,
				},
			},
		}
		if _, err := c.Drafts().Create(context.Background(), "abc-123", req); err != nil {
			t.Fatalf("Create (size=%d) error: %v", size, err)
		}
	}
}

// --- Update ---

func TestUpdateDraft(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/drafts/draft-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "draft-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.UpdateDraftRequest{
		Subject: strptr("Hello from Nylas!"),
		To:      []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
		Cc:      []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
		Body:    strptr("This is the body of my draft message."),
	}
	if _, err := c.Drafts().Update(context.Background(), "abc-123", "draft-123", req); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestUpdateDraftEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/drafts/"+encoded)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-123", "data": map[string]any{"id": raw}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.UpdateDraftRequest{
		Subject: strptr("Hello from Nylas!"),
		Body:    strptr("This is the body of my draft message."),
	}
	if _, err := c.Drafts().Update(context.Background(), "abc-123", raw, req); err != nil {
		t.Fatalf("Update encoded error: %v", err)
	}
}

func TestUpdateDraftWithAttachmentSmallAndLargeJSONOnly(t *testing.T) {
	tests := []int64{3, 3 * 1024 * 1024}

	for _, size := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/drafts/draft-123")

			var body struct {
				Attachments []map[string]any `json:"attachments"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(body.Attachments) != 1 {
				t.Fatalf("attachments len = %d", len(body.Attachments))
			}
			gotSize, _ := body.Attachments[0]["size"].(float64)
			if int64(gotSize) != size {
				t.Fatalf("attachment size = %v, want %d", gotSize, size)
			}

			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-123", "data": map[string]any{"id": "draft-123"}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		req := models.UpdateDraftRequest{
			Subject: strptr("Hello from Nylas!"),
			To:      []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
			Cc:      []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
			Body:    strptr("This is the body of my draft message."),
			Attachments: []models.CreateAttachmentRequest{
				{
					Filename:      "file1.txt",
					ContentType:   "text/plain",
					ContentBase64: strptr("this is a file"),
					Size:          size,
				},
			},
		}
		if _, err := c.Drafts().Update(context.Background(), "abc-123", "draft-123", req); err != nil {
			t.Fatalf("Update (size=%d) error: %v", size, err)
		}
	}
}

// --- Destroy ---

func TestDestroyDraft(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/drafts/draft-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-123"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Drafts().Destroy(context.Background(), "abc-123", "draft-123")
	if err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
	if out == nil || out.RequestID == "" {
		t.Fatalf("unexpected delete response: %#v", out)
	}
}

func TestDestroyDraftEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/drafts/"+encoded)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-123"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Drafts().Destroy(context.Background(), "abc-123", raw); err != nil {
		t.Fatalf("Destroy encoded error: %v", err)
	}
}

// --- Send ---

func TestSendDraft(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/drafts/draft-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "msg-1"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Drafts().Send(context.Background(), "abc-123", "draft-123"); err != nil {
		t.Fatalf("Send error: %v", err)
	}
}

func TestSendDraftEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/drafts/"+encoded)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       map[string]any{"id": "msg-1"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Drafts().Send(context.Background(), "abc-123", raw); err != nil {
		t.Fatalf("Send encoded error: %v", err)
	}
}
