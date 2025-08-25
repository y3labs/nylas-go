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
