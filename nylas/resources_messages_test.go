package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestMessagesSmartComposeProperty(t *testing.T) {
	c := newTestClient("http://example.com", "test-key")
	if c.Messages().SmartCompose() == nil {
		t.Fatalf("SmartCompose() returned nil")
	}
}

func TestMessageDeserialization(t *testing.T) {
	js := `{
		"body": "Hello, I just sent a message using Nylas!",
		"cc": [{"name": "Arya Stark", "email": "arya.stark@example.com"}],
		"date": 1635355739,
		"attachments": [{"content_type":"text/calendar","id":"4kj2jrcoj9ve5j9yxqz5cuv98","size":1708}],
		"folders": ["8l6c4d11y1p4dm4fxj52whyr9","d9zkcr2tljpu3m4qpj7l2hbr0"],
		"from": [{"name":"Daenerys Targaryen","email":"daenerys.t@example.com"}],
		"grant_id": "41009df5-bf11-4c97-aa18-b285b5f2e386",
		"id": "5d3qmne77v32r8l4phyuksl2x",
		"object": "message",
		"reply_to": [{"name":"Daenerys Targaryen","email":"daenerys.t@example.com"}],
		"snippet": "Hello, I just sent a message using Nylas!",
		"starred": true,
		"subject": "Hello from Nylas!",
		"thread_id": "1t8tv3890q4vgmwq6pmdwm8qgsaer",
		"to": [{"name":"Jon Snow","email":"j.snow@example.com"}],
		"unread": true,
		"metadata": {"custom_field":"value","another_field":123}
	}`
	var m models.Message
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.Body == nil || *m.Body != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("body = %#v", m.Body)
	}
	if len(m.Cc) != 1 || m.Cc[0].Email != "arya.stark@example.com" {
		t.Fatalf("cc = %#v", m.Cc)
	}
	if m.Date == nil || *m.Date != 1635355739 {
		t.Fatalf("date = %#v", m.Date)
	}
	if len(m.Attachments) != 1 || m.Attachments[0].ContentType == nil || *m.Attachments[0].ContentType != "text/calendar" {
		t.Fatalf("attachments = %#v", m.Attachments)
	}
	if len(m.Folders) != 2 || m.Folders[0] != "8l6c4d11y1p4dm4fxj52whyr9" {
		t.Fatalf("folders = %#v", m.Folders)
	}
	if len(m.From) != 1 || m.From[0].Email != "daenerys.t@example.com" {
		t.Fatalf("from = %#v", m.From)
	}
	if m.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("grant_id = %q", m.GrantID)
	}
	if m.ID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("id = %q", m.ID)
	}
	if m.Object != "message" {
		t.Fatalf("object = %q", m.Object)
	}
	if m.Snippet == nil || *m.Snippet != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("snippet = %#v", m.Snippet)
	}
	if m.Starred == nil || *m.Starred != true {
		t.Fatalf("starred = %#v", m.Starred)
	}
	if m.Subject == nil || *m.Subject != "Hello from Nylas!" {
		t.Fatalf("subject = %#v", m.Subject)
	}
	if m.ThreadID == nil || *m.ThreadID != "1t8tv3890q4vgmwq6pmdwm8qgsaer" {
		t.Fatalf("thread_id = %#v", m.ThreadID)
	}
	if len(m.To) != 1 || m.To[0].Email != "j.snow@example.com" {
		t.Fatalf("to = %#v", m.To)
	}
	if m.Unread == nil || *m.Unread != true {
		t.Fatalf("unread = %#v", m.Unread)
	}
	if m.Metadata == nil || m.Metadata["custom_field"] != "value" {
		t.Fatalf("metadata = %#v", m.Metadata)
	}
}

