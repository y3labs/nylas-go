package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestAvailability_Check_Global(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/availability" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
		}
		// body exists (don’t depend on model shape)
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-1",
			"data":       map[string]any{}, // minimal
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := (&AvailabilityResource{c}).Check(context.Background(), models.GetAvailabilityRequest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res == nil || res.RequestID != "rid-1" {
		t.Fatalf("bad response: %#v", res)
	}
}

func TestAvailability_Check_ForGrant(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.EscapedPath(), "/v3/grants/abc-123/availability") {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-2",
			"data":       map[string]any{},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := (&AvailabilityResource{c}).CheckForGrant(context.Background(), "abc-123", models.GetAvailabilityRequest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res == nil || res.RequestID != "rid-2" {
		t.Fatalf("bad response: %#v", res)
	}
}

func TestAvailability_Check_ErrorPath(t *testing.T) {
	// 500 so DoJSON -> error, and Check should return (nil, err).
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/availability" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
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
	res, err := (&AvailabilityResource{c}).Check(context.Background(), models.GetAvailabilityRequest{})
	if err == nil || res != nil {
		t.Fatalf("expected error and nil response, got res=%#v err=%v", res, err)
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestAvailability_Check_HeaderRequestIDFallback(t *testing.T) {
	// 200 with no request_id in JSON, but X-Request-Id header present.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/availability" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-avail-global")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{}, // minimal payload
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := (&AvailabilityResource{c}).Check(context.Background(), models.GetAvailabilityRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-avail-global" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
	if got := res.Headers.Get("X-Request-Id"); got != "rid-avail-global" {
		t.Fatalf("headers not propagated: %q", got)
	}
}

func TestAvailability_CheckForGrant_ErrorPath(t *testing.T) {
	// 400 error for the grant-scoped endpoint.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/grants/ten-1/availability" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"type":    "invalid_request",
				"message": "nope",
			},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := (&AvailabilityResource{c}).CheckForGrant(context.Background(), "ten-1", models.GetAvailabilityRequest{})
	if err == nil || res != nil {
		t.Fatalf("expected error and nil response, got res=%#v err=%v", res, err)
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestAvailability_CheckForGrant_HeaderRequestIDFallback(t *testing.T) {
	// 200 with header-only request id.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/grants/ten-2/availability" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-avail-grant")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := (&AvailabilityResource{c}).CheckForGrant(context.Background(), "ten-2", models.GetAvailabilityRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-avail-grant" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
	if got := res.Headers.Get("X-Request-Id"); got != "rid-avail-grant" {
		t.Fatalf("headers not propagated: %q", got)
	}
}
