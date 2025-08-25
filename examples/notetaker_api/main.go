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
	grant := os.Getenv("NYLAS_GRANT_ID") // optional; pass "" for global scope
	link := must("MEETING_LINK")         // e.g., "https://meet.google.com/xxx-xxxx-xxx"

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// 1) Invite a notetaker (joins 5 minutes from now, transcription on)
	join := time.Now().Add(5 * time.Minute).Unix()
	name := "Demo Notetaker"
	it, err := client.Notetakers().Invite(ctx, models.InviteNotetakerRequest{
		MeetingLink: link,
		JoinTime:    &join,
		Name:        &name,
		MeetingSettings: &models.NotetakerMeetingSettingsRequest{
			Transcription:  ptrBool(true),
			AudioRecording: ptrBool(true),
			VideoRecording: ptrBool(false),
		},
	}, grant)
	if err != nil {
		log.Fatalf("invite notetaker: %v", err)
	}
	nt := it.Data
	fmt.Printf("✓ Invited notetaker (id=%s, state=%s, request-id=%s)\n",
		nt.ID, nt.State, it.Headers.Get("x-request-id"))
	pp(nt)

	// 2) List notetakers (optionally filter/sort)
	lst, err := client.Notetakers().List(ctx, grant, &models.ListNotetakerQueryParams{
		Limit:          ptrInt(10),
		OrderBy:        ptrOrderBy(models.NotetakerOrderByCreatedAt),
		OrderDirection: ptrOrderDir(models.NotetakerOrderDirectionDESC),
	})
	if err != nil {
		log.Fatalf("list notetakers: %v", err)
	}
	fmt.Printf("\n=== Notetakers (count=%d, request-id=%s) ===\n", len(lst.Data), lst.Headers.Get("x-request-id"))
	for i, n := range lst.Data {
		fmt.Printf("[%d] id=%s name=%q state=%s join=%d\n", i, n.ID, n.Name, n.State, n.JoinTime)
	}

	// 3) Get the created notetaker explicitly
	got, err := client.Notetakers().Get(ctx, nt.ID, grant, nil)
	if err != nil {
		log.Fatalf("get notetaker %q: %v", nt.ID, err)
	}
	fmt.Printf("\n✓ Fetched notetaker (request-id=%s)\n", got.RequestID)
	pp(got.Data)

	// 4) Update the notetaker (rename)
	newName := "Demo Notetaker (Updated)"
	upd, err := client.Notetakers().Update(ctx, nt.ID, models.UpdateNotetakerRequest{
		Name: &newName,
	}, grant)
	if err != nil {
		log.Fatalf("update notetaker %q: %v", nt.ID, err)
	}
	fmt.Printf("\n✓ Updated notetaker (request-id=%s)\n", upd.RequestID)
	pp(upd.Data)

	// 5) (Optional) Get media (recording/transcript URLs) if already available
	media, err := client.Notetakers().GetMedia(ctx, nt.ID, grant)
	if err != nil {
		fmt.Printf("\n(media not ready yet or not available): %v\n", err)
	} else {
		fmt.Printf("\n✓ Media (request-id=%s)\n", media.RequestID)
		pp(media.Data)
	}

	// 6) Clean up: cancel the notetaker
	del, err := client.Notetakers().Cancel(ctx, nt.ID, grant)
	if err != nil {
		log.Fatalf("cancel notetaker %q: %v", nt.ID, err)
	}
	fmt.Printf("\n✓ Canceled notetaker (request-id=%s)\n", del.Headers.Get("x-request-id"))
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}

func pp(v any)                                                                     { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func ptrBool(b bool) *bool                                                         { return &b }
func ptrInt(n int) *int                                                            { return &n }
func ptrOrderBy(v models.NotetakerOrderBy) *models.NotetakerOrderBy                { return &v }
func ptrOrderDir(v models.NotetakerOrderDirection) *models.NotetakerOrderDirection { return &v }
