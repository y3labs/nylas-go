package nylas

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

/* -------------------- helpers -------------------- */

func mustReadAll(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}

func assertBoolPtr(t *testing.T, got *bool, want bool, field string) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %v (nil? %v)", field, boolval(got), want, got == nil)
	}
}

/* -------------------- tests -------------------- */

func TestCalendarDeserialization(t *testing.T) {
	raw := `{
		"grant_id": "abc-123-grant-id",
		"description": "Description of my new calendar",
		"hex_color": "#039BE5",
		"hex_foreground_color": "#039BE5",
		"id": "5d3qmne77v32r8l4phyuksl2x",
		"is_owned_by_user": true,
		"is_primary": true,
		"location": "Los Angeles, CA",
		"metadata": {"your-key": "value"},
		"name": "My New Calendar",
		"object": "calendar",
		"read_only": false,
		"timezone": "America/Los_Angeles"
	}`

	var cal models.Calendar
	if err := json.Unmarshal([]byte(raw), &cal); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cal.GrantID != "abc-123-grant-id" {
		t.Fatalf("grant_id = %q, want %q", cal.GrantID, "abc-123-grant-id")
	}
	assertStrPtr(t, cal.Description, "Description of my new calendar", "description")
	assertStrPtr(t, cal.HexColor, "#039BE5", "hex_color")
	assertStrPtr(t, cal.HexForegroundColor, "#039BE5", "hex_foreground_color")
	if cal.ID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("id = %q, want %q", cal.ID, "5d3qmne77v32r8l4phyuksl2x")
	}
	if cal.IsOwnedByUser != true {
		t.Fatalf("is_owned_by_user = %v, want %v", cal.IsOwnedByUser, true)
	}
	assertBoolPtr(t, cal.IsPrimary, true, "is_primary")
	assertStrPtr(t, cal.Location, "Los Angeles, CA", "location")
	if cal.Metadata["your-key"] != "value" {
		t.Fatalf("metadata[your-key] = %q, want %q", cal.Metadata["your-key"], "value")
	}
	if cal.Name != "My New Calendar" {
		t.Fatalf("name = %q, want %q", cal.Name, "My New Calendar")
	}
	if cal.Object != "calendar" {
		t.Fatalf("object = %q, want %q", cal.Object, "calendar")
	}
	if cal.ReadOnly != false {
		t.Fatalf("read_only = %v, want %v", cal.ReadOnly, false)
	}
	assertStrPtr(t, cal.Timezone, "America/Los_Angeles", "timezone")
}

func TestCalendarWithNotetakerDeserialization(t *testing.T) {
	raw := `{
		"grant_id": "abc-123-grant-id",
		"description": "Description of my new calendar",
		"id": "5d3qmne77v32r8l4phyuksl2x",
		"is_owned_by_user": true,
		"name": "My New Calendar",
		"object": "calendar",
		"read_only": false,
		"notetaker": {
			"name": "My Notetaker",
			"meeting_settings": {
				"video_recording": true,
				"audio_recording": true,
				"transcription": true
			},
			"rules": {
				"event_selection": ["internal", "external"],
				"participant_filter": {
					"participants_gte": 3,
					"participants_lte": 10
				}
			}
		}
	}`

	var cal models.Calendar
	if err := json.Unmarshal([]byte(raw), &cal); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cal.GrantID != "abc-123-grant-id" {
		t.Fatalf("grant_id = %q, want %q", cal.GrantID, "abc-123-grant-id")
	}
	if cal.ID != "5d3qmne77v32r8l4phyuksl2x" {
		t.Fatalf("id = %q, want %q", cal.ID, "5d3qmne77v32r8l4phyuksl2x")
	}
	if cal.IsOwnedByUser != true {
		t.Fatalf("is_owned_by_user = %v, want true", cal.IsOwnedByUser)
	}
	if cal.Name != "My New Calendar" {
		t.Fatalf("name = %q, want %q", cal.Name, "My New Calendar")
	}
	if cal.Object != "calendar" {
		t.Fatalf("object = %q, want %q", cal.Object, "calendar")
	}
	if cal.ReadOnly != false {
		t.Fatalf("read_only = %v, want false", cal.ReadOnly)
	}
	if cal.Notetaker == nil {
		t.Fatal("notetaker is nil")
	}
	assertStrPtr(t, cal.Notetaker.Name, "My Notetaker", "notetaker.name")
	if cal.Notetaker.MeetingSettings == nil {
		t.Fatal("notetaker.meeting_settings is nil")
	}
	ms := cal.Notetaker.MeetingSettings
	if ms == nil || !ms.VideoRecording || !ms.AudioRecording || !ms.Transcription {
		t.Fatalf("meeting_settings not all true: %#v", cal.Notetaker.MeetingSettings)
	}
	if cal.Notetaker.Rules == nil {
		t.Fatal("notetaker.rules is nil")
	}
	es := cal.Notetaker.Rules.EventSelection
	if len(es) != 2 {
		t.Fatalf("event_selection length = %d, want 2", len(es))
	}
	// membership
	hasInternal := false
	hasExternal := false
	for _, v := range es {
		if v == models.EventSelectionInternal {
			hasInternal = true
		}
		if v == models.EventSelectionExternal {
			hasExternal = true
		}
	}
	if !hasInternal || !hasExternal {
		t.Fatalf("event_selection missing internal/external: %#v", es)
	}
	if cal.Notetaker.Rules.ParticipantFilter == nil {
		t.Fatal("notetaker.rules.participant_filter is nil")
	}
	if cal.Notetaker.Rules.ParticipantFilter.ParticipantsGTE == nil ||
		*cal.Notetaker.Rules.ParticipantFilter.ParticipantsGTE != 3 {
		t.Fatalf("participants_gte = %#v, want 3", cal.Notetaker.Rules.ParticipantFilter.ParticipantsGTE)
	}
	if cal.Notetaker.Rules.ParticipantFilter.ParticipantsLTE == nil ||
		*cal.Notetaker.Rules.ParticipantFilter.ParticipantsLTE != 10 {
		t.Fatalf("participants_lte = %#v, want 10", cal.Notetaker.Rules.ParticipantFilter.ParticipantsLTE)
	}
}

