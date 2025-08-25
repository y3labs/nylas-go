package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/y3labs/nylas-go/nylas"
	"github.com/y3labs/nylas-go/nylas/models"
)

func main() {
	apiKey := must("NYLAS_API_KEY")
	grant := must("NYLAS_GRANT_ID")

	client := nylas.NewClient(apiKey)
	ctx := context.Background()

	// Limit via env: MAX=5 go run main.go
	limit := 5
	if v := os.Getenv("MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	// 1) List a few messages with STANDARD fields
	lst, err := client.Messages().List(ctx, grant, &models.ListMessagesQueryParams{
		Limit:  intPtr(limit),
		Fields: msgFieldsPtr(models.MessageFieldsStandard),
		Select: strPtr("id,subject,from,date,snippet"), // trims payload, optional
	})
	if err != nil {
		log.Fatalf("list messages: %v", err)
	}
	if len(lst.Data) == 0 {
		fmt.Println("No messages.")
		return
	}

	fmt.Printf("=== Standard messages (request-id: %s) ===\n", lst.Headers.Get("x-request-id"))
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "IDX\tID\tSUBJECT\tFROM\tDATE")
	for i, m := range lst.Data {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\n",
			i,
			m.ID,
			s(m.Subject),
			formatAddrs(m.From),
			i64(m.Date),
		)
	}
	w.Flush()

	// Choose first message id to explore field variants
	msgID := lst.Data[0].ID
	fmt.Printf("\nUsing message id: %s\n\n", msgID)

	// 2) include_headers
	withHdrs, err := client.Messages().Find(ctx, grant, msgID, &models.FindMessageQueryParams{
		Fields: msgFieldsPtr(models.MessageFieldsIncludeHeaders),
	})
	if err != nil {
		log.Fatalf("find (include_headers): %v", err)
	}
	fmt.Printf("=== include_headers (request-id: %s) ===\n", withHdrs.Headers.Get("x-request-id"))
	printHeaders(withHdrs.Data.Headers)

	// 3) include_tracking_options
	withTracking, err := client.Messages().Find(ctx, grant, msgID, &models.FindMessageQueryParams{
		Fields: msgFieldsPtr(models.MessageFieldsIncludeTracking),
	})
	if err != nil {
		log.Fatalf("find (include_tracking_options): %v", err)
	}
	fmt.Printf("\n=== include_tracking_options (request-id: %s) ===\n", withTracking.Headers.Get("x-request-id"))
	printTracking(withTracking.Data.Tracking)

	// 4) raw_mime
	withRaw, err := client.Messages().Find(ctx, grant, msgID, &models.FindMessageQueryParams{
		Fields: msgFieldsPtr(models.MessageFieldsRawMIME),
	})
	if err != nil {
		log.Fatalf("find (raw_mime): %v", err)
	}
	fmt.Printf("\n=== raw_mime (request-id: %s) ===\n", withRaw.Headers.Get("x-request-id"))
	if withRaw.Data.RawMIME != nil {
		fmt.Printf("raw_mime: %d bytes (base64url)\n", len(*withRaw.Data.RawMIME))
	} else {
		fmt.Println("raw_mime: <nil>")
	}

	// Bonus dump (uncomment for structure inspection)
	dump("STANDARD exemplar", lst.Data[0])
	dump("include_headers exemplar", withHdrs.Data)
	dump("include_tracking exemplar", withTracking.Data)
	dump("raw_mime exemplar", withRaw.Data)
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}

func intPtr(n int) *int                                         { return &n }
func msgFieldsPtr(f models.MessageFields) *models.MessageFields { return &f }
func strPtr(s string) *string                                   { return &s }

func s(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func i64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func formatAddrs(addrs []models.EmailName) string {
	if len(addrs) == 0 {
		return ""
	}
	out := ""
	for i, a := range addrs {
		if a.Name != nil && *a.Name != "" {
			out += fmt.Sprintf("%s <%s>", *a.Name, a.Email)
		} else {
			out += a.Email
		}
		if i != len(addrs)-1 {
			out += ", "
		}
	}
	return out
}

func printHeaders(hdrs []models.MessageHeader) {
	if len(hdrs) == 0 {
		fmt.Println("(no headers)")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVALUE")
	for _, h := range hdrs {
		fmt.Fprintf(w, "%s\t%s\n", h.Name, h.Value)
	}
	w.Flush()
}

func printTracking(t *models.TrackingOptions) {
	if t == nil {
		fmt.Println("(no tracking options)")
		return
	}
	b, _ := json.MarshalIndent(t, "", "  ")
	fmt.Println(string(b))
}

func dump(label string, v any) {
	fmt.Println("----", label, "----")
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
