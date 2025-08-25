package nylas

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/y3labs/nylas-go/nylas/models"
)

// AuthResource implements the Auth endpoints consistent with the Python SDK.
type AuthResource struct{ c *Client }

// -----------------------------
// Helpers (mirror Python helpers)
// -----------------------------

// hashPKCESecret replicates nylas.resources.auth._hash_pkce_secret:
// sha256(secret).hexdigest() -> base64 std encode -> trim '='
func hashPKCESecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	hex := fmt.Sprintf("%x", sum[:])                      // lowercase hex digest
	b64 := base64.StdEncoding.EncodeToString([]byte(hex)) // standard base64 (not URL-safe)
	return strings.TrimRight(b64, "=")
}

// buildAuthQuery: response_type=code, default access_type=online, join scopes with spaces.
// Adds optional fields if set (provider, prompt, include_grant_scopes, state, login_hint).
func buildAuthQuery(cfg models.URLForAuthenticationConfig) url.Values {
	q := url.Values{}
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURI)
	q.Set("response_type", "code")

	accessType := models.AccessTypeOnline
	if cfg.AccessType != nil {
		accessType = *cfg.AccessType
	}
	q.Set("access_type", string(accessType))

	if len(cfg.Scope) > 0 {
		// Python joins by space and expects %20 in URL
		q.Set("scope", strings.Join(cfg.Scope, " "))
	}
	if cfg.Provider != nil {
		q.Set("provider", string(*cfg.Provider))
	}
	if cfg.Prompt != nil {
		q.Set("prompt", string(*cfg.Prompt))
	}
	if cfg.IncludeGrantScopes != nil {
		if *cfg.IncludeGrantScopes {
			q.Set("include_grant_scopes", "true")
		} else {
			q.Set("include_grant_scopes", "false")
		}
	}
	if cfg.State != nil {
		q.Set("state", *cfg.State)
	}
	if cfg.LoginHint != nil {
		q.Set("login_hint", *cfg.LoginHint)
	}
	return q
}

// encodeQueryLikePython makes spaces appear as %20 (not '+') to match Python tests.
func encodeQueryLikePython(q url.Values) string {
	enc := q.Encode()
	// Go's url.Values.Encode uses application/x-www-form-urlencoded where spaces are '+'.
	// Python tests expect %20, so convert here.
	return strings.ReplaceAll(enc, "+", "%20")
}

func (a *AuthResource) urlAuthBuilder(q url.Values) string {
	return a.c.serverURL + "/v3/connect/auth?" + encodeQueryLikePython(q)
}

// -----------------------------
// URL builders
// -----------------------------

// URLForOAuth2 mirrors python Auth.url_for_oauth2
func (a *AuthResource) URLForOAuth2(cfg models.URLForAuthenticationConfig) string {
	q := buildAuthQuery(cfg)
	return a.urlAuthBuilder(q)
}

// URLForOAuth2PKCE mirrors python Auth.url_for_oauth2_pkce
// Returns the secret (verifier), its hashed representation, and the URL.
func (a *AuthResource) URLForOAuth2PKCE(cfg models.URLForAuthenticationConfig) (models.PkceAuthURL, error) {
	// Generate a random secret; 32 bytes -> base64 -> trim '=' to keep it token-ish
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return models.PkceAuthURL{}, err
	}
	secret := base64.StdEncoding.EncodeToString(secretBytes)
	secret = strings.TrimRight(secret, "=")

	secretHash := hashPKCESecret(secret)

	q := buildAuthQuery(cfg)
	q.Set("code_challenge", secretHash)
	q.Set("code_challenge_method", "s256") // lower-case to match Python

	url := a.urlAuthBuilder(q)
	return models.PkceAuthURL{
		Secret:     secret,
		SecretHash: secretHash,
		URL:        url,
	}, nil
}

// URLForAdminConsent mirrors python Auth.url_for_admin_consent (forces provider=microsoft, response_type=adminconsent)
func (a *AuthResource) URLForAdminConsent(cfg models.URLForAdminConsentConfig) string {
	q := buildAuthQuery(cfg.URLForAuthenticationConfig)
	// Override to admin consent
	q.Set("response_type", "adminconsent")
	// Force provider=microsoft (Python sets this before building)
	q.Set("provider", string(models.ProviderMicrosoft))
	// Credential ID is required
	q.Set("credential_id", cfg.CredentialID)
	return a.urlAuthBuilder(q)
}

