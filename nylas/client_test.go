package nylas

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestClientInit(t *testing.T) {
	c := NewClient(
		"test-key",
		WithServerURL("https://test.nylas.com"),
		WithTimeout(60*time.Second),
	)

	if got, want := c.apiKey, "test-key"; got != want {
		t.Fatalf("apiKey = %q, want %q", got, want)
	}
	if got, want := c.serverURL, "https://test.nylas.com"; got != want {
		t.Fatalf("serverURL = %q, want %q", got, want)
	}
	if got, want := c.http.Timeout, 60*time.Second; got != want {
		t.Fatalf("http.Timeout = %s, want %s", got, want)
	}
}

func TestClientInitDefaults(t *testing.T) {
	c := NewClient("test-key")

	if got, want := c.apiKey, "test-key"; got != want {
		t.Fatalf("apiKey = %q, want %q", got, want)
	}
	// Default region is US -> serverUS
	if got, want := c.serverURL, serverUS; got != want {
		t.Fatalf("serverURL = %q, want %q", got, want)
	}
	// Default timeout per NewClient (your SDK uses 30s by default)
	if got, want := c.http.Timeout, 30*time.Second; got != want {
		t.Fatalf("http.Timeout = %s, want %s", got, want)
	}
}

func TestClientInitDefaultWithEURegion(t *testing.T) {
	c := NewClient("test-key", WithRegion(models.RegionEU))

	// Default region is US -> serverUS
	// But We Explicitly Provide RegionEU and no server url
	// Should default to EU api URL
	if got, want := c.serverURL, serverEU; got != want {
		t.Fatalf("serverURL = %q, want %q", got, want)
	}
}

func TestClientInitRegionDefault(t *testing.T) {
	c := NewClient("test-key", WithRegion(models.RegionUS))

	// Default region is US -> serverUS
	if got, want := c.serverURL, serverUS; got != want {
		t.Fatalf("serverURL = %q, want %q", got, want)
	}
}

func TestClientInitNonDefaultHTTPTimeout(t *testing.T) {
	hc := &http.Client{Timeout: 30 * time.Second}
	c := NewClient("test-key", WithHTTPClient(hc), WithTimeout(60*time.Second))

	if got, want := c.http.Timeout, 60*time.Second; got != want {
		t.Fatalf("timeout = %s, want %s", got, want)
	}
}

func TestClientRegionSwitchSetsDefaultServer(t *testing.T) {
	c := NewClient("k", WithRegion(models.RegionEU))
	if got, want := c.serverURL, serverEU; got != want {
		t.Fatalf("serverURL = %q, want %q", got, want)
	}
}

func TestClientResourceAccessors(t *testing.T) {
	c := NewClient("k")

	type testCase struct {
		got any
		typ any
	}

	cases := []testCase{
		{c.Applications(), &ApplicationsResource{}},
		{c.Attachments(), &AttachmentsResource{}},
		{c.Auth(), &AuthResource{}},
		{c.Availability(), &AvailabilityResource{}},
		{c.Calendars(), &CalendarsResource{}},
		{c.Contacts(), &ContactsResource{}},
		{c.Connectors(), &ConnectorsResource{}},
		{c.Drafts(), &DraftsResource{}},
		{c.Events(), &EventsResource{}},
		{c.Folders(), &FoldersResource{}},
		{c.Grants(), &GrantsResource{}},
		{c.Messages(), &MessagesResource{}},
		{c.Notetakers(), &NotetakersResource{}},
		{c.RedirectURIs(), &RedirectURIsResource{}},
		{c.Threads(), &ThreadsResource{}},
		{c.Webhooks(), &WebhooksResource{}},
	}

	for _, tc := range cases {
		if tc.got == nil {
			t.Fatalf("resource accessor returned nil for %T", tc.typ)
		}
		gotT := reflect.TypeOf(tc.got)
		wantT := reflect.TypeOf(tc.typ)
		if gotT != wantT {
			t.Fatalf("accessor type = %v, want %v", gotT, wantT)
		}
	}
}

func TestClientSchedulerAccessor(t *testing.T) {
	c := NewClient("k")
	s := c.Scheduler()
	if s == nil {
		t.Fatal("Scheduler() returned nil")
	}
	// Optional: sanity-check the sub-accessors compile
	if s.Configurations() == nil || s.Bookings() == nil || s.Sessions() == nil {
		t.Fatal("scheduler sub-accessors returned nil")
	}
}

func TestWithUserAgent(t *testing.T) {
	c := NewClient("k", WithUserAgent("custom-UA/1.2"))
	if got, want := c.userAgent, "custom-UA/1.2"; got != want {
		t.Fatalf("userAgent = %q, want %q", got, want)
	}
}
