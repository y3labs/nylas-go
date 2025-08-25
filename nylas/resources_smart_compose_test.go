package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestMessages_SmartComposeAccessor(t *testing.T) {
	// No server needed; just ensure accessor is wired
	c := newTestClient("http://example.invalid", "test-key")
	sc := c.Messages().SmartCompose()
	if sc == nil {
		t.Fatalf("SmartCompose accessor returned nil")
	}
}

func TestSmartCompose_Deserialization(t *testing.T) {
	js := `{"suggestion":"Hello world"}`
	var resp models.ComposeMessageResponse
	if err := json.Unmarshal([]byte(js), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Suggestion != "Hello world" {
		t.Fatalf("suggestion mismatch: %q", resp.Suggestion)
	}
}

func TestSmartCompose_ComposeMessage(t *testing.T) {
	// Mock API endpoint for composing a new message
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/grant-123/messages/smart-compose")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		var got models.ComposeMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.Prompt != "Hello world" {
			t.Fatalf("prompt mismatch: %#v", got.Prompt)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"suggestion": "Hello world",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.ComposeMessageRequest{Prompt: "Hello world"}
	resp, err := c.Messages().SmartCompose().ComposeMessage(context.Background(), "grant-123", req)
	if err != nil {
		t.Fatalf("ComposeMessage error: %v", err)
	}
	if resp.RequestID == "" {
		t.Fatalf("missing request_id")
	}
	if resp.Data.Suggestion != "Hello world" {
		t.Fatalf("unexpected suggestion: %#v", resp.Data)
	}
}

func TestSmartCompose_ComposeMessageReply(t *testing.T) {
	// Mock API endpoint for composing a reply
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/grant-123/messages/message-123/smart-compose")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		var got models.ComposeMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.Prompt != "Hello world" {
			t.Fatalf("prompt mismatch: %#v", got.Prompt)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-456",
			"data": map[string]any{
				"suggestion": "Hello world",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.ComposeMessageRequest{Prompt: "Hello world"}
	resp, err := c.Messages().SmartCompose().ComposeMessageReply(context.Background(), "grant-123", "message-123", req)
	if err != nil {
		t.Fatalf("ComposeMessageReply error: %v", err)
	}
	if resp.RequestID == "" {
		t.Fatalf("missing request_id")
	}
	if resp.Data.Suggestion != "Hello world" {
		t.Fatalf("unexpected suggestion: %#v", resp.Data)
	}
}
