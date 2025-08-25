package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestNotetakerDeserialization(t *testing.T) {
	js := `{
		"id": "notetaker-123",
		"name": "Nylas Notetaker",
		"join_time": 1656090000,
		"meeting_link": "https://meet.google.com/abc-def-ghi",
		"meeting_provider": "Google Meet",
		"state": "scheduled",
		"object": "notetaker",
		"meeting_settings": {
			"video_recording": true,
			"audio_recording": true,
			"transcription": true
		}
	}`
	var nt models.Notetaker
	if err := json.Unmarshal([]byte(js), &nt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if nt.ID != "notetaker-123" || nt.Name != "Nylas Notetaker" {
		t.Fatalf("id/name mismatch: %#v", nt)
	}
	if nt.JoinTime != 1656090000 {
		t.Fatalf("join_time = %d", nt.JoinTime)
	}
	if nt.MeetingLink != "https://meet.google.com/abc-def-ghi" {
		t.Fatalf("meeting_link = %q", nt.MeetingLink)
	}
	if nt.MeetingProvider == nil || *nt.MeetingProvider != models.MeetingProviderGoogleMeet {
		t.Fatalf("meeting_provider = %#v", nt.MeetingProvider)
	}
	if nt.State != models.NotetakerStateScheduled {
		t.Fatalf("state = %q", nt.State)
	}
	if nt.Object != "notetaker" {
		t.Fatalf("object = %q", nt.Object)
	}
	ms := nt.MeetingSettings
	if !ms.VideoRecording || !ms.AudioRecording || !ms.Transcription {
		t.Fatalf("meeting_settings not all true: %#v", ms)
	}
}

func TestNotetakerStateEnum_AllValues(t *testing.T) {
	cases := []struct {
		s string
		e models.NotetakerState
	}{
		{"scheduled", models.NotetakerStateScheduled},
		{"connecting", models.NotetakerStateConnecting},
		{"waiting_for_entry", models.NotetakerStateWaitingForEntry},
		{"failed_entry", models.NotetakerStateFailedEntry},
		{"attending", models.NotetakerStateAttending},
		{"media_processing", models.NotetakerStateMediaProcessing},
		{"media_available", models.NotetakerStateMediaAvailable},
		{"media_error", models.NotetakerStateMediaError},
		{"media_deleted", models.NotetakerStateMediaDeleted},
	}
	for _, tc := range cases {
		js := `{
			"id":"n1","name":"Nylas Notetaker",
			"join_time":1656090000,"meeting_link":"https://meet",
			"state":"` + tc.s + `",
			"meeting_settings":{"video_recording":true,"audio_recording":true,"transcription":true}
		}`
		var nt models.Notetaker
		if err := json.Unmarshal([]byte(js), &nt); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if nt.State != tc.e {
			t.Fatalf("state got %q want %q", nt.State, tc.e)
		}
		if string(nt.State) != tc.s {
			t.Fatalf("string(state) got %q want %q", string(nt.State), tc.s)
		}
	}
}

func TestListNotetakers_WithAndWithoutIdentifier(t *testing.T) {
	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": []any{}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().List(context.Background(), "abc-123", nil); err != nil {
		t.Fatalf("List with identifier error: %v", err)
	}

	// Without identifier (global)
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/notetakers")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": []any{}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().List(context.Background(), "", nil); err != nil {
		t.Fatalf("List global error: %v", err)
	}
}

