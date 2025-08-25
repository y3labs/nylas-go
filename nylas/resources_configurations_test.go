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

// -----------------------
// Deserialization
// -----------------------

func TestConfigurationDeserialization(t *testing.T) {
	const cfgJSON = `
{
  "id": "abc-123-configuration-id",
  "slug": null,
  "participants": [
    {
      "email": "test@nylas.com",
      "is_organizer": true,
      "name": "Test",
      "availability": {
        "calendar_ids": ["primary"],
        "open_hours": [
          {
            "days": [0,1,2,3,4,5,6],
            "exdates": null,
            "timezone": "",
            "start": "09:00",
            "end": "17:00"
          }
        ]
      },
      "booking": { "calendar_id": "primary" },
      "timezone": ""
    }
  ],
  "requires_session_auth": false,
  "availability": {
    "duration_minutes": 30,
    "interval_minutes": 15,
    "round_to": 15,
    "availability_rules": {
      "availability_method": "collective",
      "buffer": { "before": 60, "after": 0 },
      "default_open_hours": [
        {
          "days": [0,1,2,5,6],
          "exdates": null,
          "timezone": "",
          "start": "09:00",
          "end": "18:00"
        }
      ],
      "round_robin_group_id": ""
    }
  },
  "event_booking": {
    "title": "Updated Title",
    "timezone": "utc",
    "description": "",
    "location": "none",
    "booking_type": "booking",
    "conferencing": {
      "provider": "Microsoft Teams",
      "autocreate": {
        "conf_grant_id": "",
        "conf_settings": null
      }
    },
    "hide_participants": null,
    "disable_emails": null
  },
  "scheduler": {
    "available_days_in_future": 7,
    "min_cancellation_notice": 60,
    "min_booking_notice": 120,
    "confirmation_redirect_url": "",
    "hide_rescheduling_options": false,
    "hide_cancellation_options": false,
    "hide_additional_guests": true,
    "cancellation_policy": "",
    "email_template": { "booking_confirmed": {} }
  },
  "appearance": {
    "submit_button_label": "submit",
    "thank_you_message": "thank you for your business. your booking was successful."
  }
}`

	var cfg models.Configuration
	if err := json.Unmarshal([]byte(cfgJSON), &cfg); err != nil {
		t.Fatalf("unmarshal Configuration: %v", err)
	}

	if cfg.ID != "abc-123-configuration-id" {
		t.Fatalf("ID = %q, want %q", cfg.ID, "abc-123-configuration-id")
	}
	if cfg.Slug != nil {
		t.Fatalf("Slug = %v, want nil", *cfg.Slug)
	}
	if len(cfg.Participants) != 1 {
		t.Fatalf("Participants len = %d, want 1", len(cfg.Participants))
	}
	p := cfg.Participants[0]
	if p.Email != "test@nylas.com" {
		t.Fatalf("participant.email = %q", p.Email)
	}
	if p.IsOrganizer == nil || !*p.IsOrganizer {
		t.Fatalf("participant.is_organizer = %v, want true", p.IsOrganizer)
	}
	if p.Name == nil || *p.Name != "Test" {
		t.Fatalf("participant.name = %v, want %q", p.Name, "Test")
	}
	if !reflect.DeepEqual(p.Availability.CalendarIDs, []string{"primary"}) {
		t.Fatalf("participant.availability.calendar_ids = %#v", p.Availability.CalendarIDs)
	}
	if p.Availability.OpenHours == nil || len(p.Availability.OpenHours) != 1 {
		t.Fatalf("participant.availability.open_hours invalid: %#v", p.Availability.OpenHours)
	}
	oh := p.Availability.OpenHours[0]
	if !reflect.DeepEqual(oh.Days, []int{0, 1, 2, 3, 4, 5, 6}) {
		t.Fatalf("open_hours.days = %#v", oh.Days)
	}
	if oh.Timezone != "" || oh.Start != "09:00" || oh.End != "17:00" {
		t.Fatalf("open_hours tz/start/end = %q/%q/%q, want \"\"/09:00/17:00", oh.Timezone, oh.Start, oh.End)
	}
	if p.Booking.CalendarID != "primary" {
		t.Fatalf("participant.booking.calendar_id = %q, want %q", p.Booking.CalendarID, "primary")
	}
	if p.Timezone == nil || *p.Timezone != "" {
		t.Fatalf("participant.timezone = %v, want empty string ptr", p.Timezone)
	}

	if cfg.RequiresSessionAuth == nil || *cfg.RequiresSessionAuth {
		t.Fatalf("requires_session_auth = %v, want false", cfg.RequiresSessionAuth)
	}

	// Availability
	if cfg.Availability.DurationMinutes != 30 ||
		cfg.Availability.IntervalMinutes == nil || *cfg.Availability.IntervalMinutes != 15 ||
		cfg.Availability.RoundTo == nil || *cfg.Availability.RoundTo != 15 {
		t.Fatalf("availability fields incorrect: %#v", cfg.Availability)
	}
	rules := cfg.Availability.AvailabilityRules
	if rules == nil {
		t.Fatalf("availability_rules is nil")
	}
	// We only assert a few key fields; the struct typed fields mirror the Python dict
	if rules.Buffer == nil || rules.Buffer.Before == nil || *rules.Buffer.Before != 60 ||
		rules.Buffer.After == nil || *rules.Buffer.After != 0 {
		t.Fatalf("buffer incorrect: %#v", rules.Buffer)
	}
	if len(rules.DefaultOpenHours) != 1 {
		t.Fatalf("default_open_hours len = %d", len(rules.DefaultOpenHours))
	}
	doh := rules.DefaultOpenHours[0]
	if !reflect.DeepEqual(doh.Days, []int{0, 1, 2, 5, 6}) ||
		doh.Timezone != "" || doh.Start != "09:00" || doh.End != "18:00" {
		t.Fatalf("default_open_hours incorrect: %#v", doh)
	}

	// Event booking
	eb := cfg.EventBooking
	if eb.Title != "Updated Title" ||
		eb.Timezone == nil || *eb.Timezone != "utc" ||
		eb.Location == nil || *eb.Location != "none" ||
		eb.BookingType == nil || *eb.BookingType != "booking" {
		t.Fatalf("event_booking fields incorrect: %#v", eb)
	}
	if eb.Conferencing == nil {
		t.Fatalf("event_booking.conferencing is nil")
	}
	// Accept either autocreate or details, but provider must be "Microsoft Teams"
	switch {
	case eb.Conferencing.Autocreate != nil:
		if eb.Conferencing.Autocreate.Provider != "Microsoft Teams" {
			t.Fatalf("conferencing.autocreate.provider = %q, want %q", eb.Conferencing.Autocreate.Provider, "Microsoft Teams")
		}
	case eb.Conferencing.Details != nil:
		if eb.Conferencing.Details.Provider != "Microsoft Teams" {
			t.Fatalf("conferencing.details.provider = %q, want %q", eb.Conferencing.Details.Provider, "Microsoft Teams")
		}
	default:
		t.Fatalf("conferencing missing details/autocreate: %#v", eb.Conferencing)
	}

	// Scheduler settings
	if cfg.Scheduler == nil {
		t.Fatalf("scheduler is nil")
	}
	ss := cfg.Scheduler
	if ss.AvailableDaysInFuture == nil || *ss.AvailableDaysInFuture != 7 ||
		ss.MinCancellationNotice == nil || *ss.MinCancellationNotice != 60 ||
		ss.MinBookingNotice == nil || *ss.MinBookingNotice != 120 {
		t.Fatalf("scheduler numeric settings incorrect: %#v", ss)
	}
	if cfg.Appearance == nil || cfg.Appearance["submit_button_label"] != "submit" {
		t.Fatalf("appearance submit_button_label = %q", cfg.Appearance["submit_button_label"])
	}
}

