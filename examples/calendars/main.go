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
	apiKey := must("NYLAS_API_KEY")
	grant := must("NYLAS_GRANT_ID")
	email := must("EMAIL")

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// 1) List calendars
	cl, err := client.Calendars().List(ctx, grant, nil)
	if err != nil {
		log.Fatalf("list calendars: %v", err)
	}
	if len(cl.Data) == 0 {
		fmt.Println("No calendars found.")
		return
	}

	fmt.Printf("=== Calendars (request-id: %s) ===\n", cl.Headers.Get("x-request-id"))
	for i, c := range cl.Data {
		fmt.Printf("[%d] id=%s name=%q primary=%v readOnly=%v\n",
			i, c.ID, c.Name, b(c.IsPrimary), c.ReadOnly)
	}

	// 2) Choose a calendar: NYLAS_CALENDAR_ID > primary > first
	calID := os.Getenv("NYLAS_CALENDAR_ID")
	if calID == "" {
		if primary := firstPrimary(cl.Data); primary != "" {
			calID = primary
		} else {
			calID = cl.Data[0].ID
		}
	}
	fmt.Println("Using calendar:", calID)

	// 3) Get calendar details
	cresp, err := client.Calendars().Get(ctx, grant, calID, nil)
	if err != nil {
		log.Fatalf("get calendar %q: %v", calID, err)
	}
	fmt.Printf("=== Calendar Detail (request-id: %s) ===\n", cresp.RequestID)
	pp(cresp.Data)

	// 4) Pull events for the calendar (next 7 days)

	evs, err := client.Events().List(ctx, grant, nylas.ListEventsParams{
		CalendarID: calID,
		Start:      ptrInt64(time.Now().Unix()),
		End:        ptrInt64(time.Now().Add(7 * 24 * time.Hour).Unix()),
		Limit:      ptrInt(10),
		// Select:  []string{"id","title","when"},
	})

	if err != nil {
		log.Fatalf("list events: %v", err)
	}
	fmt.Printf("=== Events (count=%d) ===\n", len(evs.Data))
	for i, e := range evs.Data {
		fmt.Printf("[%d] id=%s title=%q when=%s\n", i, e.ID, s(e.Title), summarizeWhenUnion(e.When))
	}

	// 5) (Optional) Free/Busy for this calendar via grant-scoped endpoint
	var emails []string
	emails = append(emails, email)
	fbReq := models.GetFreeBusyRequest{
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(24 * time.Hour).Unix(),
		Emails:    emails,
		// Interval: ptrInt(30), // if supported
	}
	fb, err := client.Calendars().GetFreeBusy(ctx, grant, fbReq)
	if err != nil {
		log.Printf("free-busy error: %v", err)
	} else {
		fmt.Printf("=== Free/Busy (request-id: %s) ===\n", fb.RequestID)
		pp(fb.Data)
	}

	// 6) Create → Get → Delete a demo event on the chosen calendar
	// 6) Create → Get → Delete a demo event on the chosen calendar (union-style When)
	title := "Go SDK Demo Event"
	desc := "Created via nylas-go examples"
	start := time.Now().Add(15 * time.Minute).Unix()
	end := start + 3600

	// Build the union: set exactly one arm.
	// Adjust field names if your Timespan struct uses StartTime/EndTime instead of Start/End.
	when := models.When{
		Timespan: &models.Timespan{
			StartTime: start, // or: StartTime: ptrInt64(start)
			EndTime:   end,   // or: EndTime:   ptrInt64(end)
		},
		// Time: &models.Time{Time: ptrInt64(...)}              // alternative single instant
		// Date: &models.Date{Date: "2025-08-26"}               // alternative all-day single date
		// Datespan: &models.Datespan{StartDate:"2025-08-26", EndDate:"2025-08-27"} // alternative range
	}

	created, err := client.Events().Create(ctx, grant,
		nylas.CreateEventRequest{
			Title:       &title,
			Description: &desc,
			When:        when, // ← union wrapper does the correct JSON
			Participants: []models.Participant{
				{Email: email},
			},
			Busy: ptrBool(true),
		},
		nylas.CreateEventParams{
			CalendarID:         calID,          // required
			NotifyParticipants: ptrBool(false), // no emails for the demo
			TentativeAsBusy:    ptrBool(true),
		},
	)
	if err != nil {
		log.Fatalf("create event: %v", err)
	}
	evtID := created.Data.ID
	fmt.Printf("=== Created Event (id=%s) ===\n", evtID)
	pp(created.Data)

	// GET it back (uses the same calendar)
	fetched, err := client.Events().Get(ctx, grant, evtID, nylas.FindEventParams{
		CalendarID:      calID,
		TentativeAsBusy: ptrBool(true),
	})
	if err != nil {
		log.Fatalf("get event %q: %v", evtID, err)
	}
	fmt.Printf("=== Fetched Event (request-id: %s) ===\n", fetched.Headers.Get("x-request-id"))
	fmt.Printf("id=%s title=%q when=%s busy=%v\n",
		fetched.Data.ID,
		s(fetched.Data.Title),
		summarizeWhenUnion(fetched.Data.When),
		fetched.Data.Busy,
	)

	// DELETE it
	if err := client.Events().Delete(ctx, grant, evtID, nylas.DestroyEventParams{
		CalendarID:         calID,
		NotifyParticipants: ptrBool(false),
	}); err != nil {
		log.Fatalf("delete event %q: %v", evtID, err)
	}
	fmt.Println("=== Deleted Event ===")
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
func pp(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func s(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func b(p *bool) bool          { return p != nil && *p }
func ptrInt(n int) *int       { return &n }
func ptrBool(b bool) *bool    { return &b }
func ptrInt64(n int64) *int64 { return &n }

func firstPrimary(cals []models.Calendar) string {
	for _, c := range cals {
		if b(c.IsPrimary) {
			return c.ID
		}
	}
	return ""
}

// summarizeWhen is a placeholder; if you have a typed When model, adapt it:
//
//	type When struct{ Object string; StartTime *int64; EndTime *int64; Time *int64; ... }

func summarizeWhenUnion(w models.When) string {
	switch {
	case w.Time != nil:
		if w.Time.Time != 0 {
			return time.Unix(w.Time.Time, 0).Format(time.RFC3339)
		}
	case w.Timespan != nil:
		var a, b string
		if w.Timespan.StartTime != 0 {
			a = time.Unix(w.Timespan.StartTime, 0).Format(time.RFC3339)
		}
		if w.Timespan.EndTime != 0 {
			b = time.Unix(w.Timespan.EndTime, 0).Format(time.RFC3339)
		}
		if a != "" || b != "" {
			if b != "" {
				return a + " → " + b
			}
			return a
		}
	case w.Date != nil:
		if w.Date.Date != "" {
			return w.Date.Date + " (all day)"
		}
	case w.Datespan != nil:
		if w.Datespan.StartDate != "" || w.Datespan.EndDate != "" {
			if w.Datespan.EndDate != "" {
				return w.Datespan.StartDate + " → " + w.Datespan.EndDate + " (all day)"
			}
			return w.Datespan.StartDate + " (all day)"
		}
	}
	return "<when>"
}