func TestListNotetakers_WithQueryParams_EnumAndLimit(t *testing.T) {
	state := models.NotetakerStateScheduled
	limit := 20
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("state") != string(state) || q.Get("limit") != "20" {
			t.Fatalf("query mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "k")
	_, err := c.Notetakers().List(context.Background(), "abc-123", &models.ListNotetakerQueryParams{
		State: &state, Limit: intptr(limit),
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindNotetaker_WithAndWithoutIdentifier(t *testing.T) {
	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers/notetaker-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().Get(context.Background(), "notetaker-123", "abc-123", nil); err != nil {
		t.Fatalf("Get with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/notetakers/notetaker-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().Get(context.Background(), "notetaker-123", "", nil); err != nil {
		t.Fatalf("Get global error: %v", err)
	}
}

func TestInviteNotetaker_WithAndWithoutIdentifier(t *testing.T) {
	req := models.InviteNotetakerRequest{
		MeetingLink: "https://meet.google.com/abc-def-ghi",
		JoinTime:    i64ptr(1656090000),
		Name:        strptr("Custom Notetaker"),
		MeetingSettings: &models.NotetakerMeetingSettingsRequest{
			VideoRecording: boolptr(true),
			AudioRecording: boolptr(true),
			Transcription:  boolptr(true),
		},
	}

	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/notetakers")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": map[string]any{"id": "n1"}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().Invite(context.Background(), req, "abc-123"); err != nil {
		t.Fatalf("Invite with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/notetakers")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": map[string]any{"id": "n2"}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().Invite(context.Background(), req, ""); err != nil {
		t.Fatalf("Invite global error: %v", err)
	}
}

func TestUpdateNotetaker_WithAndWithoutIdentifier(t *testing.T) {
	req := models.UpdateNotetakerRequest{
		Name:     strptr("Updated Notetaker"),
		JoinTime: i64ptr(1656100000),
		MeetingSettings: &models.NotetakerMeetingSettingsRequest{
			VideoRecording: boolptr(false),
			AudioRecording: boolptr(true),
			Transcription:  boolptr(true),
		},
	}

	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/grants/abc-123/notetakers/notetaker-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().Update(context.Background(), "notetaker-123", req, "abc-123"); err != nil {
		t.Fatalf("Update with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPatch, "/v3/notetakers/notetaker-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().Update(context.Background(), "notetaker-123", req, ""); err != nil {
		t.Fatalf("Update global error: %v", err)
	}
}

func TestLeaveMeeting_WithAndWithoutIdentifier(t *testing.T) {
	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/grants/abc-123/notetakers/notetaker-123/leave")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().Leave(context.Background(), "notetaker-123", "abc-123"); err != nil {
		t.Fatalf("Leave with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/notetakers/notetaker-123/leave")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().Leave(context.Background(), "notetaker-123", ""); err != nil {
		t.Fatalf("Leave global error: %v", err)
	}
}

func TestGetMedia_WithAndWithoutIdentifier(t *testing.T) {
	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers/notetaker-123/media")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1", "data": map[string]any{"recording": map[string]any{"url": "u"}}})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().GetMedia(context.Background(), "notetaker-123", "abc-123"); err != nil {
		t.Fatalf("GetMedia with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/notetakers/notetaker-123/media")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2", "data": map[string]any{"transcript": map[string]any{"url": "u2"}}})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().GetMedia(context.Background(), "notetaker-123", ""); err != nil {
		t.Fatalf("GetMedia global error: %v", err)
	}
}

func TestCancelNotetaker_WithAndWithoutIdentifier(t *testing.T) {
	// With identifier
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/grants/abc-123/notetakers/notetaker-123/cancel")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r1"})
	}))
	defer ts1.Close()
	c1 := newTestClient(ts1.URL, "k")
	if _, err := c1.Notetakers().Cancel(context.Background(), "notetaker-123", "abc-123"); err != nil {
		t.Fatalf("Cancel with identifier error: %v", err)
	}

	// Without identifier
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/notetakers/notetaker-123/cancel")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "r2"})
	}))
	defer ts2.Close()
	c2 := newTestClient(ts2.URL, "k")
	if _, err := c2.Notetakers().Cancel(context.Background(), "notetaker-123", ""); err != nil {
		t.Fatalf("Cancel global error: %v", err)
	}
}

func TestNotetakerMediaDeserialization(t *testing.T) {
	js := `{
		"recording": {
			"size": 21550491,
			"name": "meeting_recording.mp4",
			"type": "video/mp4",
			"created_at": 1744222418,
			"expires_at": 1744481618,
			"url": "url_for_recording",
			"ttl": 259106
		},
		"transcript": {
			"size": 862,
			"name": "raw_transcript.json",
			"type": "application/json",
			"created_at": 1744222418,
			"expires_at": 1744481618,
			"url": "url_for_transcript",
			"ttl": 259106
		}
	}`
	var m models.NotetakerMedia
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m.Recording == nil || m.Recording.URL == "" || m.Recording.Name != "meeting_recording.mp4" ||
		m.Recording.Type != "video/mp4" || m.Recording.Size != 21550491 ||
		m.Recording.CreatedAt != 1744222418 || m.Recording.ExpiresAt != 1744481618 || m.Recording.TTL != 259106 {
		t.Fatalf("recording mismatch: %#v", m.Recording)
	}

	if m.Transcript == nil || m.Transcript.URL == "" || m.Transcript.Name != "raw_transcript.json" ||
		m.Transcript.Type != "application/json" || m.Transcript.Size != 862 ||
		m.Transcript.CreatedAt != 1744222418 || m.Transcript.ExpiresAt != 1744481618 || m.Transcript.TTL != 259106 {
		t.Fatalf("transcript mismatch: %#v", m.Transcript)
	}
}

func TestMeetingProviderEnum(t *testing.T) {
	cases := []struct {
		s string
		e models.MeetingProvider
	}{
		{"Google Meet", models.MeetingProviderGoogleMeet},
		{"Zoom Meeting", models.MeetingProviderZoom},
		{"Microsoft Teams", models.MeetingProviderMicrosoftTeams},
	}

	for _, tc := range cases {
		js := `{
			"id":"n1","name":"Nylas Notetaker","join_time":1656090000,"meeting_link":"https://meet",
			"meeting_provider":"` + tc.s + `","state":"scheduled",
			"meeting_settings":{"video_recording":true,"audio_recording":true,"transcription":true}
		}`
		var nt models.Notetaker
		if err := json.Unmarshal([]byte(js), &nt); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if nt.MeetingProvider == nil || *nt.MeetingProvider != tc.e {
			t.Fatalf("meeting_provider got %#v want %q", nt.MeetingProvider, tc.e)
		}
		if nt.MeetingProvider == nil || string(*nt.MeetingProvider) != tc.s {
			t.Fatalf("provider string got %q want %q", string(*nt.MeetingProvider), tc.s)
		}
	}
}

