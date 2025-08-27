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

/* -------------------- tests -------------------- */

func TestBookingDeserialization(t *testing.T) {
	raw := `{
		"booking_id": "AAAA-BBBB-1111-2222",
		"event_id": "CCCC-DDDD-3333-4444",
		"title": "My test event",
		"organizer": {
			"name": "John Doe",
			"email": "user@example.com"
		},
		"status": "booked",
		"description": "This is an example of a description."
	}`

	var b models.Booking
	if err := json.Unmarshal([]byte(raw), &b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if b.BookingID != "AAAA-BBBB-1111-2222" {
		t.Fatalf("booking_id = %q, want %q", b.BookingID, "AAAA-BBBB-1111-2222")
	}
	if b.EventID != "CCCC-DDDD-3333-4444" {
		t.Fatalf("event_id = %q, want %q", b.EventID, "CCCC-DDDD-3333-4444")
	}
	if b.Title != "My test event" {
		t.Fatalf("title = %q, want %q", b.Title, "My test event")
	}
	if b.Organizer.Email != "user@example.com" {
		t.Fatalf("organizer.email = %q, want %q", b.Organizer.Email, "user@example.com")
	}
	assertStrPtr(t, b.Organizer.Name, "John Doe", "organizer.name")

	// if Status is custom type, the comparison still works (underlying string)
	if b.Status != "booked" {
		t.Fatalf("status = %q, want %q", b.Status, "booked")
	}
	assertStrPtr(t, b.Description, "This is an example of a description.", "description")
}

func TestBookings_Find(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/scheduling/bookings/booking-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"booking_id": "booking-123",
				"event_id":   "event-xyz",
				"title":      "Found booking",
				"organizer":  map[string]any{"email": "user@example.com"},
				"status":     "booked",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	out, err := c.Scheduler().Bookings().Find(context.Background(), "booking-123", nil)
	if err != nil {
		t.Fatalf("Find error: %v", err)
	}
	if out.Data.BookingID != "booking-123" {
		t.Fatalf("Find booking_id = %q, want %q", out.Data.BookingID, "booking-123")
	}
}

func TestBookings_Create(t *testing.T) {
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/scheduling/bookings")
		gotBody = mustReadAll(t, r.Body)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"booking_id": "new-booking",
				"event_id":   "evt-1",
				"title":      "Created",
				"organizer":  map[string]any{"email": "user@gmail.com", "name": "TEST"},
				"status":     "booked",
			},
		})
	}))
	defer ts.Close()

	var participants []models.BookingParticipant
	participant := models.BookingParticipant{
		Email: "test@nylas.com",
	}
	participants = append(participants, participant)
	c := newTestClient(ts.URL, "test-key")
	req := models.CreateBookingRequest{
		StartTime:    1730725200,
		EndTime:      1730727000,
		Participants: participants,
		Guest: models.BookingGuest{
			Name:  "TEST",
			Email: "user@gmail.com",
		},
	}
	out, err := c.Scheduler().Bookings().Create(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if out.Data.BookingID != "new-booking" {
		t.Fatalf("Create booking_id = %q, want %q", out.Data.BookingID, "new-booking")
	}

	// sanity check the JSON body has the keys from the request
	body := string(gotBody)
	for _, want := range []string{`"start_time":1730725200`, `"end_time":1730727000`, `"participants"`, `"guest"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("request body missing %s: %s", want, body)
		}
	}
}

func TestBookings_Confirm(t *testing.T) {
	var gotBody map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/scheduling/bookings/booking-123")
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"booking_id": "booking-123",
				"status":     "cancelled",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.ConfirmBookingRequest{
		Salt:   "_zfg12it",
		Status: "cancelled",
	}
	out, err := c.Scheduler().Bookings().Confirm(context.Background(), "booking-123", req, nil)
	if err != nil {
		t.Fatalf("Confirm error: %v", err)
	}
	if out.Data.Status != "cancelled" {
		t.Fatalf("Confirm status = %q, want %q", out.Data.Status, "cancelled")
	}
	if gotBody["salt"] != "_zfg12it" || gotBody["status"] != "cancelled" {
		t.Fatalf("confirm body = %#v", gotBody)
	}
}

func TestBookings_Reschedule(t *testing.T) {
	var gotBody map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, "PATCH", "/v3/scheduling/bookings/booking-123")
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"booking_id": "booking-123",
				"status":     "booked",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.RescheduleBookingRequest{
		StartTime: 1730725200,
		EndTime:   1730727000,
	}
	out, err := c.Scheduler().Bookings().Reschedule(context.Background(), "booking-123", req, nil)
	if err != nil {
		t.Fatalf("Reschedule error: %v", err)
	}
	if out.Data.BookingID != "booking-123" {
		t.Fatalf("Reschedule booking_id = %q, want %q", out.Data.BookingID, "booking-123")
	}
	if gotBody["start_time"] != float64(1730725200) || gotBody["end_time"] != float64(1730727000) {
		t.Fatalf("reschedule body = %#v", gotBody)
	}
}

func TestBookings_Destroy(t *testing.T) {
	var gotBody map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/scheduling/bookings/booking-123")
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	req := models.DeleteBookingRequest{
		CancellationReason: strptr("I am no longer available at this time."),
	}
	out, err := c.Scheduler().Bookings().Destroy(context.Background(), "booking-123", req, nil)
	if err != nil {
		t.Fatalf("Destroy error: %v", err)
	}
	if out.RequestID == "" {
		t.Fatalf("Destroy response missing request_id")
	}
	if gotBody["cancellation_reason"] != "I am no longer available at this time." {
		t.Fatalf("destroy body = %#v", gotBody)
	}
}

func TestBookings_Find_ErrorPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-1" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		// Non-2xx → DoJSON returns error
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"type": "invalid_request", "message": "boom"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	_, err := r.Find(context.Background(), "bk-1", &models.FindBookingQueryParams{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestBookings_Find_HeaderRequestIDFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-2" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		// No request_id in JSON; present only in header
		w.Header().Set("X-Request-Id", "rid-find")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	res, err := r.Find(context.Background(), "bk-2", &models.FindBookingQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-find" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
	if res.Headers.Get("X-Request-Id") != "rid-find" {
		t.Fatalf("headers not propagated")
	}
}

func TestBookings_Create_ErrorPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/scheduling/bookings" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"type": "server_error", "message": "nope"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	_, err := r.Create(context.Background(), models.CreateBookingRequest{}, &models.CreateBookingQueryParams{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestBookings_Create_HeaderRequestIDFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/scheduling/bookings" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-create")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	res, err := r.Create(context.Background(), models.CreateBookingRequest{}, &models.CreateBookingQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-create" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
}

func TestBookings_Confirm_ErrorPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-3" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"type": "forbidden", "message": "deny"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	_, err := r.Confirm(context.Background(), "bk-3", models.ConfirmBookingRequest{}, &models.ConfirmBookingQueryParams{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestBookings_Confirm_HeaderRequestIDFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-4" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-confirm")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	res, err := r.Confirm(context.Background(), "bk-4", models.ConfirmBookingRequest{}, &models.ConfirmBookingQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-confirm" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
}

func TestBookings_Reschedule_ErrorPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-5" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"type": "not_found", "message": "missing"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	_, err := r.Reschedule(context.Background(), "bk-5", models.RescheduleBookingRequest{}, &models.RescheduleBookingQueryParams{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestBookings_Reschedule_HeaderRequestIDFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-6" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-resched")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	res, err := r.Reschedule(context.Background(), "bk-6", models.RescheduleBookingRequest{}, &models.RescheduleBookingQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.RequestID != "rid-resched" {
		t.Fatalf("request id fallback failed: %#v", res)
	}
}

func TestBookings_Destroy_ErrorPath(t *testing.T) {
	// DELETE returns error; Destroy should forward it and return (nil, err).
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-7" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"type": "invalid", "message": "no"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	out, err := r.Destroy(context.Background(), "bk-7", models.DeleteBookingRequest{}, &models.DestroyBookingQueryParams{})
	if err == nil || out != nil {
		t.Fatalf("expected error and nil out, got out=%#v err=%v", out, err)
	}
	if _, ok := IsAPIError(err); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestBookings_Destroy_HeaderRequestIDFallback(t *testing.T) {
	// 200 with no request_id in JSON; header should populate DeleteResponse.RequestID.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.EscapedPath() != "/v3/scheduling/bookings/bk-8" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.Header().Set("X-Request-Id", "rid-destroy")
		_ = json.NewEncoder(w).Encode(map[string]any{}) // minimal body
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c: c}
	out, err := r.Destroy(context.Background(), "bk-8", models.DeleteBookingRequest{}, &models.DestroyBookingQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || out.RequestID != "rid-destroy" {
		t.Fatalf("request id fallback failed: %#v", out)
	}
	if out.Headers.Get("X-Request-Id") != "rid-destroy" {
		t.Fatalf("headers not propagated")
	}
}

func TestBookings_Create_ErrorPath2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/v3/scheduling/bookings" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.EscapedPath())
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "rid-err",
			"error":      map[string]any{"type": "server_error", "message": "boom"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c}
	_, err := r.Create(context.Background(), models.CreateBookingRequest{}, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !IsStatus(err, 500) {
		t.Fatalf("want 500 status, got %v", err)
	}
}

func TestBookings_Create_HeaderFallback_And_Query(t *testing.T) {
	type Q = models.CreateBookingQueryParams
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify we actually encoded query (adjust keys to match your params)
		if got := r.URL.Query().Get("timezone"); got != "America/NewYork" {
			t.Fatalf("query not encoded: %q", r.URL.String())
		}
		w.Header().Set("x-request-id", "rid-created")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "bk-1"},
		})
	}))
	defer ts.Close()

	c := NewClient("test-key", WithServerURL(ts.URL))
	r := &BookingsResource{c}
	out, err := r.Create(context.Background(),
		models.CreateBookingRequest{}, &Q{Timezone: StringPtr("America/NewYork")}, // change to a real field present in your Q
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out.RequestID != "rid-created" {
		t.Fatalf("missing header fallback: %#v", out)
	}
}
