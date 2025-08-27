package models

// ----- Enums -----

type CredentialType string

const (
	CredentialTypeAdminConsent   CredentialType = "adminconsent"
	CredentialTypeServiceAccount CredentialType = "serviceaccount"
	CredentialTypeConnector      CredentialType = "connector"
)

// ----- Core model -----

type Credential struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	CredentialType *CredentialType `json:"credential_type,omitempty"`
	HashedData     *string         `json:"hashed_data,omitempty"`
	CreatedAt      *int64          `json:"created_at,omitempty"`
	UpdatedAt      *int64          `json:"updated_at,omitempty"`
}

// ----- Provider-specific credential data shapes (optional helpers) -----

// Microsoft Admin Consent
type MicrosoftAdminConsentSettings struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Google Service Account
type GoogleServiceAccountCredential struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
}

// ----- Requests -----

// CredentialRequest mirrors Python TypedDict (credential_data is a union). Use `any` to allow any JSON-serializable struct/map.
type CredentialRequest struct {
	Name           *string        `json:"name,omitempty"`
	CredentialType CredentialType `json:"credential_type"`
	CredentialData any            `json:"credential_data"`
}

type UpdateCredentialRequest struct {
	Name           *string `json:"name,omitempty"`
	CredentialData any     `json:"credential_data,omitempty"`
}

// ----- Query params -----

// ListCredentialQueryParams: add `url` tags if your EncodeQuery reads them.
type ListCredentialQueryParams struct {
	Limit   *int    `json:"limit,omitempty" url:"limit,omitempty"`
	Offset  *int    `json:"offset,omitempty" url:"offset,omitempty"`
	OrderBy *string `json:"order_by,omitempty" url:"order_by,omitempty"`
	SortBy  *string `json:"sort_by,omitempty" url:"sort_by,omitempty"`
}