func TestCalendars_List_NoParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/calendars")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Calendars().List(context.Background(), "abc-123", nil); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestCalendars_List_WithQueryParams(t *testing.T) {
	var gotRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/calendars")
		gotRawQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data":       []any{},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	limit := 20
	_, err := c.Calendars().List(context.Background(), "abc-123", &models.ListCalendarsQueryParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	q, _ := url.ParseQuery(gotRawQuery)
	if q.Get("limit") != "20" {
		t.Fatalf("query limit = %q, want %q", q.Get("limit"), "20")
	}
}

func TestCalendars_List_WithSelect(t *testing.T) {
	var gotRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/calendars")
		gotRawQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data": []any{
				map[string]any{
					"id":          "calendar-123",
					"name":        "My Calendar",
					"description": "My calendar description",
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	sel := "id,name,description"
	_, err := c.Calendars().List(context.Background(), "abc-123", &models.ListCalendarsQueryParams{
		Select: &sel,
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	q, _ := url.ParseQuery(gotRawQuery)
	if q.Get("select") != "id,name,description" {
		t.Fatalf("select = %q, want %q", q.Get("select"), "id,name,description")
	}
}

func TestCalendars_Find(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/calendars/calendar-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data": map[string]any{
				"id":   "calendar-123",
				"name": "My Calendar",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Calendars().Get(context.Background(), "abc-123", "calendar-123", nil); err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestCalendars_Find_WithSelect(t *testing.T) {
	var gotRawQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/calendars/calendar-123")
		gotRawQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data": map[string]any{
				"id":          "calendar-123",
				"name":        "My Calendar",
				"description": "My calendar description",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	sel := "id,name,description"
	if _, err := c.Calendars().Get(context.Background(), "abc-123", "calendar-123", &models.FindCalendarQueryParams{
		Select: &sel,
	}); err != nil {
		t.Fatalf("Get error: %v", err)
	}
	q, _ := url.ParseQuery(gotRawQuery)
	if q.Get("select") != "id,name,description" {
		t.Fatalf("select = %q, want %q", q.Get("select"), "id,name,description")
	}
}

func TestCalendars_Create(t *testing.T) {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/calendars")
		gotBody = string(mustReadAll(t, r.Body))
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data": map[string]any{
				"id": "new-cal",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	body := models.CreateCalendarRequest{
		Name:        "My New Calendar",
		Description: strptr("Description of my new calendar"),
		Location:    strptr("Los Angeles, CA"),
		Timezone:    strptr("America/Los_Angeles"),
		Metadata:    map[string]string{"your-key": "value"},
	}
	if _, err := c.Calendars().Create(context.Background(), "abc-123", body); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	for _, want := range []string{
		`"name":"My New Calendar"`,
		`"description":"Description of my new calendar"`,
		`"location":"Los Angeles, CA"`,
		`"timezone":"America/Los_Angeles"`,
		`"metadata":{"your-key":"value"}`,
	} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("create body missing %s: %s", want, gotBody)
		}
	}
}

func TestCalendars_Create_WithNotetaker(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/calendars")
		_ = json.NewDecoder(r.Body).Decode(&got)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc", "data": map[string]any{"id": "cal"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	body := models.CreateCalendarRequest{
		Name:        "My New Calendar",
		Description: strptr("Description of my new calendar"),
		Location:    strptr("Los Angeles, CA"),
		Timezone:    strptr("America/Los_Angeles"),
		Notetaker: &models.NotetakerCalendarRequest{
			Name: strptr("My Notetaker"),
			MeetingSettings: &models.NotetakerCalendarSettings{
				VideoRecording: boolptr(true),
				AudioRecording: boolptr(true),
				Transcription:  boolptr(true),
			},
			Rules: &models.NotetakerCalendarRules{
				EventSelection: []models.EventSelection{
					models.EventSelectionInternal, models.EventSelectionExternal,
				},
				ParticipantFilter: &models.NotetakerCalendarParticipantFilter{
					ParticipantsGTE: intptr(3),
					ParticipantsLTE: intptr(10),
				},
			},
		},
	}
	if _, err := c.Calendars().Create(context.Background(), "abc-123", body); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	nt := got["notetaker"].(map[string]any)
	if nt["name"] != "My Notetaker" {
		t.Fatalf("notetaker.name = %v", nt["name"])
	}
}

func TestCalendars_Update(t *testing.T) {
	var got string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/calendars/calendar-123")
		got = string(mustReadAll(t, r.Body))
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc", "data": map[string]any{"id": "calendar-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	body := models.UpdateCalendarRequest{
		Name:        strptr("My Updated Calendar"),
		Description: strptr("Description of my updated calendar"),
		Location:    strptr("Los Angeles, CA"),
		Timezone:    strptr("America/Los_Angeles"),
		Metadata:    map[string]string{"your-key": "value"},
	}
	if _, err := c.Calendars().Update(context.Background(), "abc-123", "calendar-123", body); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	for _, want := range []string{
		`"name":"My Updated Calendar"`,
		`"description":"Description of my updated calendar"`,
		`"location":"Los Angeles, CA"`,
		`"timezone":"America/Los_Angeles"`,
		`"metadata":{"your-key":"value"}`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("update body missing %s: %s", want, got)
		}
	}
}

