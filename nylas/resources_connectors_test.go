package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestConnectorsCredentialsProperty(t *testing.T) {
	// No server needed; this doesn't make a request.
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	conn := c.Connectors()
	if conn == nil {
		t.Fatalf("Connectors() returned nil")
	}
	creds := conn.Credentials()
	if creds == nil {
		t.Fatalf("Connectors.Credentials() returned nil")
	}
	if _, ok := any(creds).(*CredentialsResource); !ok {
		t.Fatalf("Credentials() returned wrong type: %T", creds)
	}
}

func TestConnectorDeserialization(t *testing.T) {
	js := []byte(`{
		"provider": "google",
		"settings": {"topic_name": "abc123"},
		"scope": [
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"
		]
	}`)
	var c models.Connector
	if err := json.Unmarshal(js, &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Provider != models.ProviderGoogle {
		t.Fatalf("Provider = %q, want %q", c.Provider, models.ProviderGoogle)
	}
	if c.Settings["topic_name"] != "abc123" {
		t.Fatalf("Settings[topic_name] = %v, want %q", c.Settings["topic_name"], "abc123")
	}
	wantScope := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	if !reflect.DeepEqual(c.Scope, wantScope) {
		t.Fatalf("Scope = %#v, want %#v", c.Scope, wantScope)
	}
}

func TestListConnectors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connectors")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []map[string]any{
				{"provider": "google", "settings": map[string]any{"topic_name": "x"}, "scope": []string{}},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if out == nil || len(out.Data) != 1 {
		t.Fatalf("unexpected list response: %#v", out)
	}
}

func TestFindConnector(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connectors/google")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"provider": "google",
				"settings": map[string]any{"topic_name": "abc123"},
				"scope":    []string{},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().Get(context.Background(), models.ProviderGoogle)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if out == nil || out.Data.Provider != models.ProviderGoogle {
		t.Fatalf("unexpected get response: %#v", out)
	}
}

func TestCreateConnector(t *testing.T) {
	var captured map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connectors")
		captured = decodeJSONBody(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       captured, // echo back
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	body := models.GoogleCreateConnectorRequest{
		Provider: models.ProviderGoogle,
		Settings: models.GoogleCreateConnectorSettings{
			ClientID:     "string",
			ClientSecret: "string",
			TopicName:    strptr("string"),
		},
		Scope: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	out, err := c.Connectors().Create(context.Background(), body)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if out == nil {
		t.Fatalf("Create returned nil response")
	}

	// Validate posted JSON matches expectation
	if got := captured["provider"]; got != "google" {
		t.Fatalf("provider = %v, want %q", got, "google")
	}
	s, ok := captured["settings"].(map[string]any)
	if !ok {
		t.Fatalf("settings missing or wrong type: %#v", captured["settings"])
	}
	if s["client_id"] != "string" || s["client_secret"] != "string" || s["topic_name"] != "string" {
		t.Fatalf("unexpected settings: %#v", s)
	}
	scope, ok := captured["scope"].([]any)
	if !ok || len(scope) != 2 {
		t.Fatalf("unexpected scope: %#v", captured["scope"])
	}
}

func TestUpdateConnector(t *testing.T) {
	var captured map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/connectors/google")
		captured = decodeJSONBody(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"provider": "google",
				"settings": captured["settings"],
				"scope":    captured["scope"],
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	req := models.UpdateConnectorRequest{
		Settings: map[string]any{
			"client_id":     "string",
			"client_secret": "string",
			"topic_name":    "string",
		},
		Scope: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	out, err := c.Connectors().Update(context.Background(), models.ProviderGoogle, req)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if out == nil {
		t.Fatalf("Update returned nil response")
	}

	// Validate JSON body posted
	if _, ok := captured["settings"].(map[string]any); !ok {
		t.Fatalf("settings missing/wrong type: %#v", captured["settings"])
	}
	if _, ok := captured["scope"].([]any); !ok {
		t.Fatalf("scope missing/wrong type: %#v", captured["scope"])
	}
}

func TestDestroyConnector(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/connectors/google")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().Delete(context.Background(), models.ProviderGoogle)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if out == nil || out.RequestID == "" {
		t.Fatalf("unexpected delete response: %#v", out)
	}
}

func decodeJSONBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	defer r.Body.Close()
	var m map[string]any
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return m
}
