package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/y3labs/nylas-go/nylas"
)

func main() {
	apiKey := must("NYLAS_API_KEY")
	grant := must("NYLAS_GRANT_ID")
	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	badCal := "non-existent-calendar-123"
	fmt.Printf("Attempting to fetch events from calendar: %s\n", badCal)

	// Use the typed Events resource instead of raw DoJSON
	_, err := client.Events().List(ctx, grant, nylas.ListEventsParams{
		CalendarID: badCal, // required
		// (optionally include Start/End/Limit, not required to trigger error)
	})
	if err == nil {
		log.Fatal("expected an error, got nil")
	}

	// Surface a provider-originated error in a friendly way
	if e, ok := nylas.IsAPIError(err); ok {
		fmt.Println("Caught Nylas APIError:")
		fmt.Println("  Message       :", e.Message)
		fmt.Println("  Type          :", e.Type)                  // e.g., "invalid_request_error"
		fmt.Println("  Status Code   :", e.StatusCode)            // e.g., 404
		fmt.Println("  Request ID    :", e.RequestID)             // from header/body if present
		fmt.Println("  Provider Error:", e.ProviderErrorString()) // provider payload (if any)
		// If you want raw provider map:
		fmt.Printf("  Provider Raw   : %#v\n", e.ProviderError)
		return
	}

	if oe, ok := nylas.IsOAuthError(err); ok {
		fmt.Println("Caught Nylas OAuthError:",
			oe.ErrorDescription, "(code:", oe.ErrorCode, ")")
		return
	}

	log.Fatalf("unexpected error: %v", err)
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
