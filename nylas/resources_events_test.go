package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestEventDeserialization(t *testing.T) {
	var ev models.Event
	if err := json.Unmarshal([]byte(`{
		"busy": true,
		"calendar_id": "7d93zl2palhxqdy6e5qinsakt",
		"conferencing": {
			"provider": "Zoom Meeting",
			"details": {
				"meeting_code": "code-123456",
				"password": "password-123456",
				"url": "https://zoom.us/j/1234567890?pwd=1234567890"
			}
		},
		"created_at": 1661874192,
		"description": "Description of my new calendar",
		"hide_participants": false,
		"grant_id": "41009df5-bf11-4c97-aa18-b285b5f2e386",
		"html_link": "https://www.google.com/calendar/event?eid=bTMzcGJrNW4yYjk4bjk3OWE4Ef3feD2VuM29fMjAyMjA2MjdUMjIwMDAwWiBoYWxsYUBueWxhcy5jb20",
		"id": "5d3qmne77v32r8l4phyuksl2x_20240603T180000Z",
		"master_event_id": "5d3qmne77v32r8l4phyuksl2x",
		"location": "Roller Rink",
		"metadata": {"your_key": "your_value"},
		"object": "event",
		"organizer": {"email": "organizer@example.com", "name": ""},
		"participants": [{
			"comment": "Aristotle",
			"email": "aristotle@example.com",
			"name": "Aristotle",
			"phone_number": "+1 23456778",
			"status": "maybe"
		}],
		"read_only": false,
		"reminders": {"use_default": false, "overrides": [{"reminder_minutes": 10, "reminder_method": "email"}]},
		"recurrence": ["RRULE:FREQ=WEEKLY;BYDAY=MO", "EXDATE:20211011T000000Z"],
		"status": "confirmed",
		"title": "Birthday Party",
		"updated_at": 1661874192,
		"visibility": "private",
		"when": {
			"start_time": 1661874192,
			"end_time": 1661877792,
			"start_timezone": "America/New_York",
			"end_timezone": "America/New_York",
			"object": "timespan"
		}
	}`), &ev); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}

	if !ev.Busy {
		t.Fatalf("Busy = false")
	}
	if ev.CalendarID != "7d93zl2palhxqdy6e5qinsakt" {
		t.Fatalf("CalendarID = %q", ev.CalendarID)
	}
	if ev.Conferencing == nil || ev.Conferencing.Details == nil {
		t.Fatalf("Conferencing details missing: %#v", ev.Conferencing)
	}
	if string(ev.Conferencing.Details.Provider) != "Zoom Meeting" {
		t.Fatalf("Conferencing.Provider = %q", ev.Conferencing.Details.Provider)
	}
	if ev.Conferencing.Details.Details["meeting_code"] != "code-123456" {
		t.Fatalf("meeting_code = %#v", ev.Conferencing.Details.Details["meeting_code"])
	}
	if ev.Conferencing.Details.Details["password"] != "password-123456" {
		t.Fatalf("password = %#v", ev.Conferencing.Details.Details["password"])
	}
	if ev.Conferencing.Details.Details["url"] != "https://zoom.us/j/1234567890?pwd=1234567890" {
		t.Fatalf("url = %#v", ev.Conferencing.Details.Details["url"])
	}
	if ev.CreatedAt == nil || *ev.CreatedAt != 1661874192 {
		t.Fatalf("CreatedAt = %#v", ev.CreatedAt)
	}
	if ev.Description == nil || *ev.Description != "Description of my new calendar" {
		t.Fatalf("Description = %q", strval(ev.Description))
	}
	if ev.HideParticipants == nil || *ev.HideParticipants {
		t.Fatalf("HideParticipants = %#v", ev.HideParticipants)
	}
	if ev.GrantID != "41009df5-bf11-4c97-aa18-b285b5f2e386" {
		t.Fatalf("GrantID = %q", ev.GrantID)
	}
	if ev.HTMLLink == nil || *ev.HTMLLink == "" {
		t.Fatalf("HTMLLink = %q", strval(ev.HTMLLink))
	}
	if ev.ID != "5d3qmne77v32r8l4phyuksl2x_20240603T180000Z" {
		t.Fatalf("ID = %q", ev.ID)
	}
	if ev.MasterEventID == nil || *ev.MasterEventID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("MasterEventID = %q", strval(ev.MasterEventID))
	}
	if ev.Location == nil || *ev.Location != "Roller Rink" {
		t.Fatalf("Location = %q", strval(ev.Location))
	}
	if ev.Object != "event" {
		t.Fatalf("Object = %q", ev.Object)
	}
	if len(ev.Participants) != 1 || strval(ev.Participants[0].Name) != "Aristotle" || strval(ev.Participants[0].Comment) != "Aristotle" ||
		strval(ev.Participants[0].PhoneNumber) != "+1 23456778" {
		t.Fatalf("Participants[0] = %#v", ev.Participants[0])
	}
	if ev.Participants[0].Email != "aristotle@example.com" {
		t.Fatalf("Participant.Email = %q", ev.Participants[0].Email)
	}
	if ev.Participants[0].Status == nil || string(*ev.Participants[0].Status) != "maybe" {
		t.Fatalf("Participant.Status = %#v", ev.Participants[0].Status)
	}
	if ev.ReadOnly == nil || *ev.ReadOnly {
		t.Fatalf("ReadOnly = %#v", ev.ReadOnly)
	}
	if ev.Reminders == nil || ev.Reminders.UseDefault {
		t.Fatalf("Reminders.UseDefault = %#v", ev.Reminders)
	}
	if len(ev.Reminders.Overrides) != 1 || ev.Reminders.Overrides[0].ReminderMinutes == nil ||
		*ev.Reminders.Overrides[0].ReminderMinutes != 10 ||
		ev.Reminders.Overrides[0].ReminderMethod == nil ||
		*ev.Reminders.Overrides[0].ReminderMethod != "email" {
		t.Fatalf("Reminders.Overrides = %#v", ev.Reminders.Overrides)
	}
	if len(ev.Recurrence) != 2 || ev.Recurrence[0] != "RRULE:FREQ=WEEKLY;BYDAY=MO" || ev.Recurrence[1] != "EXDATE:20211011T000000Z" {
		t.Fatalf("Recurrence = %#v", ev.Recurrence)
	}
	if ev.Status == nil || string(*ev.Status) != "confirmed" {
		t.Fatalf("Status = %#v", ev.Status)
	}
	if ev.Title == nil || *ev.Title != "Birthday Party" {
		t.Fatalf("Title = %q", strval(ev.Title))
	}
	if ev.UpdatedAt == nil || *ev.UpdatedAt != 1661874192 {
		t.Fatalf("UpdatedAt = %#v", ev.UpdatedAt)
	}
	if ev.Visibility == nil || string(*ev.Visibility) != "private" {
		t.Fatalf("Visibility = %#v", ev.Visibility)
	}
	if ev.When.Timespan == nil || ev.When.Timespan.StartTime != 1661874192 || ev.When.Timespan.EndTime != 1661877792 {
		t.Fatalf("When.Timespan = %#v", ev.When.Timespan)
	}
	if strval(ev.When.Timespan.StartTimezone) != "America/New_York" || strval(ev.When.Timespan.EndTimezone) != "America/New_York" {
		t.Fatalf("When.Timespan TZ = %#v", ev.When.Timespan)
	}
}

func TestListEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/events")
		q := r.URL.Query()
		if q.Get("calendar_id") != "abc-123" || q.Get("limit") != "20" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	params := ListEventsParams{CalendarID: "abc-123", Limit: intptr(20)}
	if _, err := c.Events().List(context.Background(), "abc-123", params); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListEventsWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/events")
		if r.URL.Query().Get("limit") != "20" {
			t.Fatalf("limit not set: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	params := ListEventsParams{CalendarID: "abc-123", Limit: intptr(20)}
	if _, err := c.Events().List(context.Background(), "abc-123", params); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListEventsWithSelectParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/events")
		if r.URL.Query().Get("select") != "id,title,description,when" {
			t.Fatalf("select not set: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []any{map[string]any{
				"id":    "event-123",
				"title": "Team Meeting",
				"when":  map[string]any{"object": "timespan", "start_time": 1625097600, "end_time": 1625101200},
			}},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	params := ListEventsParams{CalendarID: "abc-123", Select: strptr("id,title,description,when")}
	if res, err := c.Events().List(context.Background(), "abc-123", params); err != nil || res == nil {
		t.Fatalf("List with select error: %v, res=%#v", err, res)
	}
}

func TestListImportEventsVariants(t *testing.T) {
	// plain with calendar_id
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/events/import")
			q := r.URL.Query()
			if q.Get("calendar_id") != "primary" {
				t.Fatalf("query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		_, _ = c.Events().ListImport(context.Background(), "grant-123", ListImportEventsParams{CalendarID: "primary"})
	}

	// with select
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/events/import")
			q := r.URL.Query()
			if q.Get("select") != "id,title,participants" || q.Get("calendar_id") != "primary" {
				t.Fatalf("query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		_, _ = c.Events().ListImport(context.Background(), "grant-123", ListImportEventsParams{
			CalendarID: "primary", Select: strptr("id,title,participants"),
		})
	}

	// with limit
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/events/import")
			if r.URL.Query().Get("limit") != "100" {
				t.Fatalf("limit not set: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		_, _ = c.Events().ListImport(context.Background(), "grant-123", ListImportEventsParams{
			CalendarID: "primary", Limit: intptr(100),
		})
	}

	// with time filters
	{
		start, end := int64(1672531200), int64(1704067199)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/events/import")
			q := r.URL.Query()
			if q.Get("start") != "1672531200" || q.Get("end") != "1704067199" {
				t.Fatalf("query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		_, _ = c.Events().ListImport(context.Background(), "grant-123", ListImportEventsParams{
			CalendarID: "primary", Start: &start, End: &end,
		})
	}

	// with all params
	{
		_, _ = int64(1672531200), int64(1704067199)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertMethodPath(t, r, http.MethodGet, "/v3/grants/grant-123/events/import")
			q := r.URL.Query()
			if q.Get("calendar_id") != "primary" || q.Get("limit") != "50" ||
				q.Get("start") != "1672531200" || q.Get("end") != "1704067199" ||
				q.Get("select") != "id,title,participants,when" ||
				q.Get("page_token") != "next-page-token-123" {
				t.Fatalf("query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": []any{}})
		}))
		defer ts.Close()

		c := newTestClient(ts.URL, "test-key")
		_, _ = c.Events().ListImport(context.Background(), "grant-123", ListImportEventsParams{
			CalendarID: "primary",
			Limit:      intptr(50),
			Start:      i64ptr(1672531200),
			End:        i64ptr(1704067199),
			Select:     strptr("id,title,participants,when"),
			PageToken:  strptr("next-page-token-123"),
		})
	}
}

func TestFindEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/events/event-123")
		if r.URL.Query().Get("calendar_id") != "abc-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": "event-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Events().Get(context.Background(), "abc-123", "event-123", FindEventParams{CalendarID: "abc-123"})
	if err != nil || out.Data.ID != "event-123" {
		t.Fatalf("Get error: %v, out=%#v", err, out)
	}
}

func TestCreateEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/events")
		if r.URL.Query().Get("calendar_id") != "abc-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["description"] != "Description of my new event" || body["location"] != "Los Angeles, CA" {
			t.Fatalf("body = %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": "ev"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := CreateEventRequest{
		When: models.When{
			Timespan: &models.Timespan{
				Object:        "timespan",
				StartTime:     1661874192,
				EndTime:       1661877792,
				StartTimezone: strptr("America/New_York"),
				EndTimezone:   strptr("America/New_York"),
			},
		},
		Description: strptr("Description of my new event"),
		Location:    strptr("Los Angeles, CA"),
		Metadata:    map[string]any{"your-key": "value"},
	}
	params := CreateEventParams{CalendarID: "abc-123"}
	if _, err := c.Events().Create(context.Background(), "abc-123", req, params); err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func TestUpdateEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// PATCH in Go SDK (Python used PUT)
		assertMethodPath(t, r, http.MethodPatch, "/v3/grants/abc-123/events/event-123")
		if r.URL.Query().Get("calendar_id") != "abc-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["description"] != "Updated description of my event" {
			t.Fatalf("body = %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": "event-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := UpdateEventRequest{
		When: &models.When{
			Timespan: &models.Timespan{
				Object:        "timespan",
				StartTime:     1661874192,
				EndTime:       1661877792,
				StartTimezone: strptr("America/New_York"),
				EndTimezone:   strptr("America/New_York"),
			},
		},
		Description: strptr("Updated description of my event"),
		Location:    strptr("Los Angeles, CA"),
		Metadata:    map[string]any{"your-key": "value"},
	}
	params := UpdateEventParams{CalendarID: "abc-123"}
	if _, err := c.Events().Update(context.Background(), "abc-123", "event-123", req, params); err != nil {
		t.Fatalf("Update error: %v", err)
	}
}

func TestDestroyEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/events/event-123")
		if r.URL.Query().Get("calendar_id") != "abc-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if err := c.Events().Delete(context.Background(), "abc-123", "event-123", DestroyEventParams{CalendarID: "abc-123"}); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestSendRSVP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/events/event-123/send-rsvp")
		if r.URL.Query().Get("calendar_id") != "abc-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["status"] != "yes" {
			t.Fatalf("status body = %#v", body)
		}
		w.Header().Set("x-request-id", "abc-req")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "ignored"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Events().SendRSVP(context.Background(), "abc-123", "event-123",
		SendRSVPRequest{Status: models.SendRSVPStatus("yes")},
		SendRSVPParams{CalendarID: "abc-123"},
	)
	if err != nil || out.RequestID != "abc-req" {
		t.Fatalf("SendRSVP err=%v, out=%#v", err, out)
	}
}

func TestEventWithNotetakerDeserialization(t *testing.T) {
	var ev models.Event
	if err := json.Unmarshal([]byte(`{
		"id": "event-123",
		"grant_id": "grant-123",
		"calendar_id": "calendar-123",
		"busy": true,
		"participants": [{"email": "test@example.com", "name": "Test User", "status": "yes"}],
		"when": {"start_time": 1497916800, "end_time": 1497920400, "object": "timespan"},
		"title": "Test Event with Notetaker",
		"notetaker": {
			"id": "notetaker-123",
			"name": "Custom Notetaker",
			"meeting_settings": {"video_recording": true, "audio_recording": true, "transcription": true}
		}
	}`), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.ID != "event-123" || ev.GrantID != "grant-123" || ev.CalendarID != "calendar-123" || !ev.Busy {
		t.Fatalf("core fields mismatch: %#v", ev)
	}
	if ev.Title == nil || *ev.Title != "Test Event with Notetaker" {
		t.Fatalf("Title = %q", strval(ev.Title))
	}
	if ev.Notetaker == nil || ev.Notetaker.MeetingSettings == nil {
		t.Fatalf("Notetaker missing: %#v", ev.Notetaker)
	}
	ms := ev.Notetaker.MeetingSettings
	if ms == nil || !ms.VideoRecording || !ms.AudioRecording || !ms.Transcription {
		t.Fatalf("meeting_settings not all true: %#v", ms)
	}
}

func TestCreateEventWithNotetaker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/events")
		if r.URL.Query().Get("calendar_id") != "calendar-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		n, ok := body["notetaker"].(map[string]any)
		if !ok {
			t.Fatalf("notetaker missing in body: %#v", body)
		}
		ms, ok := n["meeting_settings"].(map[string]any)
		if !ok || ms["video_recording"] != true || ms["audio_recording"] != true || ms["transcription"] != true {
			t.Fatalf("meeting_settings = %#v", ms)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": "ev"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := CreateEventRequest{
		Title: strptr("Test Event with Notetaker"),
		When: models.When{Timespan: &models.Timespan{
			Object:    "timespan",
			StartTime: 1497916800, EndTime: 1497920400,
		}},
		Participants: []models.Participant{{Email: "test@example.com", Name: strptr("Test User")}},
		Notetaker: &CreateEventNotetaker{
			Name: strptr("Custom Notetaker"),
			MeetingSettings: &models.NotetakerMeetingSettings{
				VideoRecording: true,
				AudioRecording: true,
				Transcription:  true,
			},
		},
	}
	if _, err := c.Events().Create(context.Background(), "abc-123", req, CreateEventParams{CalendarID: "calendar-123"}); err != nil {
		t.Fatalf("Create with notetaker error: %v", err)
	}
}