// -----------------------------
// Token & info flows
// -----------------------------

// ExchangeCodeForToken mirrors python exchange_code_for_token
func (a *AuthResource) ExchangeCodeForToken(ctx context.Context, req models.CodeExchangeRequest) (*models.CodeExchangeResponse, error) {
	// Default client_secret to API key if not provided
	if req.ClientSecret == nil && a.c.apiKey != "" {
		cs := a.c.apiKey
		req.ClientSecret = &cs
	}

	// Add grant_type=authorization_code
	body := map[string]any{
		"redirect_uri": req.RedirectURI,
		"code":         req.Code,
		"client_id":    req.ClientID,
		"grant_type":   "authorization_code",
	}
	if req.ClientSecret != nil && *req.ClientSecret != "" {
		body["client_secret"] = *req.ClientSecret
	}
	if req.CodeVerifier != nil && *req.CodeVerifier != "" {
		body["code_verifier"] = *req.CodeVerifier
	}

	out, _, err := DoJSON[models.CodeExchangeResponse](a.c, ctx, http.MethodPost, "/v3/connect/token", nil, body, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshAccessToken mirrors python refresh_access_token
func (a *AuthResource) RefreshAccessToken(ctx context.Context, req models.TokenExchangeRequest) (*models.CodeExchangeResponse, error) {
	if req.ClientSecret == nil && a.c.apiKey != "" {
		cs := a.c.apiKey
		req.ClientSecret = &cs
	}

	body := map[string]any{
		"redirect_uri":  req.RedirectURI,
		"refresh_token": req.RefreshToken,
		"client_id":     req.ClientID,
		"grant_type":    "refresh_token",
	}
	if req.ClientSecret != nil && *req.ClientSecret != "" {
		body["client_secret"] = *req.ClientSecret
	}

	out, _, err := DoJSON[models.CodeExchangeResponse](a.c, ctx, http.MethodPost, "/v3/connect/token", nil, body, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IDTokenInfo mirrors python id_token_info (GET with ?id_token=...; returns Response envelope)
func (a *AuthResource) IDTokenInfo(ctx context.Context, idToken string) (*Response[models.TokenInfoResponse], error) {
	q := url.Values{}
	q.Set("id_token", idToken)

	out, headers, err := DoJSON[Response[models.TokenInfoResponse]](a.c, ctx, http.MethodGet, "/v3/connect/tokeninfo", q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// ValidateAccessToken mirrors python validate_access_token (GET with ?access_token=...)
func (a *AuthResource) ValidateAccessToken(ctx context.Context, accessToken string) (*Response[models.TokenInfoResponse], error) {
	q := url.Values{}
	q.Set("access_token", accessToken)

	out, headers, err := DoJSON[Response[models.TokenInfoResponse]](a.c, ctx, http.MethodGet, "/v3/connect/tokeninfo", q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Revoke mirrors python revoke (POST with ?token=...)
func (a *AuthResource) Revoke(ctx context.Context, token string) (bool, error) {
	q := url.Values{}
	q.Set("token", token)
	// No JSON body; expect 2xx
	_, _, err := DoJSON[map[string]any](a.c, ctx, http.MethodPost, "/v3/connect/revoke", q, nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CustomAuthentication mirrors python custom_authentication (POST /v3/connect/custom, returns Response[Grant])
func (a *AuthResource) CustomAuthentication(ctx context.Context, body map[string]any) (*Response[models.Grant], error) {
	out, headers, err := DoJSON[Response[models.Grant]](a.c, ctx, http.MethodPost, "/v3/connect/custom", nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// DetectProvider mirrors python detect_provider (POST /v3/providers/detect with query params; returns Response envelope)
func (a *AuthResource) DetectProvider(ctx context.Context, params models.ProviderDetectParams) (*Response[models.ProviderDetectResponse], error) {
	// Send as query params just like Python
	q := url.Values{}
	q.Set("email", params.Email)
	if params.AllProviderTypes != nil {
		if *params.AllProviderTypes {
			q.Set("all_provider_types", "true")
		} else {
			q.Set("all_provider_types", "false")
		}
	}

	out, headers, err := DoJSON[Response[models.ProviderDetectResponse]](a.c, ctx, http.MethodPost, "/v3/providers/detect", q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
