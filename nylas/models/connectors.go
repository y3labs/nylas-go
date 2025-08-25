package models

// Connector mirrors nylas.models.connectors.Connector
type Connector struct {
	Provider Provider       `json:"provider"`
	Settings map[string]any `json:"settings,omitempty"`
	Scope    []string       `json:"scope,omitempty"`
}

// BaseCreateConnectorRequest (kept for parity; concrete requests below include Provider directly)
type BaseCreateConnectorRequest struct {
	Provider Provider `json:"provider"`
}

// GoogleCreateConnectorSettings mirrors the Python settings object.
type GoogleCreateConnectorSettings struct {
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	TopicName    *string `json:"topic_name,omitempty"`
}

// MicrosoftCreateConnectorSettings mirrors the Python settings object.
type MicrosoftCreateConnectorSettings struct {
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	Tenant       *string `json:"tenant,omitempty"`
}

// GoogleCreateConnectorRequest corresponds to CreateConnectorRequest variant for Google.
type GoogleCreateConnectorRequest struct {
	Provider Provider                      `json:"provider"` // should be "google"
	Settings GoogleCreateConnectorSettings `json:"settings"`
	Scope    []string                      `json:"scope,omitempty"`
}

// MicrosoftCreateConnectorRequest corresponds to CreateConnectorRequest variant for Microsoft.
type MicrosoftCreateConnectorRequest struct {
	Provider Provider                         `json:"provider"` // should be "microsoft"
	Settings MicrosoftCreateConnectorSettings `json:"settings"`
	Scope    []string                         `json:"scope,omitempty"`
}

// ImapCreateConnectorRequest corresponds to the IMAP variant (no extra settings in Python model).
type ImapCreateConnectorRequest struct {
	Provider Provider `json:"provider"` // should be "imap"
}

// VirtualCalendarsCreateConnectorRequest corresponds to the Virtual Calendars variant.
type VirtualCalendarsCreateConnectorRequest struct {
	Provider Provider `json:"provider"` // should be "virtual-calendar"
}

// UpdateConnectorRequest mirrors the Python TypedDict (all optional).
type UpdateConnectorRequest struct {
	Name     *string        `json:"name,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
	Scope    []string       `json:"scope,omitempty"`
}

// ListConnectorQueryParams mirrors ListQueryParams (limit/page_token).
// Include `url` tags if you use a generic query encoder that reads them.
type ListConnectorQueryParams struct {
	Limit     *int    `json:"limit,omitempty" url:"limit,omitempty"`
	PageToken *string `json:"page_token,omitempty" url:"page_token,omitempty"`
}
