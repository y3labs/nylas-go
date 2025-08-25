package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

func TestRedirectURI_Deserialization(t *testing.T) {
	js := `{
		"id": "0556d035-6cb6-4262-a035-6b77e11cf8fc",
		"url": "http://localhost/abc",
		"platform": "web",
		"settings": {
			"origin": "string",
			"bundle_id": "string",
			"app_store_id": "string",
			"team_id": "string",
			"package_name": "string",
			"sha1_certificate_fingerprint": "string"
		}
	}`
	var ru models.RedirectURI
	if err := json.Unmarshal([]byte(js), &ru); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if ru.ID != "0556d035-6cb6-4262-a035-6b77e11cf8fc" {
		t.Fatalf("id mismatch: %q", ru.ID)
	}
	if ru.URL != "http://localhost/abc" {
		t.Fatalf("url mismatch: %q", ru.URL)
	}
	if ru.Platform != "web" {
		t.Fatalf("platform mismatch: %q", ru.Platform)
	}
	if ru.Settings.Origin != "string" ||
		ru.Settings.BundleID != "string" ||
		ru.Settings.AppStoreID != "string" ||
		ru.Settings.TeamID != "string" ||
		ru.Settings.PackageName != "string" ||
		ru.Settings.SHA1CertificateFingerprint != "string" {
		t.Fatalf("settings mismatch: %#v", ru.Settings)
	}
}

func TestRedirectURIs_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/applications/redirect-uris")
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": []any{
				map[string]any{
					"id":       "rid-1",
					"url":      "http://localhost/abc",
					"platform": "web",
					"settings": map[string]any{
						"origin":                       "string",
						"bundle_id":                    "string",
						"app_store_id":                 "string",
						"team_id":                      "string",
						"package_name":                 "string",
						"sha1_certificate_fingerprint": "string",
					},
				},
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.RedirectURIs().List(context.Background())
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(res.Data) != 1 || res.Data[0].ID != "rid-1" {
		t.Fatalf("unexpected list data: %#v", res.Data)
	}
}

func TestRedirectURIs_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/applications/redirect-uris/redirect_uri-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":       "redirect_uri-123",
				"url":      "http://localhost/abc",
				"platform": "web",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.RedirectURIs().Get(context.Background(), "redirect_uri-123")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if res.Data.ID != "redirect_uri-123" {
		t.Fatalf("unexpected id: %#v", res.Data)
	}
}

func TestRedirectURIs_Create(t *testing.T) {
	want := models.CreateRedirectURIRequest{
		URL:      "http://localhost/abc",
		Platform: "web",
		Settings: models.WritableRedirectURISettings{
			Origin:                     "string",
			BundleID:                   "string",
			AppStoreID:                 "string",
			TeamID:                     "string",
			PackageName:                "string",
			SHA1CertificateFingerprint: "string",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/applications/redirect-uris")
		var got models.CreateRedirectURIRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.URL != want.URL || got.Platform != want.Platform ||
			got.Settings.Origin != want.Settings.Origin ||
			got.Settings.BundleID != want.Settings.BundleID ||
			got.Settings.AppStoreID != want.Settings.AppStoreID ||
			got.Settings.TeamID != want.Settings.TeamID ||
			got.Settings.PackageName != want.Settings.PackageName ||
			got.Settings.SHA1CertificateFingerprint != want.Settings.SHA1CertificateFingerprint {
			t.Fatalf("create body mismatch:\n got: %#v\nwant: %#v", got, want)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":       "new-redirect-uri",
				"url":      got.URL,
				"platform": got.Platform,
				"settings": got.Settings,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.RedirectURIs().Create(context.Background(), want)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if res.Data.ID != "new-redirect-uri" || res.Data.URL != want.URL || res.Data.Platform != want.Platform {
		t.Fatalf("unexpected create resp: %#v", res.Data)
	}
}

func TestRedirectURIs_Update(t *testing.T) {
	want := models.UpdateRedirectURIRequest{
		URL:      "http://localhost/abc",
		Platform: "web",
		Settings: models.WritableRedirectURISettings{
			Origin:                     "string",
			BundleID:                   "string",
			AppStoreID:                 "string",
			TeamID:                     "string",
			PackageName:                "string",
			SHA1CertificateFingerprint: "string",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPut, "/v3/applications/redirect-uris/redirect_uri-123")
		var got models.UpdateRedirectURIRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.URL != want.URL || got.Platform != want.Platform ||
			got.Settings.Origin != want.Settings.Origin ||
			got.Settings.BundleID != want.Settings.BundleID ||
			got.Settings.AppStoreID != want.Settings.AppStoreID ||
			got.Settings.TeamID != want.Settings.TeamID ||
			got.Settings.PackageName != want.Settings.PackageName ||
			got.Settings.SHA1CertificateFingerprint != want.Settings.SHA1CertificateFingerprint {
			t.Fatalf("update body mismatch:\n got: %#v\nwant: %#v", got, want)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":       "redirect_uri-123",
				"url":      got.URL,
				"platform": got.Platform,
				"settings": got.Settings,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.RedirectURIs().Update(context.Background(), "redirect_uri-123", want)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if res.Data.ID != "redirect_uri-123" || res.Data.URL != want.URL || res.Data.Platform != want.Platform {
		t.Fatalf("unexpected update resp: %#v", res.Data)
	}
}

func TestRedirectURIs_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodDelete, "/v3/applications/redirect-uris/redirect_uri-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-del-1",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.RedirectURIs().Delete(context.Background(), "redirect_uri-123")
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if res == nil || res.RequestID == "" {
		t.Fatalf("expected request_id in delete response, got: %#v", res)
	}
}
