package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/y3labs/nylas-go/nylas"
)

func main() {
	apiKey := must("NYLAS_API_KEY")
	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// ------------------------------------------------------
	// 1) List grants (limit=5 for demo)
	// ------------------------------------------------------
	limit := 5
	grants, err := client.Grants().List(ctx, &nylas.ListGrantsParams{Limit: &limit})
	die(err)

	fmt.Printf("✓ Listed %d grants (request-id=%s)\n",
		len(grants.Data), grants.Headers.Get("x-request-id"))
	pp(grants.Data)

	if len(grants.Data) == 0 {
		fmt.Println("No grants available.")
		return
	}

	// ------------------------------------------------------
	// 2) Get details of the first grant
	// ------------------------------------------------------
	firstID := grants.Data[0].ID
	grant, err := client.Grants().Get(ctx, firstID)
	die(err)

	fmt.Printf("\n✓ Fetched grant %s (request-id=%s)\n",
		grant.Data.ID, grant.Headers.Get("x-request-id"))
	pp(grant.Data)
}

/* ---------- helpers ---------- */

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
func pp(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
