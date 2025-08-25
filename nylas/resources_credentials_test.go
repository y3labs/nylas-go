package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestCredentialDeserialization(t *testing.T) {
	var c models.Credential
	if err := json.Unmarshal([]byte(`{
		"id":"e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
		"name":"My first Google credential",
		"created_at":1617817109,
		"updated_at":1617817109
	}`), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.ID != "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47" {
		t.Fatalf("ID = %q", c.ID)
	}
	if c.Name != "My first Google credential" {
		t.Fatalf("Name = %q", c.Name)
	}
	if c.CreatedAt == nil || *c.CreatedAt != 1617817109 {
		t.Fatalf("CreatedAt = %#v, want 1617817109", c.CreatedAt)
	}
	if c.UpdatedAt == nil || *c.UpdatedAt != 1617817109 {
		t.Fatalf("UpdatedAt = %#v, want 1617817109", c.UpdatedAt)
	}
}

func TestListCredentials(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connectors/google/creds")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": []any{
				map[string]any{
					"id":         "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
					"name":       "My first Google credential",
					"created_at": 1617817109,
					"updated_at": 1617817109,
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().Credentials().List(context.Background(), models.ProviderGoogle, nil)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if out == nil || out.RequestID == "" || len(out.Data) != 1 {
		t.Fatalf("unexpected list response: %#v", out)
	}
}

func TestFindCredential(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connectors/google/creds/abc-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":         "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
				"name":       "My first Google credential",
				"created_at": 1617817109,
				"updated_at": 1617817109,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().Credentials().Get(context.Background(), models.ProviderGoogle, "abc-123")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if out == nil || out.RequestID == "" || out.Data.ID == "" {
		t.Fatalf("unexpected get response: %#v", out)
	}
}

func TestCreateCredential(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connectors/google/creds")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["credential_type"] != string(models.CredentialTypeServiceAccount) {
			t.Fatalf("credential_type = %v, want %q", body["credential_type"], models.CredentialTypeServiceAccount)
		}
		if body["name"] != "My first Google credential" {
			t.Fatalf("name = %v", body["name"])
		}
		cd, ok := body["credential_data"].(map[string]any)
		if !ok {
			t.Fatalf("credential_data type: %T", body["credential_data"])
		}
		if cd["private_key_id"] != "string" || cd["private_key"] != "string" || cd["client_email"] != "string" {
			t.Fatalf("credential_data = %#v", cd)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":   "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
				"name": "My first Google credential",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	name := "My first Google credential"
	req := models.CredentialRequest{
		Name:           &name,
		CredentialType: models.CredentialTypeServiceAccount,
		CredentialData: map[string]string{
			"private_key_id": "string",
			"private_key":    "string",
			"client_email":   "string",
		},
	}
	out, err := c.Connectors().Credentials().Create(context.Background(), models.ProviderGoogle, req)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if out == nil || out.RequestID == "" || out.Data.ID == "" {
		t.Fatalf("unexpected create response: %#v", out)
	}
}

func TestUpdateCredential(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/connectors/google/creds/abc-123")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["name"] != "My first Google credential" {
			t.Fatalf("name = %v", body["name"])
		}
		cd, ok := body["credential_data"].(map[string]any)
		if !ok {
			t.Fatalf("credential_data type: %T", body["credential_data"])
		}
		if cd["private_key_id"] != "string" || cd["private_key"] != "string" || cd["client_email"] != "string" {
			t.Fatalf("credential_data = %#v", cd)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":   "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
				"name": "My first Google credential",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	name := "My first Google credential"
	body := models.UpdateCredentialRequest{
		Name: &name,
		CredentialData: map[string]string{
			"private_key_id": "string",
			"private_key":    "string",
			"client_email":   "string",
		},
	}
	out, err := c.Connectors().Credentials().Update(context.Background(), models.ProviderGoogle, "abc-123", body)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if out == nil || out.RequestID == "" || out.Data.ID == "" {
		t.Fatalf("unexpected update response: %#v", out)
	}
}

func TestDeleteCredential(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/connectors/google/creds/abc-123")
		w.Header().Set("x-request-id", "req-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Connectors().Credentials().Delete(context.Background(), models.ProviderGoogle, "abc-123")
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if out == nil || out.RequestID == "" {
		t.Fatalf("unexpected delete response: %#v", out)
	}
}