func TestCalendars_Update_WithNotetaker(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/grants/abc-123/calendars/calendar-123")
		_ = json.NewDecoder(r.Body).Decode(&got)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc", "data": map[string]any{"id": "calendar-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	body := models.UpdateCalendarRequest{
		Name: strptr("My Updated Calendar"),
		Notetaker: &models.NotetakerCalendarRequest{
			Name: strptr("Updated Notetaker"),
			MeetingSettings: &models.NotetakerCalendarSettings{
				VideoRecording: boolptr(false),
				AudioRecording: boolptr(true),
				Transcription:  boolptr(false),
			},
			Rules: &models.NotetakerCalendarRules{
				EventSelection: []models.EventSelection{models.EventSelectionAll},
				ParticipantFilter: &models.NotetakerCalendarParticipantFilter{
					ParticipantsGTE: intptr(2),
				},
			},
		},
	}
	if _, err := c.Calendars().Update(context.Background(), "abc-123", "calendar-123", body); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	nt := got["notetaker"].(map[string]any)
	if nt["name"] != "Updated Notetaker" {
		t.Fatalf("notetaker.name = %v", nt["name"])
	}
	ms := nt["meeting_settings"].(map[string]any)
	if ms["video_recording"] != false || ms["audio_recording"] != true || ms["transcription"] != false {
		t.Fatalf("meeting_settings = %#v", ms)
	}
}

func TestCalendars_Destroy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/calendars/calendar-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "abc"})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	if _, err := c.Calendars().Delete(context.Background(), "abc-123", "calendar-123"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestCalendars_GetAvailability(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/calendars/availability")
		_ = json.NewDecoder(r.Body).Decode(&got)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data":       map[string]any{"ok": true},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	// Build a minimal availability request compatible with your models.
	body := models.GetAvailabilityRequest{
		StartTime:       1497916800,
		EndTime:         1498003200,
		DurationMinutes: 60,
		IntervalMinutes: intptr(30),
		RoundTo:         nil,
		AvailabilityRules: &models.AvailabilityRules{
			AvailabilityMethod: nil,
			Buffer:             nil,
			DefaultOpenHours: []models.OpenHours{
				{
					Days:     []int{0},
					Timezone: "America/Los_Angeles",
					Start:    "09:00",
					End:      "17:00",
					Exdates:  []string{"2021-03-01"},
				},
			},
			RoundRobinGroupID: nil,
			TentativeAsBusy:   boolptr(false),
		},
	}
	if _, err := c.Calendars().GetAvailability(context.Background(), body); err != nil {
		t.Fatalf("GetAvailability error: %v", err)
	}
	// spot check keys made it through
	if got["start_time"] == nil || got["end_time"] == nil || got["duration_minutes"] == nil {
		t.Fatalf("availability body missing fields: %#v", got)
	}
}

func TestCalendars_GetFreeBusy(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc123/calendars/free-busy")
		_ = json.NewDecoder(r.Body).Decode(&got)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc",
			"data": []any{
				map[string]any{
					"email": "user1@example.com",
					"time_slots": []any{
						map[string]any{"start_time": 1690898400, "end_time": 1690902000, "status": "busy", "object": "time_slot"},
					},
					"object": "free_busy",
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.GetFreeBusyRequest{
		Emails:    []string{"test@gmail.com", "test2@gmail.com"},
		StartTime: 1497916800,
		EndTime:   1498003200,
	}
	if _, err := c.Calendars().GetFreeBusy(context.Background(), "abc123", req); err != nil {
		t.Fatalf("GetFreeBusy error: %v", err)
	}
	// ensure body echoed
	if got["emails"] == nil {
		t.Fatalf("free-busy body missing emails: %#v", got)
	}
}