func TestStateEnumComparisonHelpers(t *testing.T) {
	js := `{
		"id":"n1","name":"Nylas Notetaker","join_time":1656090000,
		"meeting_link":"https://meet.google.com/abc-def-ghi","state":"scheduled",
		"meeting_settings":{"video_recording":true,"audio_recording":true,"transcription":true}
	}`
	var nt models.Notetaker
	if err := json.Unmarshal([]byte(js), &nt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !nt.IsState(models.NotetakerStateScheduled) || !nt.IsScheduled() || nt.IsAttending() || nt.HasMediaAvailable() {
		t.Fatalf("state helpers mismatch: %#v", nt.State)
	}
}

func TestHelperMethodsVariousStates(t *testing.T) {
	// Attending
	js1 := `{
		"id":"n2","name":"Nylas Notetaker","join_time":1656090000,"meeting_link":"https://zoom.us/j/123",
		"meeting_provider":"Zoom Meeting","state":"attending",
		"meeting_settings":{"video_recording":true,"audio_recording":true,"transcription":true}
	}`
	var a models.Notetaker
	if err := json.Unmarshal([]byte(js1), &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !a.IsState(models.NotetakerStateAttending) || a.IsScheduled() || !a.IsAttending() || a.HasMediaAvailable() {
		t.Fatalf("attending helpers mismatch")
	}

	// Media available
	js2 := `{
		"id":"n3","name":"Nylas Notetaker","join_time":1656090000,"meeting_link":"https://teams.microsoft.com/...",
		"meeting_provider":"Microsoft Teams","state":"media_available",
		"meeting_settings":{"video_recording":true,"audio_recording":true,"transcription":true}
	}`
	var m models.Notetaker
	if err := json.Unmarshal([]byte(js2), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !m.IsState(models.NotetakerStateMediaAvailable) || m.IsScheduled() || m.IsAttending() || !m.HasMediaAvailable() {
		t.Fatalf("media_available helpers mismatch")
	}
}

func TestListNotetakers_TimeFilters(t *testing.T) {
	start := int64(1704067200) // Jan 1, 2024
	end := int64(1704153600)   // Jan 2, 2024
	limit := 20
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("join_time_start") != "1704067200" || q.Get("join_time_end") != "1704153600" || q.Get("limit") != "20" {
			t.Fatalf("query mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "k")
	_, err := c.Notetakers().List(context.Background(), "abc-123", &models.ListNotetakerQueryParams{
		JoinTimeStart: &start, JoinTimeEnd: &end, Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestNotetakerLeaveResponseDeserialization(t *testing.T) {
	js := `{"id":"notetaker-123","message":"Notetaker has left the meeting","object":"notetaker_leave_response"}`
	var lr models.NotetakerLeaveResponse
	if err := json.Unmarshal([]byte(js), &lr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if lr.ID != "notetaker-123" || lr.Message != "Notetaker has left the meeting" || lr.Object != "notetaker_leave_response" {
		t.Fatalf("leave response mismatch: %#v", lr)
	}
}

func TestListNotetakers_OrderParams(t *testing.T) {
	orderBy := models.NotetakerOrderByName
	orderDir := models.NotetakerOrderDirectionDESC
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("order_by") != string(orderBy) || q.Get("order_direction") != string(orderDir) {
			t.Fatalf("order query mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "k")
	_, err := c.Notetakers().List(context.Background(), "abc-123", &models.ListNotetakerQueryParams{
		OrderBy: &orderBy, OrderDirection: &orderDir,
	})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestListNotetakers_DefaultOrder_NoQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/grants/abc-123/notetakers")
		// ensure no query params
		if raw := r.URL.RawQuery; raw != "" {
			t.Fatalf("unexpected query: %s", raw)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": []any{}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "k")
	if _, err := c.Notetakers().List(context.Background(), "abc-123", nil); err != nil {
		t.Fatalf("List error: %v", err)
	}
}

func TestFindNotetaker_WithSelectParam(t *testing.T) {
	sel := "id,name,join_time"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("select") != sel {
			t.Fatalf("select mismatch: %s", r.URL.RawQuery)
		}
		assertMethodPath(t, r, http.MethodGet, "/v3/notetakers/notetaker-123")
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "rid", "data": map[string]any{"id": "notetaker-123"}})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "k")
	_, err := c.Notetakers().Get(context.Background(), "notetaker-123", "", &models.FindNotetakerQueryParams{Select: &sel})
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

// Ensure url is imported (used above) and not optimized away.
var _ = url.PathEscape
