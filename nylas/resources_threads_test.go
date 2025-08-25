package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestThreadDeserialization(t *testing.T) {
	js := `{
		"grant_id": "ca8f1733-6063-40cc-a2e3-ec7274abef11",
		"id": "7ml84jdmfnw20sq59f30hirhe",
		"object": "thread",
		"has_attachments": false,
		"has_drafts": false,
		"earliest_message_date": 1634149514,
		"latest_message_received_date": 1634832749,
		"latest_message_sent_date": 1635174399,
		"participants": [
			{"email": "daenerys.t@example.com", "name": "Daenerys Targaryen"}
		],
		"snippet": "jnlnnn --Sent with Nylas",
		"starred": false,
		"subject": "Dinner Wednesday?",
		"unread": false,
		"message_ids": ["njeb79kFFzli09", "998abue3mGH4sk"],
		"draft_ids": ["a809kmmoW90Dx"],
		"folders": ["8l6c4d11y1p4dm4fxj52whyr9", "d9zkcr2tljpu3m4qpj7l2hbr0"],
		"latest_draft_or_message": {
			"body": "Hello, I just sent a message using Nylas!",
			"cc": [{"name": "Arya Stark", "email": "arya.stark@example.com"}],
			"date": 1635355739,
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
			"id": "njeb79kFFzli09",
			"object": "message",
			"reply_to": [{"name": "Daenerys Targaryen", "email": "daenerys.t@example.com"}],
			"snippet": "Hello, I just sent a message using Nylas!",
			"starred": true,
			"subject": "Hello from Nylas!",
			"thread_id": "1t8tv3890q4vgmwq6pmdwm8qgsaer",
			"to": [{"name": "Jon Snow", "email": "j.snow@example.com"}],
			"unread": true
		}
	}`

	var th models.Thread
	if err := json.Unmarshal([]byte(js), &th); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if th.GrantID != "ca8f1733-6063-40cc-a2e3-ec7274abef11" {
		t.Fatalf("grant_id mismatch: %q", th.GrantID)
	}
	if th.ID != "7ml84jdmfnw20sq59f30hirhe" {
		t.Fatalf("id mismatch: %q", th.ID)
	}
	if th.Object != "thread" {
		t.Fatalf("object mismatch: %q", th.Object)
	}
	if th.HasAttachments == nil || *th.HasAttachments != false {
		t.Fatalf("has_attachments mismatch: %#v", th.HasAttachments)
	}
	if th.HasDrafts != false {
		t.Fatalf("has_drafts mismatch: %v", th.HasDrafts)
	}
	if th.EarliestMessageDate == nil || *th.EarliestMessageDate != 1634149514 {
		t.Fatalf("earliest_message_date mismatch: %#v", th.EarliestMessageDate)
	}
	if th.LatestMessageReceivedDate == nil || *th.LatestMessageReceivedDate != 1634832749 {
		t.Fatalf("latest_message_received_date mismatch: %#v", th.LatestMessageReceivedDate)
	}
	if th.LatestMessageSentDate == nil || *th.LatestMessageSentDate != 1635174399 {
		t.Fatalf("latest_message_sent_date mismatch: %#v", th.LatestMessageSentDate)
	}
	if len(th.Participants) != 1 || th.Participants[0].Email != "daenerys.t@example.com" || (th.Participants[0].Name == nil || *th.Participants[0].Name != "Daenerys Targaryen") {
		t.Fatalf("participants mismatch: %#v", th.Participants)
	}
	if th.Snippet == nil || *th.Snippet != "jnlnnn --Sent with Nylas" {
		t.Fatalf("snippet mismatch: %#v", th.Snippet)
	}
	if th.Starred != false {
		t.Fatalf("starred mismatch: %v", th.Starred)
	}
	if th.Subject == nil || *th.Subject != "Dinner Wednesday?" {
		t.Fatalf("subject mismatch: %#v", th.Subject)
	}
	if th.Unread != false {
		t.Fatalf("unread mismatch: %v", th.Unread)
	}
	if len(th.MessageIDs) != 2 || th.MessageIDs[0] != "njeb79kFFzli09" || th.MessageIDs[1] != "998abue3mGH4sk" {
		t.Fatalf("message_ids mismatch: %#v", th.MessageIDs)
	}
	if len(th.DraftIDs) != 1 || th.DraftIDs[0] != "a809kmmoW90Dx" {
		t.Fatalf("draft_ids mismatch: %#v", th.DraftIDs)
	}
	if len(th.Folders) != 2 {
		t.Fatalf("folders mismatch: %#v", th.Folders)
	}

	// latest_draft_or_message -> Message
	if th.LatestDraftOrMessage.Message == nil {
		t.Fatalf("expected latest_draft_or_message.message, got nil")
	}
	m := th.LatestDraftOrMessage.Message
	if m.Body == nil || *m.Body == "" || *m.Body != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("latest message body mismatch: %#v", m.Body)
	}
	if m.Date == nil || *m.Date != 1635355739 {
		t.Fatalf("latest message date mismatch: %#v", m.Date)
	}

	att := m.Attachments[0]

	if att.ContentType == nil || *att.ContentType != "text/calendar" {
		t.Fatalf("attachments content_type mismatch: %#v", att)
	}
	if att.Size == nil || *att.Size != 1708 {
		t.Fatalf("attachments size mismatch: %#v", att)
	}
	if len(m.Folders) != 2 {
		t.Fatalf("message.folders mismatch: %#v", m.Folders)
	}
	if len(m.From) != 1 || m.From[0].Email != "daenerys.t@example.com" {
		t.Fatalf("from mismatch: %#v", m.From)
	}
	if m.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("message.grant_id mismatch: %q", m.GrantID)
	}
	if m.ID != "njeb79kFFzli09" || m.Object != "message" {
		t.Fatalf("message id/object mismatch: id=%q object=%q", m.ID, m.Object)
	}
	if len(m.ReplyTo) != 1 || m.ReplyTo[0].Email != "daenerys.t@example.com" {
		t.Fatalf("reply_to mismatch: %#v", m.ReplyTo)
	}
	if m.Snippet == nil || *m.Snippet != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("message.snippet mismatch: %#v", m.Snippet)
	}
	if m.Starred == nil || *m.Starred != true {
		t.Fatalf("message.starred mismatch: %#v", m.Starred)
	}
	if m.Subject == nil || *m.Subject != "Hello from Nylas!" {
		t.Fatalf("message.subject mismatch: %#v", m.Subject)
	}
	if m.ThreadID == nil || *m.ThreadID != "1t8tv3890q4vgmwq6pmdwm8qgsaer" {
		t.Fatalf("message.thread_id mismatch: %#v", m.ThreadID)
	}
	if len(m.To) != 1 || m.To[0].Email != "j.snow@example.com" {
		t.Fatalf("to mismatch: %#v", m.To)
	}
	if m.Unread == nil || *m.Unread != true {
		t.Fatalf("message.unread mismatch: %#v", m.Unread)
	}
}

