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
	rcpt := os.Getenv("EMAIL")
	if rcpt == "" {
		rcpt = "recipient@example.com"
	}

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// -----------------------------
	// 1) Create a draft with metadata
	// -----------------------------
	body := "Hello from Go metadata demo at " + time.Now().Format(time.RFC3339)
	subj := "Metadata Demo (Go • Draft)"
	opens := true
	draftMeta := map[string]any{"demo_id": "abc123", "env": "demo", "source": "go-example"}

	draftReq := models.CreateDraftRequest{
		Body:    &body,
		Subject: &subj,
		To:      []models.EmailName{{Email: rcpt}},
		TrackingOptions: &models.TrackingOptions{
			Opens: &opens,
		},
		Metadata: draftMeta,
	}
	created, err := client.Drafts().Create(ctx, grant, draftReq)
	if err != nil {
		log.Fatalf("create draft: %v", err)
	}
	fmt.Printf("✓ Draft created (id=%s, request-id=%s)\n", created.Data.ID, created.Headers.Get("x-request-id"))
	pp(created.Data)

	// -----------------------------------------
	// 2) Send the draft -> returns Message object
	// -----------------------------------------
	sentFromDraft, err := client.Drafts().Send(ctx, grant, created.Data.ID)
	if err != nil {
		log.Fatalf("send draft: %v", err)
	}
	fmt.Printf("\n✓ Draft sent (message id=%s, request-id=%s)\n", sentFromDraft.Data.ID, sentFromDraft.Headers.Get("x-request-id"))
	fmt.Println("Message.metadata from draft:")
	pp(sentFromDraft.Data.Metadata)

	// ----------------------------------------------------------------
	// 3) Send a NEW message directly via Messages().Send with metadata
	//    (Demonstrates metadata at send-time without going through Drafts)
	// ----------------------------------------------------------------
	body2 := "Hello again — sent directly via Messages.Send"
	subj2 := "Metadata Demo (Go • Direct Send)"
	sendMeta := map[string]any{"campaign": "spring", "tier": "gold"}

	msgReq := models.SendMessageRequest{
		// Embedded CreateDraftRequest fields:
		CreateDraftRequest: models.CreateDraftRequest{
			Body:     &body2,
			Subject:  &subj2,
			To:       []models.EmailName{{Email: rcpt}},
			Metadata: sendMeta,
			TrackingOptions: &models.TrackingOptions{
				Opens: &opens,
			},
		},
		// If your API allows setting From here, do it; otherwise it uses the grant's identity.
		// From: []models.EmailName{{Email: "you@example.com"}},
		// UseDraft: ptrBool(false), // optional, defaults to false if omitted
	}
	sentDirect, err := client.Messages().Send(ctx, grant, msgReq)
	if err != nil {
		log.Fatalf("messages.send: %v", err)
	}
	fmt.Printf("\n✓ Direct send (message id=%s, request-id=%s)\n", sentDirect.Data.ID, sentDirect.Headers.Get("x-request-id"))
	fmt.Println("Message.metadata from direct send:")
	pp(sentDirect.Data.Metadata)

	// ----------------------------------------------------------
	// 4) Fetch the direct-sent message to confirm metadata fields
	// ----------------------------------------------------------
	found, err := client.Messages().Find(ctx, grant, sentDirect.Data.ID, &models.FindMessageQueryParams{
		// Fields: msgFieldsPtr(models.MessageFieldsStandard), // optional
		// Select: strPtr("id,subject,metadata"),               // optional to trim response
	})
	if err != nil {
		log.Fatalf("find message: %v", err)
	}
	fmt.Printf("\n✓ Fetched message (request-id=%s)\n", found.RequestID)
	fmt.Println("Fetched metadata:")
	pp(found.Data.Metadata)
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}

func pp(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }

// Optional helpers if you need them later:
// func msgFieldsPtr(f models.MessageFields) *models.MessageFields { return &f }
// func strPtr(s string) *string { return &s }
// func ptrBool(b bool) *bool { return &b }
