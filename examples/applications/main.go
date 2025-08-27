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
	// Required env:
	//   NYLAS_API_KEY
	apiKey := must("NYLAS_API_KEY")

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// 1) Fetch application info
	app, err := client.Applications().Info(ctx)
	die(err)

	fmt.Printf("✓ Applications.Info (request-id=%s)\n", app.Headers.Get("x-request-id"))
	fmt.Printf("ApplicationID: %s\n", app.Data.ApplicationID)
	fmt.Printf("OrganizationID: %s\n", app.Data.OrganizationID)
	fmt.Printf("Region: %s\n", app.Data.Region)
	fmt.Printf("Environment: %s\n", app.Data.Environment)
	fmt.Printf("Branding.Name: %s\n", app.Data.Branding.Name)

	// If you want the full structure:
	fmt.Println("\n=== Full Application ===")
	pp(app.Data)

	uris, err := client.Applications().RedirectURIs().List(ctx)
	die(err)
	fmt.Printf("\nRedirect URIs (count=%d, request-id=%s)\n", len(uris.Data), uris.Headers.Get("x-request-id"))
	pp(uris.Data)
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