// -----------------------
// Resource: Configurations
// -----------------------

func TestListConfigurations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/scheduling/configurations")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Scheduler().Configurations().List(context.Background(), "grant-123", nil)
	if err != nil {
		t.Fatalf("List configurations error: %v", err)
	}
}

func TestFindConfiguration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/scheduling/configurations/config-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id": "config-123",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Scheduler().Configurations().Find(context.Background(), "grant-123", "config-123")
	if err != nil {
		t.Fatalf("Find configuration error: %v", err)
	}
}

func TestCreateConfiguration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/grant-123/scheduling/configurations")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		// Spot-check a few key fields the Python test sets
		if v, ok := body["requires_session_auth"]; !ok || v.(bool) != false {
			t.Fatalf("requires_session_auth missing/!=false: %#v", body["requires_session_auth"])
		}
		parts, ok := body["participants"].([]any)
		if !ok || len(parts) != 1 {
			t.Fatalf("participants invalid: %#v", body["participants"])
		}
		p0 := parts[0].(map[string]any)
		if p0["email"] != "test@nylas.com" || p0["is_organizer"].(bool) != true {
			t.Fatalf("participant fields mismatch: %#v", p0)
		}
		eb := body["event_booking"].(map[string]any)
		if eb["title"] != "My test event" {
			t.Fatalf("event_booking.title = %v", eb["title"])
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id": "config-created",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	req := models.CreateConfigurationRequest{
		RequiresSessionAuth: boolptr(false),
		Participants: []models.ConfigParticipant{
			{
				Name:        strptr("Test"),
				Email:       "test@nylas.com",
				IsOrganizer: boolptr(true),
				Availability: models.ParticipantAvailability{
					CalendarIDs: []string{"primary"},
				},
				Booking: models.ParticipantBooking{CalendarID: "primary"},
			},
		},
		Availability: models.Availability{
			DurationMinutes: 30,
		},
		EventBooking: models.EventBooking{
			Title: "My test event",
		},
	}

	_, err := c.Scheduler().Configurations().Create(context.Background(), "grant-123", req)
	if err != nil {
		t.Fatalf("Create configuration error: %v", err)
	}
}

func TestUpdateConfiguration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/grant-123/scheduling/configurations/config-123")

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		eb := body["event_booking"].(map[string]any)
		if eb["title"] != "My test event" {
			t.Fatalf("event_booking.title = %v", eb["title"])
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id": "config-123",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.UpdateConfigurationRequest{
		EventBooking: &models.EventBooking{Title: "My test event"},
	}
	_, err := c.Scheduler().Configurations().Update(context.Background(), "grant-123", "config-123", req)
	if err != nil {
		t.Fatalf("Update configuration error: %v", err)
	}
}

func TestDestroyConfiguration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/grant-123/scheduling/configurations/config-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc-123"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Scheduler().Configurations().Destroy(context.Background(), "grant-123", "config-123")
	if err != nil {
		t.Fatalf("Destroy configuration error: %v", err)
	}
}
