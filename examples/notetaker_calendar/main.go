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

	// MEETING_LINK is optional. If empty, we’ll use the Event.Location we set below.
	// You can also point this at an actual Meet/Zoom/Teams URL.
	meetingLink := os.Getenv("MEETING_LINK")

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// 1) Choose a calendar (env > primary > first)
	calID := os.Getenv("NYLAS_CALENDAR_ID")
	if calID == "" {
		cals, err := client.Calendars().List(ctx, grant, nil)
		if err != nil {
			log.Fatalf("list calendars: %v", err)
		}
		if len(cals.Data) == 0 {
			log.Fatal("no calendars available")
		}
		calID = firstPrimary(cals.Data)
		if calID == "" {
			calID = cals.Data[0].ID
		}
	}
	fmt.Println("Using calendar:", calID)

	// 2) Create a 1-hour event starting 20 minutes from now.
	start := time.Now().Add(20 * time.Minute).Unix()
	end := start + 3600
	title := "Meeting with Notetaker (Go SDK)"
	location := meetingLink
	if location == "" {
		location = "https://meet.example.com/abc-123" // demo link if none provided
	}

	when := models.When{
		Timespan: &models.Timespan{
			StartTime: start,
			EndTime:   end,
		},
	}

	created, err := client.Events().Create(ctx, grant,
		nylas.CreateEventRequest{
			Title:       &title,
			When:        when,
			Location:    &location, // stash the meeting link here for notetaker
			Busy:        ptrBool(true),
			Description: strPtr("Created via nylas-go notetaker_calendar example"),
		},
		nylas.CreateEventParams{
			CalendarID:         calID,
			NotifyParticipants: ptrBool(false),
		},
	)
	if err != nil {
		log.Fatalf("create event: %v", err)
	}
	ev := created.Data
	fmt.Printf("✓ Event created (id=%s)\n", ev.ID)
	pp(ev)

	// 3) Invite a notetaker to the event’s meeting link
	join := start // have the notetaker join at event start
	ntName := "Demo Notetaker for " + ev.ID

	inv, err := client.Notetakers().Invite(ctx, models.InviteNotetakerRequest{
		MeetingLink: locationFromEvent(ev),
		JoinTime:    &join,
		Name:        &ntName,
		MeetingSettings: &models.NotetakerMeetingSettingsRequest{
			Transcription:  ptrBool(true),
			AudioRecording: ptrBool(true),
			VideoRecording: ptrBool(false),
		},
	}, grant)
	if err != nil {
		log.Fatalf("invite notetaker: %v", err)
	}
	nt := inv.Data
	fmt.Printf("\n✓ Notetaker invited (id=%s, state=%s)\n", nt.ID, nt.State)
	pp(nt)

	// 4) (Optional) Fetch notetaker media (recording/transcript) — likely not ready if event is in the future
	media, err := client.Notetakers().GetMedia(ctx, nt.ID, grant)
	if err != nil {
		fmt.Printf("\n(media not ready or not available yet): %v\n", err)
	} else {
		fmt.Printf("\n✓ Notetaker media\n")
		pp(media.Data)
	}

	// 5) Cleanup: cancel notetaker, then delete event
	delNT, err := client.Notetakers().Cancel(ctx, nt.ID, grant)
	if err != nil {
		log.Fatalf("cancel notetaker: %v", err)
	}
	fmt.Printf("\n✓ Canceled notetaker (request-id=%s)\n", delNT.Headers.Get("x-request-id"))

	if err := client.Events().Delete(ctx, grant, ev.ID, nylas.DestroyEventParams{
		CalendarID:         calID,
		NotifyParticipants: ptrBool(false),
	}); err != nil {
		log.Fatalf("delete event %q: %v", ev.ID, err)
	}
	fmt.Println("✓ Deleted event")
}

/* ---------- helpers ---------- */

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}

func ptrBool(b bool) *bool    { return &b }
func strPtr(s string) *string { return &s }

func pp(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }

func firstPrimary(cals []models.Calendar) string {
	for _, c := range cals {
		if c.IsPrimary != nil && *c.IsPrimary {
			return c.ID
		}
	}
	return ""
}

// If your Event model also includes Conferencing with a URL, you can prefer it here.
// This helper currently returns Location, which we set when creating the event.
func locationFromEvent(e models.Event) string {
	if e.Location != nil && *e.Location != "" {
		return *e.Location
	}
	// Uncomment if your model has conferencing details:
	// if e.Conferencing != nil && e.Conferencing.Details != nil && e.Conferencing.Details.URL != "" {
	//     return e.Conferencing.Details.URL
	// }
	return ""
}
