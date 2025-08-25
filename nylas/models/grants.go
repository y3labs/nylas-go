package models

type Grant struct {
	ID             string                 `json:"id"`
	Provider       string                 `json:"provider"`               // OAuth provider (string enum upstream)
	Scope          []string               `json:"scope"`                  // scopes for the grant
	AccountID      *string                `json:"account_id,omitempty"`   // migrated v2 account id (optional)
	GrantStatus    *string                `json:"grant_status,omitempty"` // current status
	Email          *string                `json:"email,omitempty"`
	UserAgent      *string                `json:"user_agent,omitempty"`
	IP             *string                `json:"ip,omitempty"`
	State          *string                `json:"state,omitempty"`
	CreatedAt      *int64                 `json:"created_at,omitempty"` // unix seconds
	UpdatedAt      *int64                 `json:"updated_at,omitempty"` // unix seconds
	ProviderUserID *string                `json:"provider_user_id,omitempty"`
	Settings       map[string]interface{} `json:"settings,omitempty"` // provider-specific settings
}