func TestListThreads(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       []map[string]any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Threads().List(context.Background(), "abc-123", nil); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListThreadsWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads")
		if got := r.URL.Query().Get("to"); got != "abc@gmail.com" {
			t.Fatalf("expected to=abc@gmail.com, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123", "data": []map[string]any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	to := "abc@gmail.com"
	_, err := c.Threads().List(context.Background(), "abc-123", &models.ListThreadsQueryParams{To: &to})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListThreadsWithSelectParam(t *testing.T) {
	wantSelect := "id,has_attachments,earliest_message_date,participants,snippet,unread,subject,message_ids,folders"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads")
		if got := r.URL.Query().Get("select"); got != wantSelect {
			t.Fatalf("expected select=%q, got %q", wantSelect, got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{
					"id":                         "thread-123",
					"has_attachments":            false,
					"earliest_message_date":      1634149514,
					"participants":               []map[string]any{{"email": "test@example.com", "name": "Test User"}},
					"snippet":                    "Test snippet",
					"unread":                     false,
					"subject":                    "Test subject",
					"message_ids":                []string{"msg-123"},
					"folders":                    []string{"folder-123"},
					"latest_message_received_at": nil,
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	selectStr := wantSelect
	res, err := c.Threads().List(context.Background(), "abc-123", &models.ListThreadsQueryParams{Select: &selectStr})
	if err != nil || res == nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListThreadsWithEarliestMessageDateParam(t *testing.T) {
	var tsVal int64 = 1672531200
	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads")
		if got := r.URL.Query().Get("earliest_message_date"); got != "1672531200" {
			t.Fatalf("expected earliest_message_date=1672531200, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{
					"id":                    "thread-123",
					"has_attachments":       false,
					"earliest_message_date": 1672617600,
					"participants":          []map[string]any{{"email": "test@example.com", "name": "Test User"}},
					"snippet":               "Test snippet",
					"unread":                false,
					"subject":               "Test subject",
					"message_ids":           []string{"msg-123"},
					"folders":               []string{"folder-123"},
				},
			},
		})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	_, err := c.Threads().List(context.Background(), "abc-123", &models.ListThreadsQueryParams{EarliestMessageDate: &tsVal})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListThreadsWithoutEarliestMessageDateInResponse(t *testing.T) {
	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{
					"id":          "thread-123",
					"grant_id":    "test-grant-id",
					"has_drafts":  false,
					"starred":     false,
					"unread":      false,
					"message_ids": []string{"msg-123"},
					"folders":     []string{"folder-123"},
					"latest_draft_or_message": map[string]any{
						"body":      "Test message body",
						"date":      1672617600,
						"from":      []map[string]any{{"name": "Test User", "email": "test@example.com"}},
						"grant_id":  "test-grant-id",
						"id":        "msg-123",
						"object":    "message",
						"subject":   "Test subject",
						"thread_id": "thread-123",
						"to":        []map[string]any{{"name": "Recipient", "email": "recipient@example.com"}},
						"unread":    false,
					},
					"has_attachments": false,
					"participants":    []map[string]any{{"email": "test@example.com", "name": "Test User"}},
					"snippet":         "Test snippet",
					"subject":         "Test subject",
				},
			},
		})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	res, err := c.Threads().List(context.Background(), "abc-123", nil)
	if err != nil || res == nil || len(res.Data) != 1 {
		t.Fatalf("List error or unexpected data: err=%v data=%#v", err, res)
	}
}

func TestFindThread(t *testing.T) {
	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/threads/thread-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req",
			"data":       map[string]any{"id": "thread-123", "grant_id": "abc-123", "has_drafts": false, "starred": false, "unread": false, "message_ids": []string{}, "folders": []string{}, "latest_draft_or_message": nil},
		})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	if _, err := c.Threads().Find(context.Background(), "abc-123", "thread-123"); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestFindThreadEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/v3/grants/abc-123/threads/" + encoded
		assertMethodPath(t, r, http.MethodGet, want)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": raw, "grant_id": "abc-123", "has_drafts": false, "starred": false, "unread": false, "message_ids": []string{}, "folders": []string{}, "latest_draft_or_message": nil}})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	if _, err := c.Threads().Find(context.Background(), "abc-123", raw); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestUpdateThread(t *testing.T) {
	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/threads/thread-123")

		var got models.UpdateThreadRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got.Starred == nil || *got.Starred != true {
			t.Fatalf("starred mismatch: %#v", got.Starred)
		}
		if got.Unread == nil || *got.Unread != false {
			t.Fatalf("unread mismatch: %#v", got.Unread)
		}
		if len(got.Folders) != 1 || got.Folders[0] != "folder-123" {
			t.Fatalf("folders mismatch: %#v", got.Folders)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req",
			"data":       map[string]any{"id": "thread-123"},
		})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	st, un := true, false
	body := models.UpdateThreadRequest{
		Starred: &st,
		Unread:  &un,
		Folders: []string{"folder-123"},
	}
	if _, err := c.Threads().Update(context.Background(), "abc-123", "thread-123", body); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestUpdateThreadEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/v3/grants/abc-123/threads/" + encoded
		assertMethodPath(t, r, http.MethodPut, want)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": raw}})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	st, un := true, false
	body := models.UpdateThreadRequest{Starred: &st, Unread: &un, Folders: []string{"folder-123"}}
	if _, err := c.Threads().Update(context.Background(), "abc-123", raw, body); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDestroyThread(t *testing.T) {
	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/threads/thread-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req"})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	if _, err := c.Threads().Destroy(context.Background(), "abc-123", "thread-123"); err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
}

func TestDestroyThreadEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	encoded := url.PathEscape(raw)

	tsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/v3/grants/abc-123/threads/" + encoded
		assertMethodPath(t, r, http.MethodDelete, want)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req"})
	}))
	defer tsrv.Close()

	c := newTestClient(tsrv.URL, "test-key")
	if _, err := c.Threads().Destroy(context.Background(), "abc-123", raw); err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
}
