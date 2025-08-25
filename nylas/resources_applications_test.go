package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestApplicationsInfo(t *testing.T) {
	// Arrange: mock API server
	type payload struct {
		RequestID string `json:"request_id"`
		Data      any    `json:"data"`
	}
	reply := payload{
		RequestID: "req-123",
		Data: map[string]any{
			"application_id":  "ad410018-d306-43f9-8361-fa5d7b2172e0",
			"organization_id": "f5db4482-dbbe-4b32-b347-61c260d803ce",
			"region":          "us",
			"environment":     "production",
			"branding": map[string]any{
				"name":        "My application",
				"icon_url":    "https://my-app.com/my-icon.png",
				"website_url": "https://my-app.com",
				"description": "Online banking application.",
			},
			"hosted_authentication": map[string]any{
				"background_image_url": "https://my-app.com/bg.jpg",
				"alignment":            "left",
				"color_primary":        "#dc0000",
				"color_secondary":      "#000056",
				"title":                "string",
				"subtitle":             "string",
				"background_color":     "#003400",
				"spacing":              5,
			},
			"callback_uris": []any{
				map[string]any{
					"id":       "0556d035-6cb6-4262-a035-6b77e11cf8fc",
					"url":      "string",
					"platform": "web",
					"settings": map[string]any{
						"origin":                       "string",
						"bundle_id":                    "string",
						"package_name":                 "string",
						"sha1_certificate_fingerprint": "string",
					},
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert path & method
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/applications" {
			t.Fatalf("path = %s, want /v3/applications", r.URL.Path)
		}
		// (Optional) assert auth header is present
		if got := r.Header.Get("Authorization"); got == "" {
			t.Fatalf("missing Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	}))
	defer ts.Close()

	// Act: call SDK
	c := NewClient("test-key", WithServerURL(ts.URL))
	res, err := c.Applications().Info(context.Background())
	if err != nil {
		t.Fatalf("Applications.Info error: %v", err)
	}

	// Assert: envelope
	if res == nil {
		t.Fatal("nil response")
	}
	if got, want := res.RequestID, "req-123"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}

	// Assert: strongly typed data
	var d models.Application = res.Data
	if got, want := d.ApplicationID, "ad410018-d306-43f9-8361-fa5d7b2172e0"; got != want {
		t.Fatalf("ApplicationID = %q, want %q", got, want)
	}
	if got, want := d.OrganizationID, "f5db4482-dbbe-4b32-b347-61c260d803ce"; got != want {
		t.Fatalf("OrganizationID = %q, want %q", got, want)
	}
	if got, want := d.Region, models.RegionUS; got != want {
		t.Fatalf("Region = %q, want %q", got, want)
	}
	if got, want := d.Environment, models.ProdEnv; got != want {
		t.Fatalf("Environment = %q, want %q", got, want)
	}

	// Branding
	if got, want := d.Branding.Name, "My application"; got != want {
		t.Fatalf("Branding.Name = %q, want %q", got, want)
	}
	if got, want := d.Branding.IconURL, "https://my-app.com/my-icon.png"; got != want {
		t.Fatalf("Branding.IconURL = %q, want %q", got, want)
	}
	if got, want := d.Branding.WebsiteURL, "https://my-app.com"; got != want {
		t.Fatalf("Branding.WebsiteURL = %q, want %q", got, want)
	}
	if got, want := d.Branding.Description, "Online banking application."; got != want {
		t.Fatalf("Branding.Description = %q, want %q", got, want)
	}

	// Hosted Auth
	if got, want := d.HostedAuthentication.BackgroundImageURL, "https://my-app.com/bg.jpg"; got != want {
		t.Fatalf("HostedAuth.BackgroundImageURL = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.Alignment, "left"; got != want {
		t.Fatalf("HostedAuth.Alignment = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.ColorPrimary, "#dc0000"; got != want {
		t.Fatalf("HostedAuth.ColorPrimary = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.ColorSecondary, "#000056"; got != want {
		t.Fatalf("HostedAuth.ColorSecondary = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.Title, "string"; got != want {
		t.Fatalf("HostedAuth.Title = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.Subtitle, "string"; got != want {
		t.Fatalf("HostedAuth.Subtitle = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.BackgroundColor, "#003400"; got != want {
		t.Fatalf("HostedAuth.BackgroundColor = %q, want %q", got, want)
	}
	if got, want := d.HostedAuthentication.Spacing, 5; got != want {
		t.Fatalf("HostedAuth.Spacing = %d, want %d", got, want)
	}

	// Callback URIs
	if len(d.CallbackURIs) != 1 {
		t.Fatalf("CallbackURIs len = %d, want 1", len(d.CallbackURIs))
	}
	cb := d.CallbackURIs[0]
	if got, want := cb.ID, "0556d035-6cb6-4262-a035-6b77e11cf8fc"; got != want {
		t.Fatalf("CallbackURIs[0].ID = %q, want %q", got, want)
	}
	if got, want := cb.URL, "string"; got != want {
		t.Fatalf("CallbackURIs[0].URL = %q, want %q", got, want)
	}
	if got, want := cb.Platform, "web"; got != want {
		t.Fatalf("CallbackURIs[0].Platform = %q, want %q", got, want)
	}
	if got, want := cb.Settings.Origin, "string"; got != want {
		t.Fatalf("CallbackURIs[0].Settings.Origin = %q, want %q", got, want)
	}
	if got, want := cb.Settings.BundleID, "string"; got != want {
		t.Fatalf("CallbackURIs[0].Settings.BundleID = %q, want %q", got, want)
	}
	if got, want := cb.Settings.PackageName, "string"; got != want {
		t.Fatalf("CallbackURIs[0].Settings.PackageName = %q, want %q", got, want)
	}
	if got, want := cb.Settings.SHA1CertificateFingerprint, "string"; got != want {
		t.Fatalf("CallbackURIs[0].Settings.SHA1CertificateFingerprint = %q, want %q", got, want)
	}
}

func TestApplicationsRedirectURIsAccessor(t *testing.T) {
	// This Test isnt really necessary as RedirectURIs is exposed on the Client
	app := &ApplicationsResource{c: NewClient("k")}
	ru := app.RedirectURIs() // requires tiny forwarder on ApplicationsResource
	if ru == nil {
		t.Fatal("Applications.RedirectURIs() returned nil")
	}
}
