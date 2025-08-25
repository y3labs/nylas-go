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

	max := 50
	if v := os.Getenv("MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			max = n
		}
	}

	total := 0
	var cursor *string

	for {
		resp, err := client.Messages().List(ctx, grant, &models.ListMessagesQueryParams{
			PageToken: cursor, // leave nil first page; pass through afterward
			Limit:     intPtr(15),
			// Fields: []string{"id","subject","from","date"}, // uncomment if your SDK supports field selection
		})
		if err != nil {
			log.Fatalf("list messages: %v", err)
		}
		if len(resp.Data) == 0 {
			if total == 0 {
				fmt.Println("No messages.")
			}
			break
		}

		// Pretty table
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		if total == 0 {
			fmt.Fprintln(w, "IDX\tID\tSUBJECT\tFROM")
		}
		for i, m := range resp.Data {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				total+i,
				m.ID,
				s(m.Subject),
				formatAddrs(m.From),
			)
		}
		w.Flush()

		total += len(resp.Data)
		if total >= max || resp.NextCursor == nil || *resp.NextCursor == "" {
			break
		}
		cursor = resp.NextCursor
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
func intPtr(n int) *int { return &n }

// s safely de-references *string
func s(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// formatAddrs renders []models.EmailName as "Name <email>, ..."
func formatAddrs(addrs []models.EmailName) string {
	if len(addrs) == 0 {
		return ""
	}
	out := make([]string, 0, len(addrs))
	for _, a := range addrs {
		name := s(a.Name)
		if name != "" {
			out = append(out, fmt.Sprintf("%s <%s>", name, a.Email))
		} else {
			out = append(out, a.Email)
		}
	}
	// join without importing strings for one-liner clarity
	b, _ := json.Marshal(out)      // ["A <a@x>","b@x"]
	return string(b[1 : len(b)-1]) // A <a@x>, "b@x"  (quick & dirty)
}