func TestListMessages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{"id": "m1", "grant_id": "abc-123"},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().List(context.Background(), "abc-123", nil); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListMessagesWithQueryParams(t *testing.T) {
	sub := "Hello from Nylas!"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("subject") != sub {
			t.Fatalf("subject query mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Messages().List(context.Background(), "abc-123", &models.ListMessagesQueryParams{Subject: &sub})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListMessagesWithSelectParam(t *testing.T) {
	sel := "id,subject,from,to"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("select") != sel {
			t.Fatalf("select mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": []any{}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().List(context.Background(), "abc-123", &models.ListMessagesQueryParams{Select: &sel}); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/message-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"id": "message-123", "grant_id": "abc-123"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Find(context.Background(), "abc-123", "message-123", nil); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestFindMessageEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	enc := "/v3/grants/abc-123/messages/" + url.PathEscape(raw)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, enc)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": raw}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Find(context.Background(), "abc-123", raw, nil); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestFindMessageWithQueryParams(t *testing.T) {
	fields := models.MessageFieldsStandard
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") != string(fields) {
			t.Fatalf("fields mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/message-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "message-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Find(context.Background(), "abc-123", "message-123", &models.FindMessageQueryParams{Fields: &fields}); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestUpdateMessage(t *testing.T) {
	req := models.UpdateMessageRequest{
		Starred:  boolptr(true),
		Unread:   boolptr(false),
		Folders:  []string{"folder-123"},
		Metadata: map[string]any{"foo": "bar"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/messages/message-123")
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["starred"] != true || body["unread"] != false {
			t.Fatalf("unexpected body: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "message-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Update(context.Background(), "abc-123", "message-123", req); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestUpdateMessageEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	enc := "/v3/grants/abc-123/messages/" + url.PathEscape(raw)
	req := models.UpdateMessageRequest{Starred: boolptr(true)}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, enc)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": raw}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Update(context.Background(), "abc-123", raw, req); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDeleteMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/messages/message-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Destroy(context.Background(), "abc-123", "message-123"); err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
}

func TestDeleteMessageEncodedID(t *testing.T) {
	raw := "<!&!AAAAAAAAAAAuAAAAAAAAABQ/wHZyqaNCptfKg5rnNAoBAMO2jhD3dRHOtM0AqgC7tuYAAAAAAA4AABAAAACTn3BxdTQ/T4N/0BgqPmf+AQAAAAA=@example.com>"
	enc := "/v3/grants/abc-123/messages/" + url.PathEscape(raw)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, enc)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Destroy(context.Background(), "abc-123", raw); err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
}

func TestSendMessageJSON(t *testing.T) {
	subject := "Hello from Nylas!"
	body := "This is the body of my draft message."
	req := models.SendMessageRequest{
		CreateDraftRequest: models.CreateDraftRequest{
			Subject:  &subject,
			Body:     &body,
			To:       []models.EmailName{{Name: strptr("Jon Snow"), Email: "jsnow@gmail.com"}},
			Cc:       []models.EmailName{{Name: strptr("Arya Stark"), Email: "astark@gmail.com"}},
			Metadata: map[string]any{"custom_field": "value", "another_field": 123},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/messages/send")
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["subject"] != subject {
			t.Fatalf("subject mismatch: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "m1"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Send(context.Background(), "abc-123", req); err != nil {
		t.Fatalf("Send error: %v", err)
	}
}

func TestSendMessageMultipart(t *testing.T) {
	req := map[string]any{
		"subject": "Hello from Nylas!",
		"to":      []map[string]any{{"name": "Jon Snow", "email": "jsnow@gmail.com"}},
		"cc":      []map[string]any{{"name": "Arya Stark", "email": "astark@gmail.com"}},
		"body":    "This is the body of my draft message.",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure multipart form-data
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data; boundary=") {
			t.Fatalf("expected multipart content type, got: %s", ct)
		}
		// Read body and ensure parts exist
		all, _ := io.ReadAll(r.Body)
		s := string(all)
		if !strings.Contains(s, `name="message"`) || !strings.Contains(s, `name="file0"`) {
			t.Fatalf("multipart parts missing, got: %s", s)
		}
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/messages/send")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "m2"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	files := []FormFile{
		{Filename: "file1.txt", ContentType: "text/plain", Reader: bytes.NewBufferString("this is a file")},
	}
	if _, err := c.Messages().SendMultipart(context.Background(), "abc-123", req, files); err != nil {
		t.Fatalf("SendMultipart error: %v", err)
	}
}

func TestListScheduledMessages(t *testing.T) {
	resp := map[string]any{
		"request_id": "dd3ec9a2-8f15-403d-b269-32b1f1beb9f5",
		"data": []map[string]any{
			{
				"schedule_id": "8cd56334-6d95-432c-86d1-c5dab0ce98be",
				"status":      map[string]any{"code": "pending", "description": "schedule send awaiting send at time"},
			},
			{
				"schedule_id": "rb856334-6d95-432c-86d1-c5dab0ce98be",
				"status":      map[string]any{"code": "success", "description": "schedule send succeeded"},
				"close_time":  1690579819,
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/schedules")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Messages().ListScheduledMessages(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("ListScheduledMessages error: %v", err)
	}
	if out.RequestID != resp["request_id"].(string) {
		t.Fatalf("request_id = %q", out.RequestID)
	}
	if len(out.Data) != 2 {
		t.Fatalf("len(data) = %d", len(out.Data))
	}
	if out.Data[0].ScheduleID != "8cd56334-6d95-432c-86d1-c5dab0ce98be" ||
		out.Data[0].Status.Code != "pending" {
		t.Fatalf("first entry = %#v", out.Data[0])
	}
	if out.Data[1].ScheduleID != "rb856334-6d95-432c-86d1-c5dab0ce98be" ||
		out.Data[1].Status.Code != "success" ||
		out.Data[1].CloseTime == nil || *out.Data[1].CloseTime != 1690579819 {
		t.Fatalf("second entry = %#v", out.Data[1])
	}
}

func TestFindScheduledMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/schedules/schedule-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"schedule_id": "schedule-123"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().FindScheduledMessage(context.Background(), "abc-123", "schedule-123"); err != nil {
		t.Fatalf("FindScheduledMessage error: %v", err)
	}
}

func TestStopScheduledMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/messages/schedules/schedule-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"message": "stopped"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().StopScheduledMessage(context.Background(), "abc-123", "schedule-123"); err != nil {
		t.Fatalf("StopScheduledMessage error: %v", err)
	}
}

func TestCleanMessages(t *testing.T) {
	req := models.CleanMessagesRequest{
		MessageID:               []string{"message-1", "message-2"},
		IgnoreImages:            boolptr(true),
		IgnoreLinks:             boolptr(true),
		IgnoreTables:            boolptr(true),
		ImagesAsMarkdown:        boolptr(true),
		RemoveConclusionPhrases: boolptr(true),
	}

	// Build a realistic response
	body1 := "Hello, I just sent a message using Nylas!"
	name := "Daenerys Targaryen"
	resp := map[string]any{
		"request_id": "rid-1",
		"data": []map[string]any{
			{
				"body":         body1,
				"from":         []map[string]any{{"name": name, "email": "daenerys.t@example.com"}},
				"object":       "message",
				"id":           "message-1",
				"grant_id":     "41009df5-bf11-4c97-aa18-b285b5f2e386",
				"conversation": "cleaned example",
			},
			{
				"id":           "message-2",
				"grant_id":     "41009df5-bf11-4c97-aa18-b285b5f2e386",
				"conversation": "another example",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/messages/clean")
		var got models.CleanMessagesRequest
		_ = json.NewDecoder(r.Body).Decode(&got)
		if len(got.MessageID) != 2 || got.MessageID[0] != "message-1" {
			t.Fatalf("request mismatch: %#v", got)
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Messages().CleanMessages(context.Background(), "abc-123", req)
	if err != nil {
		t.Fatalf("CleanMessages error: %v", err)
	}
	if len(out.Data) != 2 {
		t.Fatalf("len(data) = %d", len(out.Data))
	}
	if out.Data[0].Body == nil || *out.Data[0].Body != body1 {
		t.Fatalf("body[0] = %#v", out.Data[0].Body)
	}
	if len(out.Data[0].From) != 1 || out.Data[0].From[0].Email != "daenerys.t@example.com" || out.Data[0].From[0].Name == nil || *out.Data[0].From[0].Name != name {
		t.Fatalf("from[0] = %#v", out.Data[0].From)
	}
	if out.Data[0].Object != "message" || out.Data[0].ID != "message-1" || out.Data[0].GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("message[0] = %#v", out.Data[0])
	}
	if out.Data[0].Conversation != "cleaned example" || out.Data[1].Conversation != "another example" {
		t.Fatalf("conversation mismatch: %#v | %#v", out.Data[0].Conversation, out.Data[1].Conversation)
	}
}

func TestListMessagesWithFieldsIncludeTrackingOptions(t *testing.T) {
	fields := models.MessageFieldsIncludeTracking
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") != string(fields) {
			t.Fatalf("fields mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": []any{}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().List(context.Background(), "abc-123", &models.ListMessagesQueryParams{Fields: &fields}); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListMessagesWithFieldsRawMIME(t *testing.T) {
	fields := models.MessageFieldsRawMIME
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") != string(fields) {
			t.Fatalf("fields mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": []any{}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().List(context.Background(), "abc-123", &models.ListMessagesQueryParams{Fields: &fields}); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindMessageWithIncludeTrackingOptionsField(t *testing.T) {
	fields := models.MessageFieldsIncludeTracking
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") != string(fields) {
			t.Fatalf("fields mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/message-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "message-123"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Find(context.Background(), "abc-123", "message-123", &models.FindMessageQueryParams{Fields: &fields}); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestFindMessageWithRawMIMEField(t *testing.T) {
	fields := models.MessageFieldsRawMIME
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") != string(fields) {
			t.Fatalf("fields mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/messages/message-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "x", "data": map[string]any{"id": "message-123"}})
	}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Messages().Find(context.Background(), "abc-123", "message-123", &models.FindMessageQueryParams{Fields: &fields}); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestMessageDeserializationWithTrackingOptions(t *testing.T) {
	js := `{
		"grant_id": "g1",
		"id": "m1",
		"object": "message",
		"subject": "Hello from Nylas!",
		"tracking_options": {
			"opens": true,
			"thread_replies": false,
			"links": true,
			"label": "Marketing Campaign"
		}
	}`
	var m models.Message
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.Tracking == nil || m.Tracking.Opens == nil || *m.Tracking.Opens != true ||
		m.Tracking.ThreadReplies == nil || *m.Tracking.ThreadReplies != false ||
		m.Tracking.Links == nil || *m.Tracking.Links != true ||
		m.Tracking.Label == nil || *m.Tracking.Label != "Marketing Campaign" {
		t.Fatalf("tracking_options = %#v", m.Tracking)
	}
}

func TestMessageDeserializationWithRawMIME(t *testing.T) {
	raw := "TUlNRS1WZXJzaW9uOiAxLjAKQ29udGVudC1UeXBlOiB0ZXh0L3BsYWluOyBjaGFyc2V0PXV0Zi04CgpIZWxsbyBXb3JsZCE="
	js := `{"grant_id":"g1","id":"m1","object":"message","raw_mime":"` + raw + `"}`
	var m models.Message
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.RawMIME == nil || *m.RawMIME != raw {
		t.Fatalf("raw_mime = %#v", m.RawMIME)
	}
}

func TestMessageDeserializationBackwardsCompat(t *testing.T) {
	js := `{
		"grant_id": "g1",
		"id": "m1",
		"object": "message",
		"body": "Hello, I just sent a message using Nylas!",
		"subject": "Hello from Nylas!"
	}`
	var m models.Message
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.Tracking != nil || m.RawMIME != nil {
		t.Fatalf("unexpected tracking/raw_mime: %#v %#v", m.Tracking, m.RawMIME)
	}
	if m.Body == nil || *m.Body != "Hello, I just sent a message using Nylas!" {
		t.Fatalf("body = %#v", m.Body)
	}
	if m.Subject == nil || *m.Subject != "Hello from Nylas!" {
		t.Fatalf("subject = %#v", m.Subject)
	}
}
