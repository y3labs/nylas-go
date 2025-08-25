package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/y3labs/nylas-go/nylas"
	"github.com/y3labs/nylas-go/nylas/models"
)

func main() {
	// Env you’ll need:
	//   NYLAS_API_KEY, NYLAS_GRANT_ID
	//   HOST_EMAIL, HOST_CALENDAR_ID  (host/organizer)
	//   GUEST_EMAIL                   (the guest to book)
	//   TIMEZONE                      (e.g., "America/New_York")
	apiKey := must("NYLAS_API_KEY")
	grant := must("NYLAS_GRANT_ID")

	hostEmail := must("EMAIL")
	//hostCalID := must("HOST_CALENDAR_ID")
	guestEmail := getenvDefault("GUEST_EMAIL", "guest@example.com")
	tz := getenvDefault("TIMEZONE", "America/New_York")

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// -------------------------------------------------------------------
	// Pre ) Get Calendar ID if Not Provided in ENV
	// -------------------------------------------------------------------
	hostCalID := os.Getenv("HOST_CALENDAR_ID")
	if hostCalID == "" {
		cl, err := client.Calendars().List(ctx, grant, nil)
		if err != nil {
			log.Fatalf("list calendars: %v", err)
		}
		if len(cl.Data) == 0 {
			fmt.Println("No calendars found.")
			return
		}
		if primary := firstPrimary(cl.Data); primary != "" {
			hostCalID = primary
		} else {
			hostCalID = cl.Data[0].ID
		}
	}
	fmt.Println("Using calendar:", hostCalID)
	// --------------------------------------------------------------------
	// 1) Create a minimal Configuration (one host with one calendar)
	// --------------------------------------------------------------------
	cfgReq := models.CreateConfigurationRequest{
		Participants: []models.ConfigParticipant{
			{
				Email: hostEmail,
				Availability: models.ParticipantAvailability{
					CalendarIDs: []string{hostCalID},
					// OpenHours: []models.OpenHours{...} // optional
				},
				Booking: models.ParticipantBooking{
					CalendarID: hostCalID,
				},
				Timezone:    strPtr(tz),
				IsOrganizer: ptrBool(true), // optional; host is organizer by default
			},
		},
		Availability: models.Availability{
			DurationMinutes: 30, // 30-min slots
			// IntervalMinutes: ptrInt(15), // optional slot interval
		},
		EventBooking: models.EventBooking{
			Title:         "Consultation",
			Timezone:      strPtr(tz),
			BookingType:   bookingTypePtr(models.BookingTypeBooking), // immediate confirmation (no organizer confirm flow)
			DisableEmails: ptrBool(false),                            // let Nylas send email if configured for your tenant
			Reminders: []models.BookingReminder{
				{Type: models.BookingReminderEmail, MinutesBeforeEvent: 30, Recipient: recipientPtr(models.BookingRecipientGuest)},
			},
		},
		Slug: strPtr("go-sdk-demo"), // optional
		Scheduler: &models.SchedulerSettings{
			AdditionalFields: map[string]models.AdditionalField{
				"company": {
					Label:    "Company",
					Type:     models.AdditionalFieldText, // text | email | phone_number | dropdown | ...
					Required: false,
				},
			},
		},
	}

	cfg, err := client.Scheduler().Configurations().Create(ctx, grant, cfgReq)
	die(err)
	fmt.Printf("✓ Configuration created (id=%s, request-id=%s)\n", cfg.Data.ID, cfg.Headers.Get("x-request-id"))
	pp(cfg.Data)

	// --------------------------------------------------------------------
	// 2) (Optional) Create a short-lived Session if the config requires it
	// --------------------------------------------------------------------
	var sess *nylas.Response[models.Session]
	if b(cfg.Data.RequiresSessionAuth) {
		ttlMinutes := 25 // <= 30
		s, err := client.Scheduler().Sessions().Create(ctx, models.CreateSessionRequest{
			ConfigurationID: &cfg.Data.ID,
			TimeToLive:      &ttlMinutes,
		})
		if err != nil {
			fmt.Println("(session create failed):", err)
		} else {
			sess = s
			fmt.Printf("\n✓ Session created (session_id=%s, request-id=%s)\n",
				sess.Data.SessionID, sess.Headers.Get("x-request-id"))
		}
	} else {
		fmt.Println("(skipping session: configuration does not require session auth)")
	}

	// --------------------------------------------------------------------
	// 3) Create a booking precisely at 12:00 PM America/New_York
	// --------------------------------------------------------------------

	// Pull rules from the created configuration
	noticeMin := 60
	if cfg.Data.Scheduler != nil && cfg.Data.Scheduler.MinBookingNotice != nil {
		noticeMin = *cfg.Data.Scheduler.MinBookingNotice
	}
	slotMin := 30
	if cfg.Data.Availability.DurationMinutes != 0 {
		slotMin = cfg.Data.Availability.DurationMinutes
	}
	intervalMin := slotMin
	if cfg.Data.Availability.IntervalMinutes != nil && *cfg.Data.Availability.IntervalMinutes > 0 {
		intervalMin = *cfg.Data.Availability.IntervalMinutes
	}

	// Timezone location for wall-clock math
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("invalid TIMEZONE %q: %v", tz, err)
	}

	// Helper: return the next date/time at given H:M in tz that is >= now + notice
	nextWallclock := func(hour, minute int, noticeMinutes int) time.Time {
		now := time.Now().In(loc)
		minAllowed := now.Add(time.Duration(noticeMinutes) * time.Minute)

		// Start with "today at hour:minute"
		cand := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
		if cand.Before(minAllowed) {
			// move to tomorrow at the same wall-clock time
			cand = cand.Add(24 * time.Hour)
			// If still before minAllowed (e.g., very large notice), keep adding days
			for cand.Before(minAllowed) {
				cand = cand.Add(24 * time.Hour)
			}
		}

		// Align to interval boundary just in case interval != 30 and someone changed it
		if over := cand.Minute() % intervalMin; over != 0 {
			cand = cand.Add(time.Duration(intervalMin-over) * time.Minute)
		}
		return cand
	}

	// 12:00 PM start and end (slotMin long)
	noon := nextWallclock(12, 0, noticeMin+5) // +5 min buffer to avoid edge rejections
	start := noon.Unix()
	end := start + int64(slotMin*60)

	// Build and create the booking
	bookingReq := models.CreateBookingRequest{
		StartTime: start,
		EndTime:   end,
		Guest:     models.BookingGuest{Email: guestEmail, Name: "Guest User"},
		Participants: []models.BookingParticipant{
			{Email: hostEmail},
		},
		Timezone:      &tz,
		EmailLanguage: emailLangPtr(models.EmailLangEN),
		AdditionalFields: map[string]string{
			"company": "Y3 Labs",
		},
	}

	booking, err := client.Scheduler().Bookings().Create(ctx, bookingReq, &models.CreateBookingQueryParams{
		ConfigurationID: &cfg.Data.ID,
		Timezone:        &tz,
	})
	die(err)
	fmt.Printf("\n✓ Booking created for %s (id=%s, status=%s, request-id=%s)\n",
		noon.Format(time.RFC1123Z), booking.Data.BookingID, booking.Data.Status, booking.Headers.Get("x-request-id"))
	pp(booking.Data)

	// --------------------------------------------------------------------
	// 4) Reschedule that booking to 2:00 PM America/New_York (same date)
	// --------------------------------------------------------------------

	// Take the date we actually booked (not "today"), then set to 14:00
	bookedDate := time.Unix(start, 0).In(loc)
	twoPM := time.Date(bookedDate.Year(), bookedDate.Month(), bookedDate.Day(), 14, 0, 0, 0, loc)

	// Ensure we still respect min-notice relative to "now"
	minAllowed := time.Now().In(loc).Add(time.Duration(noticeMin+5) * time.Minute)
	if twoPM.Before(minAllowed) {
		// If somehow 2pm is now too soon, jump to the next day at 2pm
		twoPM = twoPM.Add(24 * time.Hour)
		// keep jumping days if your notice is very large
		for twoPM.Before(minAllowed) {
			twoPM = twoPM.Add(24 * time.Hour)
		}
	}
	// Align to interval boundary as a safety
	if over := twoPM.Minute() % intervalMin; over != 0 {
		twoPM = twoPM.Add(time.Duration(intervalMin-over) * time.Minute)
	}

	newStart := twoPM.Unix()
	newEnd := newStart + int64(slotMin*60)

	res, err := client.Scheduler().Bookings().Reschedule(ctx, booking.Data.BookingID, models.RescheduleBookingRequest{
		StartTime: newStart,
		EndTime:   newEnd,
	}, &models.RescheduleBookingQueryParams{
		ConfigurationID: &cfg.Data.ID,
	})
	die(err)
	fmt.Printf("\n✓ Booking rescheduled to %s (id=%s, status=%s, request-id=%s)\n",
		twoPM.Format(time.RFC1123Z), res.Data.BookingID, res.Data.Status, res.Headers.Get("x-request-id"))
	pp(res.Data)

	// --------------------------------------------------------------------
	// 5) Fetch the booking (Find)
	// --------------------------------------------------------------------
	found, err := client.Scheduler().Bookings().Find(ctx, booking.Data.BookingID, &models.FindBookingQueryParams{
		ConfigurationID: &cfg.Data.ID,
	})
	die(err)
	fmt.Printf("\n✓ Booking fetched (id=%s, status=%s, request-id=%s)\n",
		found.Data.BookingID, found.Data.Status, found.Headers.Get("x-request-id"))
	pp(found.Data)

	// --------------------------------------------------------------------
	// 6) Cancel the booking (Destroy with body)
	// --------------------------------------------------------------------
	reason := "Demo cleanup from Go example"
	cancel, err := client.Scheduler().Bookings().Destroy(ctx, booking.Data.BookingID, models.DeleteBookingRequest{
		CancellationReason: &reason,
	}, &models.DestroyBookingQueryParams{
		ConfigurationID: &cfg.Data.ID,
	})
	die(err)
	fmt.Printf("\n✓ Booking cancelled (request-id=%s)\n", cancel.Headers.Get("x-request-id"))

	// --------------------------------------------------------------------
	// 7) Cleanup: delete configuration, destroy session
	// --------------------------------------------------------------------
	delCfg, err := client.Scheduler().Configurations().Destroy(ctx, grant, cfg.Data.ID)
	if err != nil {
		fmt.Println("(config delete failed/skipped):", err)
	} else {
		fmt.Printf("✓ Configuration deleted (request-id=%s)\n", delCfg.Headers.Get("x-request-id"))
	}

	if sess != nil {
		if _, err := client.Scheduler().Sessions().Destroy(ctx, sess.Data.SessionID); err != nil {
			fmt.Println("(session delete failed/skipped):", err)
		} else {
			fmt.Println("✓ Session deleted")
		}
	}
}

/* ---------- helpers ---------- */

func firstPrimary(cals []models.Calendar) string {
	for _, c := range cals {
		if b(c.IsPrimary) {
			return c.ID
		}
	}
	return ""
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func pp(v any)                                                                { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func b(p *bool) bool                                                          { return p != nil && *p }
func ptrBool(b bool) *bool                                                    { return &b }
func ptrInt(n int) *int                                                       { return &n }
func strPtr(s string) *string                                                 { return &s }
func bookingTypePtr(b models.BookingType) *models.BookingType                 { return &b }
func recipientPtr(r models.BookingRecipientType) *models.BookingRecipientType { return &r }
func emailLangPtr(l models.EmailLanguage) *models.EmailLanguage               { return &l }