func TestUpdateEventWithNotetaker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/grants/abc-123/events/event-123")
		if r.URL.Query().Get("calendar_id") != "calendar-123" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		n := body["notetaker"].(map[string]any)
		ms := n["meeting_settings"].(map[string]any)
		if ms["video_recording"] != false || ms["audio_recording"] != true || ms["transcription"] != false {
			t.Fatalf("meeting_settings = %#v", ms)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req", "data": map[string]any{"id": "event-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := UpdateEventRequest{
		Title: strptr("Updated Test Event"),
		Notetaker: &models.EventNotetaker{
			ID:   strptr("notetaker-123"),
			Name: strptr("Updated Notetaker"),
			MeetingSettings: &models.NotetakerMeetingSettings{
				VideoRecording: false,
				AudioRecording: true,
				Transcription:  false,
			},
		},
	}
	if _, err := c.Events().Update(context.Background(), "abc-123", "event-123", req, UpdateEventParams{CalendarID: "calendar-123"}); err != nil {
		t.Fatalf("Update with notetaker error: %v", err)
	}
}

func TestEventConferencingEdgeCases(t *testing.T) {
	// Empty conferencing -> nil
	{
		var ev models.Event
		if err := json.Unmarshal([]byte(`{
			"id":"test-event-id","grant_id":"g","calendar_id":"c","busy":true,
			"participants":[{"email":"test@example.com","name":"Test User","status":"yes"}],
		     "when":{"object":"timespan","start_time":1497916800,"end_time":1497920400},
			"conferencing":{}, "title":"X"
		}`), &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Conferencing != nil {
			t.Fatalf("expected nil conferencing, got %#v", ev.Conferencing)
		}
	}
	// Details without provider -> nil
	{
		var ev models.Event
		if err := json.Unmarshal([]byte(`{
			"id":"test-event-id","grant_id":"g","calendar_id":"c","busy":true,
			"participants":[{"email":"test@example.com","name":"Test User","status":"yes"}],
			"when":{"object":"timespan","start_time":1497916800,"end_time":1497920400},
			"conferencing":{"details":{"meeting_code":"code","password":"p","url":"https://x"}}, "title":"Y"
		}`), &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Conferencing != nil {
			t.Fatalf("expected nil conferencing, got %#v", ev.Conferencing)
		}
	}
	// Autocreate without provider -> nil
	{
		var ev models.Event
		if err := json.Unmarshal([]byte(`{
			"id":"test-event-id","grant_id":"g","calendar_id":"c","busy":true,
			"participants":[{"email":"test@example.com","name":"Test User","status":"yes"}],
			"when":{"object":"timespan","start_time":1497916800,"end_time":1497920400},
			"conferencing":{"autocreate":{}}, "title":"Z"
		}`), &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Conferencing != nil {
			t.Fatalf("expected nil conferencing, got %#v", ev.Conferencing)
		}
	}
	// Unknown fields only -> nil
	{
		var ev models.Event
		if err := json.Unmarshal([]byte(`{
			"id":"test-event-id","grant_id":"g","calendar_id":"c","busy":true,
			"participants":[{"email":"test@example.com","name":"Test User","status":"yes"}],
			"when":{"object":"timespan","start_time":1497916800,"end_time":1497920400},
			"conferencing":{"unknown_field":"value"}, "title":"Q"
		}`), &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Conferencing != nil {
			t.Fatalf("expected nil conferencing, got %#v", ev.Conferencing)
		}
	}
}
