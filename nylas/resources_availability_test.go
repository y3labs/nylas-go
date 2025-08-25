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
