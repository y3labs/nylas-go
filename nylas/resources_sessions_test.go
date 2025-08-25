package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestSession_Deserialization(t *testing.T) {
	js := `{"session_id":"session-id"}`
	var s models.Session
	if err := json.Unmarshal([]byte(js), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.SessionID != "session-id" {
		t.Fatalf("session_id mismatch: %q", s.SessionID)
	}
}

func TestSessions_Create(t *testing.T) {
	// Mock server validates method, path and request body; returns a created session
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/scheduling/sessions")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		var got models.CreateSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.ConfigurationID == nil || *got.ConfigurationID != "configuration-123" {
			t.Fatalf("ConfigurationID mismatch: %#v", got.ConfigurationID)
		}
		if got.TimeToLive == nil || *got.TimeToLive != 30 {
			t.Fatalf("TimeToLive mismatch: %#v", got.TimeToLive)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"session_id": "session-id",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	cfg := "configuration-123"
	ttl := 30
	req := models.CreateSessionRequest{
		ConfigurationID: &cfg,
		TimeToLive:      &ttl,
	}

	res, err := c.Scheduler().Sessions().Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if res.Data.SessionID != "session-id" {
		t.Fatalf("unexpected session: %#v", res.Data)
	}
	if res.RequestID == "" {
		t.Fatalf("missing request_id on response")
	}
}

func TestSessions_Destroy(t *testing.T) {
	// Mock server validates method/path, returns a basic delete response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/scheduling/sessions/session-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-del-1",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.Scheduler().Sessions().Destroy(context.Background(), "session-123")
	if err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
	if res == nil || res.RequestID == "" {
		t.Fatalf("expected request_id in delete response, got: %#v", res)
	}
}
