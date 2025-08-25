package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/y3labs/nylas-go/nylas/models"
)

// --- small helpers ---

func mustReadJSONBody(t *testing.T, r *http.Request, v any) {
	t.Helper()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

func assertQueryHas(t *testing.T, r *http.Request, key, want string) {
	t.Helper()
	got := r.URL.Query().Get(key)
	if got != want {
		t.Fatalf("query[%s] = %q, want %q", key, got, want)
	}
}

// --- tests ---

func Test_hashPKCESecret(t *testing.T) {
	// Matches the Python test exact value for "nylas"
	want := "ZTk2YmY2Njg2YTNjMzUxMGU5ZTkyN2RiNzA2OWNiMWNiYTliOTliMDIyZjQ5NDgzYTZjZTMyNzA4MDllNjhhMg"
	got := hashPKCESecret("nylas")
	if got != want {
		t.Fatalf("hashPKCESecret = %q, want %q", got, want)
	}
}

func TestURLForOAuth2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")

	cfg := models.URLForAuthenticationConfig{
		ClientID:    "abc-123",
		RedirectURI: "https://example.com/oauth/callback",
		Scope:       []string{"email.read_only", "calendar", "contacts"},
		LoginHint:   strptr("test@gmail.com"),
		Provider:    ptr(models.ProviderGoogle),
		Prompt:      ptr(models.PromptSelectProviderDetect),
		State:       strptr("abc-123-state"),
	}
	u := c.Auth().URLForOAuth2(cfg)

	got, err := url.Parse(u)
	if err != nil {
		t.Fatalf("parse got URL: %v", err)
	}

	if got.Scheme != "http" || got.Host != strings.TrimPrefix(ts.URL, "http://") || got.Path != "/v3/connect/auth" {
		t.Fatalf("unexpected base/path: %s", u)
	}

	wantQ := url.Values{}
	wantQ.Set("client_id", "abc-123")
	wantQ.Set("redirect_uri", "https://example.com/oauth/callback")
	wantQ.Set("response_type", "code")
	wantQ.Set("access_type", "online")
	wantQ.Set("scope", "email.read_only calendar contacts")
	wantQ.Set("login_hint", "test@gmail.com")
	wantQ.Set("provider", "google")
	wantQ.Set("prompt", "select_provider,detect")
	wantQ.Set("state", "abc-123-state")

	gotQ := got.Query()

	if !reflect.DeepEqual(gotQ, wantQ) {
		t.Errorf("query mismatch:\n got:  %v\n want: %v", gotQ, wantQ)
	}
}

func TestURLForAdminConsent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	c := newTestClient(ts.URL, "test-key")

	cfg := models.URLForAdminConsentConfig{
		URLForAuthenticationConfig: models.URLForAuthenticationConfig{
			ClientID:    "abc-123",
			RedirectURI: "https://example.com/oauth/callback",
			Scope:       []string{"email.read_only", "calendar", "contacts"},
			LoginHint:   ptr("test@gmail.com"),
			Prompt:      ptr(models.PromptSelectProviderDetect),
			State:       ptr("abc-123-state"),
		},
		CredentialID: "cred-123",
	}

	u := c.Auth().URLForAdminConsent(cfg)

	got, err := url.Parse(u)
	if err != nil {
		t.Fatalf("parse got URL: %v", err)
	}
	// Base/path checks
	host := strings.TrimPrefix(ts.URL, "http://")
	if got.Scheme != "http" || got.Host != host || got.Path != "/v3/connect/auth" {
		t.Fatalf("unexpected base/path: %s", u)
	}

	// Build expected query (unencoded; url.Parse(got) will decode for us)
	wantQ := url.Values{}
	wantQ.Set("provider", "microsoft")
	wantQ.Set("credential_id", "cred-123")
	wantQ.Set("client_id", "abc-123")
	wantQ.Set("redirect_uri", "https://example.com/oauth/callback")
	wantQ.Set("scope", "email.read_only calendar contacts")
	wantQ.Set("login_hint", "test@gmail.com")
	wantQ.Set("prompt", "select_provider,detect")
	wantQ.Set("state", "abc-123-state")
	wantQ.Set("response_type", "adminconsent")
	wantQ.Set("access_type", "online")

	assertQueryEqual(t, wantQ, got.Query())
}

