package models

// ----- Enums -----

type AccessType string

const (
	AccessTypeOnline  AccessType = "online"
	AccessTypeOffline AccessType = "offline"
)

type Provider string

const (
	ProviderGoogle          Provider = "google"
	ProviderIMAP            Provider = "imap"
	ProviderMicrosoft       Provider = "microsoft"
	ProviderICloud          Provider = "icloud"
	ProviderVirtualCalendar Provider = "virtual-calendar"
	ProviderYahoo           Provider = "yahoo"
	ProviderEWS             Provider = "ews"
	ProviderZoom            Provider = "zoom"
)

type Prompt string

const (
	PromptSelectProvider       Prompt = "select_provider"
	PromptDetect               Prompt = "detect"
	PromptSelectProviderDetect Prompt = "select_provider,detect"
	PromptDetectSelectProvider Prompt = "detect,select_provider"
)

// ----- Request/Config payloads -----

// URLForAuthenticationConfig mirrors the Python TypedDict with optional fields as pointers/slices.
type URLForAuthenticationConfig struct {
	ClientID           string      `json:"client_id"`
	RedirectURI        string      `json:"redirect_uri"`
	Provider           *Provider   `json:"provider,omitempty"`
	AccessType         *AccessType `json:"access_type,omitempty"`
	Prompt             *Prompt     `json:"prompt,omitempty"`
	Scope              []string    `json:"scope,omitempty"`
	IncludeGrantScopes *bool       `json:"include_grant_scopes,omitempty"`
	State              *string     `json:"state,omitempty"`
	LoginHint          *string     `json:"login_hint,omitempty"`
}

// URLForAdminConsentConfig extends URLForAuthenticationConfig with a required CredentialID.
type URLForAdminConsentConfig struct {
	URLForAuthenticationConfig
	CredentialID string `json:"credential_id"`
}

type CodeExchangeRequest struct {
	RedirectURI  string  `json:"redirect_uri"`
	Code         string  `json:"code"`
	ClientID     string  `json:"client_id"`
	ClientSecret *string `json:"client_secret,omitempty"`
	CodeVerifier *string `json:"code_verifier,omitempty"`
}

type TokenExchangeRequest struct {
	RedirectURI  string  `json:"redirect_uri"`
	RefreshToken string  `json:"refresh_token"`
	ClientID     string  `json:"client_id"`
	ClientSecret *string `json:"client_secret,omitempty"`
}

// ----- Responses -----

type CodeExchangeResponse struct {
	AccessToken  string    `json:"access_token"`
	GrantID      string    `json:"grant_id"`
	ExpiresIn    int       `json:"expires_in"`
	Email        *string   `json:"email,omitempty"`
	RefreshToken *string   `json:"refresh_token,omitempty"`
	Scope        *string   `json:"scope,omitempty"`
	IDToken      *string   `json:"id_token,omitempty"`
	TokenType    *string   `json:"token_type,omitempty"`
	Provider     *Provider `json:"provider,omitempty"`
}

type TokenInfoResponse struct {
	Iss   string  `json:"iss"`
	Aud   string  `json:"aud"`
	Iat   int64   `json:"iat"`
	Exp   int64   `json:"exp"`
	Sub   *string `json:"sub,omitempty"`
	Email *string `json:"email,omitempty"`
}

type PkceAuthURL struct {
	Secret     string `json:"secret"`
	SecretHash string `json:"secret_hash"`
	URL        string `json:"url"`
}

// ----- Provider detection -----

type ProviderDetectParams struct {
	Email            string `json:"email"`
	AllProviderTypes *bool  `json:"all_provider_types,omitempty"`
}

type ProviderDetectResponse struct {
	EmailAddress string  `json:"email_address"`
	Detected     bool    `json:"detected"`
	Provider     *string `json:"provider,omitempty"`
	Type         *string `json:"type,omitempty"`
}
