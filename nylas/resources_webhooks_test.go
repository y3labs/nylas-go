package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestWebhookDeserialization(t *testing.T) {
	js := `{
		"id": "UMWjAjMeWQ4D8gYF2moonK4486",
		"description": "Production webhook destination",
		"trigger_types": ["calendar.created"],
		"webhook_url": "https://example.com/webhooks",
		"status": "active",
		"notification_email_addresses": ["jane@example.com", "joe@example.com"],
		"status_updated_at": 1234567890,
		"created_at": 1234567890,
		"updated_at": 1234567890
	}`

	var w models.Webhook
	if err := json.Unmarshal([]byte(js), &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if w.ID != "UMWjAjMeWQ4D8gYF2moonK4486" {
		t.Fatalf("id mismatch: %q", w.ID)
	}
	if w.Description == nil || *w.Description != "Production webhook destination" {
		t.Fatalf("description mismatch: %#v", w.Description)
	}
	if len(w.TriggerTypes) != 1 || w.TriggerTypes[0] != models.TriggerCalendarCreated {
		t.Fatalf("trigger_types mismatch: %#v", w.TriggerTypes)
	}
	if w.WebhookURL != "https://example.com/webhooks" {
		t.Fatalf("webhook_url mismatch: %q", w.WebhookURL)
	}
	if w.Status != models.WebhookStatusActive {
		t.Fatalf("status mismatch: %q", w.Status)
	}
	if got := w.NotificationEmailAddresses; len(got) != 2 || got[0] != "jane@example.com" || got[1] != "joe@example.com" {
		t.Fatalf("notification_email_addresses mismatch: %#v", got)
	}
	if w.StatusUpdatedAt != 1234567890 || w.CreatedAt != 1234567890 || w.UpdatedAt != 1234567890 {
		t.Fatalf("timestamps mismatch: status_updated_at=%d created_at=%d updated_at=%d", w.StatusUpdatedAt, w.CreatedAt, w.UpdatedAt)
	}
}

func TestWebhookDeserialization_AllTriggers(t *testing.T) {
	// Use explicit strings to mirror the Python test list (including legacy entries).
	trigs := []models.WebhookTrigger{
		"booking.created",
		"booking.pending",
		"booking.rescheduled",
		"booking.cancelled",
		"booking.reminder",
		"calendar.created",
		"calendar.updated",
		"calendar.deleted",
		"contact.updated",
		"contact.deleted",
		"event.created",
		"event.updated",
		"event.deleted",
		"grant.created",
		"grant.updated",
		"grant.deleted",
		"grant.expired",
		"message.send_success",
		"message.send_failed",
		"message.bounce_detected",
		"message.created",
		"message.updated",
		"message.opened",
		"message.link_clicked",
		"message.opened.legacy",
		"message.link_clicked.legacy",
		"message.intelligence.order",
		"message.intelligence.tracking",
		"message.intelligence.return",
		"thread.replied",
		"thread.replied.legacy",
		"folder.created",
		"folder.updated",
		"folder.deleted",
	}

	payload := map[string]any{
		"id":                           "UMWjAjMeWQ4D8gYF2moonK4486",
		"description":                  "Production webhook destination",
		"trigger_types":                trigs,
		"webhook_url":                  "https://example.com/webhooks",
		"status":                       "active",
		"notification_email_addresses": []string{"jane@example.com", "joe@example.com"},
		"status_updated_at":            1234567890,
		"created_at":                   1234567890,
		"updated_at":                   1234567890,
	}
	b, _ := json.Marshal(payload)

	var w models.Webhook
	if err := json.Unmarshal(b, &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if w.ID != "UMWjAjMeWQ4D8gYF2moonK4486" ||
		w.Description == nil || *w.Description != "Production webhook destination" ||
		w.WebhookURL != "https://example.com/webhooks" ||
		w.Status != models.WebhookStatusActive {
		t.Fatalf("basic fields mismatch: %#v", w)
	}
	if len(w.TriggerTypes) != len(trigs) {
		t.Fatalf("trigger_types length mismatch: %d vs %d", len(w.TriggerTypes), len(trigs))
	}
	for i := range trigs {
		if w.TriggerTypes[i] != trigs[i] {
			t.Fatalf("trigger_types[%d] mismatch: %q vs %q", i, w.TriggerTypes[i], trigs[i])
		}
	}
}

func TestListWebhooks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/webhooks")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": []map[string]any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().List(context.Background()); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/webhooks/webhook-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"id":                           "webhook-123",
				"trigger_types":                []string{"calendar.created"},
				"webhook_url":                  "https://example.com/webhooks",
				"status":                       "active",
				"notification_email_addresses": []string{"jane@example.com"},
				"status_updated_at":            1,
				"created_at":                   1,
				"updated_at":                   1,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().Find(context.Background(), "webhook-123"); err != nil {
		t.Fatalf("Find error: %v", err)
	}
}

func TestCreateWebhook(t *testing.T) {
	req := models.CreateWebhookRequest{
		TriggerTypes: []models.WebhookTrigger{models.TriggerEventCreated},
		WebhookURL:   "https://example.com/webhooks",
		Description:  ptr("Production webhook destination"),
		NotificationEmailAddresses: []string{
			"jane@test.com",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/webhooks")

		var got models.CreateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(got.TriggerTypes) != 1 || got.TriggerTypes[0] != models.TriggerEventCreated {
			t.Fatalf("trigger_types mismatch: %#v", got.TriggerTypes)
		}
		if got.WebhookURL != "https://example.com/webhooks" {
			t.Fatalf("webhook_url mismatch: %q", got.WebhookURL)
		}
		if got.Description == nil || *got.Description != "Production webhook destination" {
			t.Fatalf("description mismatch: %#v", got.Description)
		}
		if len(got.NotificationEmailAddresses) != 1 || got.NotificationEmailAddresses[0] != "jane@test.com" {
			t.Fatalf("notification_email_addresses mismatch: %#v", got.NotificationEmailAddresses)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"id":            "new-webhook",
				"webhook_url":   got.WebhookURL,
				"trigger_types": got.TriggerTypes,
				"status":        "active",
				"created_at":    1, "updated_at": 1, "status_updated_at": 1,
				"webhook_secret": "secret",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().Create(context.Background(), req); err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func TestUpdateWebhook(t *testing.T) {
	req := models.UpdateWebhookRequest{
		TriggerTypes: []models.WebhookTrigger{models.TriggerEventCreated},
		WebhookURL:   ptr("https://example.com/webhooks"),
		Description:  ptr("Production webhook destination"),
		NotificationEmailAddresses: []string{
			"jane@test.com",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/webhooks/webhook-123")

		var got models.UpdateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(got.TriggerTypes) != 1 || got.TriggerTypes[0] != models.TriggerEventCreated {
			t.Fatalf("trigger_types mismatch: %#v", got.TriggerTypes)
		}
		if got.WebhookURL == nil || *got.WebhookURL != "https://example.com/webhooks" {
			t.Fatalf("webhook_url mismatch: %#v", got.WebhookURL)
		}
		if got.Description == nil || *got.Description != "Production webhook destination" {
			t.Fatalf("description mismatch: %#v", got.Description)
		}
		if len(got.NotificationEmailAddresses) != 1 || got.NotificationEmailAddresses[0] != "jane@test.com" {
			t.Fatalf("notification_email_addresses mismatch: %#v", got.NotificationEmailAddresses)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"id":            "webhook-123",
				"webhook_url":   *got.WebhookURL,
				"trigger_types": got.TriggerTypes,
				"status":        "active",
				"created_at":    1, "updated_at": 1, "status_updated_at": 1,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().Update(context.Background(), "webhook-123", req); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDestroyWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/webhooks/webhook-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"status": "ok",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	resp, err := c.Webhooks().Destroy(context.Background(), "webhook-123")
	if err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
	if resp == nil || resp.RequestID == "" || resp.Data == nil || resp.Data.Status != "ok" {
		t.Fatalf("unexpected destroy response: %#v", resp)
	}
}

func TestRotateWebhookSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/webhooks/webhook-123/rotate-secret")

		// Ensure an empty JSON object was sent
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty object body, got: %#v", got)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"id":             "webhook-123",
				"webhook_url":    "https://example.com/webhooks",
				"trigger_types":  []string{"calendar.created"},
				"status":         "active",
				"webhook_secret": "new-secret",
				"created_at":     1, "updated_at": 1, "status_updated_at": 1,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().RotateSecret(context.Background(), "webhook-123"); err != nil {
		t.Fatalf("RotateSecret error: %v", err)
	}
}

func TestWebhookIPAddresses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/webhooks/ip-addresses")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid",
			"data": map[string]any{
				"ip_addresses": []string{"1.2.3.4", "5.6.7.8"},
				"updated_at":   123,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "key")
	if _, err := c.Webhooks().IPAddresses(context.Background()); err != nil {
		t.Fatalf("IPAddresses error: %v", err)
	}
}

func TestExtractChallengeParameter(t *testing.T) {
	got, err := ExtractChallengeParameter("https://example.com/webhooks?challenge=abc123")
	if err != nil {
		t.Fatalf("ExtractChallengeParameter error: %v", err)
	}
	if got != "abc123" {
		t.Fatalf("challenge mismatch: %q", got)
	}
}

func TestExtractChallengeParameter_NoChallenge(t *testing.T) {
	_, err := ExtractChallengeParameter("https://example.com/webhooks")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "invalid URL or no challenge parameter found" {
		t.Fatalf("unexpected error: %v", err)
	}
}