func TestExchangeCodeForToken(t *testing.T) {
	// Echo back a CodeExchange-like JSON (not wrapped in {"data":...})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connect/token")
		var got map[string]any
		mustReadJSONBody(t, r, &got)
		if got["client_id"] != "abc-123" || got["client_secret"] != "secret" ||
			got["code"] != "code" || got["redirect_uri"] != "https://example.com/oauth/callback" ||
			got["grant_type"] != "authorization_code" {
			t.Fatalf("unexpected request body: %#v", got)
		}
		w.Header().Set("x-request-id", "req-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "nylas_access_token",
			"expires_in":    3600,
			"id_token":      "jwt_token",
			"refresh_token": "nylas_refresh_token",
			"scope":         "https://www.googleapis.com/auth/gmail.readonly profile",
			"token_type":    "Bearer",
			"grant_id":      "grant_123",
			"provider":      "google",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Auth().ExchangeCodeForToken(context.Background(), models.CodeExchangeRequest{
		ClientID:     "abc-123",
		ClientSecret: strptr("secret"),
		Code:         "code",
		RedirectURI:  "https://example.com/oauth/callback",
	})
	if err != nil {
		t.Fatalf("ExchangeCodeForToken error: %v", err)
	}
}

func TestExchangeCodeForToken_NoSecret_UsesAPIKey(t *testing.T) {
	var body map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mustReadJSONBody(t, r, &body)
		// client_secret should be set from API key
		if body["client_secret"] != "nylas-api-key" {
			t.Fatalf("client_secret = %v, want nylas-api-key", body["client_secret"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "nylas-api-key")
	_, _ = c.Auth().ExchangeCodeForToken(context.Background(), models.CodeExchangeRequest{
		ClientID:    "abc-123",
		Code:        "code",
		RedirectURI: "https://example.com/oauth/callback",
	})
}

func TestIDTokenInfo(t *testing.T) {
	// Returns {"request_id": "...", "data": {...}}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connect/tokeninfo")
		assertQueryHas(t, r, "id_token", "id-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"iss":   "https://nylas.com",
				"aud":   "http://localhost:3030",
				"sub":   "Jaf84d88-£274-46cc-bbc9-aed7dac061c7",
				"email": "user@example.com",
				"iat":   1692094848,
				"exp":   1692095173,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.Auth().IDTokenInfo(context.Background(), "id-123")
	if err != nil {
		t.Fatalf("IDTokenInfo error: %v", err)
	}
	if res == nil || strval(res.Data.Email) != "user@example.com" {
		t.Fatalf("unexpected token info: %#v", res)
	}
}

func TestValidateAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodGet, "/v3/connect/tokeninfo")
		assertQueryHas(t, r, "access_token", "atk-123")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data":       map[string]any{"email": "user@example.com"},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Auth().ValidateAccessToken(context.Background(), "atk-123")
	if err != nil {
		t.Fatalf("ValidateAccessToken error: %v", err)
	}
}

func TestRefreshAccessToken(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connect/token")
		mustReadJSONBody(t, r, &got)
		if got["grant_type"] != "refresh_token" ||
			got["refresh_token"] != "refresh-12345" ||
			got["client_id"] != "abc-123" ||
			got["client_secret"] != "secret" ||
			got["redirect_uri"] != "https://example.com/oauth/callback" {
			t.Fatalf("unexpected refresh body: %#v", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	_, err := c.Auth().RefreshAccessToken(context.Background(), models.TokenExchangeRequest{
		RedirectURI:  "https://example.com/oauth/callback",
		RefreshToken: "refresh-12345",
		ClientID:     "abc-123",
		ClientSecret: strptr("secret"),
	})
	if err != nil {
		t.Fatalf("RefreshAccessToken error: %v", err)
	}
}

func TestRefreshAccessToken_NoSecret_UsesAPIKey(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mustReadJSONBody(t, r, &got)
		if got["client_secret"] != "nylas-api-key" {
			t.Fatalf("client_secret = %v, want nylas-api-key", got["client_secret"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "nylas-api-key")
	_, _ = c.Auth().RefreshAccessToken(context.Background(), models.TokenExchangeRequest{
		RedirectURI:  "https://example.com/oauth/callback",
		RefreshToken: "refresh-12345",
		ClientID:     "abc-123",
	})
}

func TestCustomAuthentication(t *testing.T) {
	// returns grant object wrapped in {"request_id","data":{...}}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connect/custom")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"id":           "e19f8e1a-eb1c-41c0-b6a6-d2e59daf7f47",
				"provider":     "google",
				"grant_status": "valid",
				"email":        "email@example.com",
				"scope":        []string{"Mail.Read", "User.Read", "offline_access"},
				"user_agent":   "string",
				"ip":           "string",
				"state":        "my-state",
				"created_at":   1617817109,
				"updated_at":   1617817109,
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	res, err := c.Auth().CustomAuthentication(context.Background(), map[string]any{
		"provider": "google", "settings": map[string]any{"foo": "bar"},
	})
	if err != nil {
		t.Fatalf("CustomAuthentication error: %v", err)
	}
	if res.Data.ID == "" || res.Data.Provider == "" {
		t.Fatalf("unexpected grant: %#v", res.Data)
	}
}

func TestRevoke(t *testing.T) {
	// Expect POST with query ?token=access_token (no json body)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/connect/revoke")
		if r.URL.Query().Get("token") != "access_token" {
			t.Fatalf("token qparam = %q, want access_token", r.URL.Query().Get("token"))
		}
		// 200 OK, empty or minimal body is fine
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")
	ok, err := c.Auth().Revoke(context.Background(), "access_token")
	if err != nil {
		t.Fatalf("Revoke error: %v", err)
	}
	if !ok {
		t.Fatalf("Revoke returned false, want true")
	}
}

func TestDetectProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/v3/providers/detect")
		q := r.URL.Query()
		if q.Get("email") != "test@gmail.com" || q.Get("all_provider_types") != "true" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "abc-123",
			"data": map[string]any{
				"email_address": "test@gmail.com",
				"detected":      true,
				"provider":      "google",
				"type":          "string",
			},
		})
	}))
	defer ts.Close()

	c := newTestClient(ts.URL, "test-key")

	params := models.ProviderDetectParams{
		Email:            "test@gmail.com",
		AllProviderTypes: boolptr(true), // helper returning *bool
	}
	res, err := c.Auth().DetectProvider(context.Background(), params)
	if err != nil {
		t.Fatalf("DetectProvider error: %v", err)
	}

	// res.Data is models.ProviderDetectResponse
	if !res.Data.Detected {
		t.Fatalf("expected detected=true, got false")
	}
	if res.Data.Provider == nil || strings.ToLower(*res.Data.Provider) != "google" {
		t.Fatalf("unexpected provider: %#v", res.Data.Provider)
	}
	if res.Data.EmailAddress != "test@gmail.com" {
		t.Fatalf("unexpected email_address: %q", res.Data.EmailAddress)
	}
}

