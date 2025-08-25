package nylas

import (
	"context"
	"net/http"
	"time"

	"github.com/y3labs/nylas-go/nylas/models"
)

// Region identifies an API region.

// Default server URLs per region.
const (
	serverUS = "https://api.us.nylas.com"
	serverEU = "https://api.eu.nylas.com"
)

// Option applies a functional option to Client.
type Option func(*Client)

// Client is the root SDK client.
type Client struct {
	apiKey    string
	serverURL string
	region    models.Region
	http      *http.Client
	userAgent string
}

// NewClient constructs a new Client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey: apiKey,
		region: models.RegionUS,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.serverURL == "" {
		// default by region
		switch c.region {
		case models.RegionEU:
			c.serverURL = serverEU
		default:
			c.serverURL = serverUS
		}
	}
	if c.http == nil {
		c.http = &http.Client{Timeout: 30 * time.Second}
	}
	if c.userAgent == "" {
		c.userAgent = "nylas-go/0.1.0"
	}
	return c
}

// WithRegion sets the region and corresponding base URL (unless already overridden).
func WithRegion(r models.Region) Option {
	return func(c *Client) {
		c.region = r
		if c.serverURL == "" {
			switch r {
			case models.RegionEU:
				c.serverURL = serverEU
			default:
				c.serverURL = serverUS
			}
		}
	}
}

// WithServerURL overrides the base server URL.
func WithServerURL(url string) Option {
	return func(c *Client) { c.serverURL = url }
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithUserAgent sets the User-Agent header value.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithTimeout sets a request timeout on the underlying http.Client if none provided.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.http == nil {
			c.http = &http.Client{Timeout: d}
			return
		}
		c.http.Timeout = d
	}
}

// internal: base execution

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	// Default headers (only when missing)
	if req.Header.Get("Authorization") == "" && c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if c.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	return c.http.Do(req)
}

// region accessors for resources
func (c *Client) Applications() *ApplicationsResource { return &ApplicationsResource{c} }

func (c *Client) Attachments() *AttachmentsResource { return &AttachmentsResource{c} }

func (c *Client) Auth() *AuthResource { return &AuthResource{c} }

func (c *Client) Availability() *AvailabilityResource { return &AvailabilityResource{c} }

func (c *Client) Calendars() *CalendarsResource { return &CalendarsResource{c} }

func (c *Client) Contacts() *ContactsResource { return &ContactsResource{c} }

func (c *Client) Connectors() *ConnectorsResource { return &ConnectorsResource{c} }

func (c *Client) Drafts() *DraftsResource { return &DraftsResource{c} }

func (c *Client) Events() *EventsResource { return &EventsResource{c} }

func (c *Client) Folders() *FoldersResource { return &FoldersResource{c} }

func (c *Client) Grants() *GrantsResource { return &GrantsResource{c} }

func (c *Client) Messages() *MessagesResource { return &MessagesResource{c} }

func (c *Client) Notetakers() *NotetakersResource { return &NotetakersResource{c} }

func (c *Client) RedirectURIs() *RedirectURIsResource { return &RedirectURIsResource{c} }

func (c *Client) Scheduler() *Scheduler { return &Scheduler{c} }

func (c *Client) Threads() *ThreadsResource { return &ThreadsResource{c} }

func (c *Client) Webhooks() *WebhooksResource { return &WebhooksResource{c} }
