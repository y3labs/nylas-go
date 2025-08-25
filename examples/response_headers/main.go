package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/y3labs/nylas-go/nylas"
)

func main() {
	apiKey := must("NYLAS_API_KEY")
	grant := must("NYLAS_GRANT_ID")
	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// (1) Headers from list response
	lst, err := client.Messages().List(ctx, grant, nil)
	die(err)
	fmt.Println("List -> request-id:", reqID(lst.RequestID, lst.Headers),
		"rate-remaining:", rateRemaining(lst.Headers))
	// dump all headers (uncomment if troubleshooting)
	dumpHeaders(lst.Headers)

	// (2) Headers from single-item response
	if len(lst.Data) > 0 {
		id := lst.Data[0].ID
		one, err := client.Messages().Find(ctx, grant, id, nil)
		die(err)
		fmt.Println("Find -> request-id:", reqID(one.RequestID, one.Headers),
			"rate-remaining:", rateRemaining(one.Headers))
		dumpHeaders(one.Headers)
	} else {
		fmt.Println("Find -> skipped (no messages)")
	}

	// (3) Headers from error response (force an error)
	_, err = client.Events().List(ctx, grant, nylas.ListEventsParams{
		CalendarID: "does-not-exist",
	})
	if err != nil {
		if e, ok := nylas.IsAPIError(err); ok {
			fmt.Println("Error -> request-id:", e.RequestID, "status:", e.StatusCode, "type:", e.Type)
		} else {
			fmt.Println("Error ->", err)
		}
	} else {
		log.Fatal("expected an error, got nil")
	}
}

/* helpers */

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// prefer body field, fall back to header
func reqID(bodyID string, h map[string][]string) string {
	if bodyID != "" {
		return bodyID
	}
	// http.Header is case-insensitive; use both common spellings defensively
	if v := getHeader(h, "x-request-id"); v != "" {
		return v
	}
	return ""
}

func rateRemaining(h map[string][]string) string {
	// Providers sometimes use one of these
	if v := getHeader(h, "x-ratelimit-remaining"); v != "" {
		return v
	}
	if v := getHeader(h, "x-rate-limit-remaining"); v != "" {
		return v
	}
	return ""
}

func getHeader(h map[string][]string, key string) string {
	// mimic http.Header.Get semantics
	for k, vals := range h {
		if len(vals) == 0 {
			continue
		}
		if equalFold(k, key) {
			return vals[0]
		}
	}
	return ""
}

// simple ASCII case-fold (good enough for header keys)
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

func dumpHeaders(h http.Header) {
	for k, v := range h {
		fmt.Printf("%s: %v\n", k, v)
	}
}
