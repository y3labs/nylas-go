package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestGrantDeserialization(t *testing.T) {
	js := `{
		"id": "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
		"provider": "google",
		"grant_status": "valid",
		"email": "email@example.com",
		"scope": ["Mail.Read", "User.Read", "offline_access"],
		"user_agent": "string",
		"ip": "string",
		"state": "my-state",
		"created_at": 1617817109,
		"updated_at": 1617817109
	}`

	var g models.Grant
	if err := json.Unmarshal([]byte(js), &g); err != nil {
		t.Fatalf("unmarshal grant: %v", err)
	}

	if g.ID != "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47" {
		t.Fatalf("id = %q", g.ID)
	}
	if g.Provider != "google" {
		t.Fatalf("provider = %q", g.Provider)
	}
	if g.GrantStatus == nil || *g.GrantStatus != "valid" {
		t.Fatalf("grant_status = %#v", g.GrantStatus)
	}
	if g.Email == nil || *g.Email != "email@example.com" {
		t.Fatalf("email = %#v", g.Email)
	}
	wantScope := []string{"Mail.Read", "User.Read", "offline_access"}
	if len(g.Scope) != len(wantScope) ||
		g.Scope[0] != wantScope[0] ||
		g.Scope[1] != wantScope[1] ||
		g.Scope[2] != wantScope[2] {
		t.Fatalf("scope = %#v", g.Scope)
	}
	if g.UserAgent == nil || *g.UserAgent != "string" {
		t.Fatalf("user_agent = %#v", g.UserAgent)
	}
	if g.IP == nil || *g.IP != "string" {
		t.Fatalf("ip = %#v", g.IP)
	}
	if g.State == nil || *g.State != "my-state" {
		t.Fatalf("state = %#v", g.State)
	}
	if g.CreatedAt == nil || *g.CreatedAt != 1617817109 {
		t.Fatalf("created_at = %#v", g.CreatedAt)
	}
	if g.UpdatedAt == nil || *g.UpdatedAt != 1617817109 {
		t.Fatalf("updated_at = %#v", g.UpdatedAt)
	}
}

func TestListGrants(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{"id": "g1", "provider": "google", "scope": []string{"a"}},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Grants().List(context.Background(), nil); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestGetGrant(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":       "grant-123",
				"provider": "google",
				"scope":    []string{"Mail.Read"},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Grants().Get(context.Background(), "grant-123"); err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestUpdateGrant(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/grants/grant-123")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["settings"]; !ok {
			t.Fatalf("missing settings in body: %#v", body)
		}
		if scope, ok := body["scope"].([]any); !ok || len(scope) != 2 {
			t.Fatalf("unexpected scope in body: %#v", body["scope"])
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":       "grant-123",
				"provider": "google",
				"scope":    []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := UpdateGrantRequest{
		Settings: map[string]interface{}{
			"client_id":     "string",
			"client_secret": "string",
		},
		Scope: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	if _, err := c.Grants().Update(context.Background(), "grant-123", req); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDeleteGrant(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/grant-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if err := c.Grants().Delete(context.Background(), "grant-123"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}