func TestURLForOAuth2PKCE(t *testing.T) {
	c := NewClient("test-key", WithServerURL("http://localhost"))
	a := c.Auth() // use the real resource

	cfg := models.URLForAuthenticationConfig{
		ClientID:    "abc",
		RedirectURI: "https://app/cb",
		Scope:       []string{"email.read_only"},
	}

	out, err := a.URLForOAuth2PKCE(cfg)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out.Secret == "" || out.SecretHash == "" || out.URL == "" {
		t.Fatalf("missing fields: %#v", out)
	}

	// hash must match helper
	if want := hashPKCESecret(out.Secret); want != out.SecretHash {
		t.Fatalf("hash mismatch: got %q want %q", out.SecretHash, want)
	}

	// URL must contain the PKCE params
	u, err := url.Parse(out.URL)
	if err != nil {
		t.Fatalf("bad URL: %v", err)
	}
	q := u.Query()
	if strings.ToLower(q.Get("code_challenge_method")) != "s256" {
		t.Fatalf("missing/incorrect code_challenge_method: %q", q.Get("code_challenge_method"))
	}
	if q.Get("code_challenge") != out.SecretHash {
		t.Fatalf("missing/incorrect code_challenge: %q", q.Get("code_challenge"))
	}

	// sanity check a couple of base auth params
	if q.Get("client_id") != "abc" || q.Get("redirect_uri") != "https://app/cb" {
		t.Fatalf("base auth params missing: client_id=%q redirect_uri=%q", q.Get("client_id"), q.Get("redirect_uri"))
	}
}

func assertQueryEqual(t *testing.T, want, got url.Values) {
	t.Helper()

	// Check for missing or mismatched values
	for k, wantVals := range want {
		gotVals, ok := got[k]
		if !ok {
			t.Errorf("missing param %q", k)
			continue
		}
		if !reflect.DeepEqual(gotVals, wantVals) {
			t.Errorf("param %q: got %v, want %v", k, gotVals, wantVals)
		}
	}
	// Check for unexpected extras
	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf("unexpected extra param %q=%v", k, got[k])
		}
	}
}
