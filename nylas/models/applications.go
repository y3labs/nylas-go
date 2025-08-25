package models

type Environment string

const (
	ProdEnv Environment = "production"
	StagEnv Environment = "staging"
	DevEnv  Environment = "development"
	SandEnv Environment = "sandbox"
)

type Branding struct {
	Name        string `json:"name"`
	IconURL     string `json:"icon_url,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
	Description string `json:"description,omitempty"`
}

type HostedAuthentication struct {
	BackgroundImageURL string `json:"background_image_url,omitempty"`
	Alignment          string `json:"alignment,omitempty"`
	ColorPrimary       string `json:"color_primary,omitempty"`
	ColorSecondary     string `json:"color_secondary,omitempty"`
	Title              string `json:"title,omitempty"`
	Subtitle           string `json:"subtitle,omitempty"`
	BackgroundColor    string `json:"background_color,omitempty"`
	Spacing            int    `json:"spacing,omitempty"`
}

type Application struct {
	ApplicationID        string               `json:"application_id"`
	OrganizationID       string               `json:"organization_id"`
	Region               Region               `json:"region"`
	Environment          Environment          `json:"environment"`
	Branding             Branding             `json:"branding"`
	HostedAuthentication HostedAuthentication `json:"hosted_authentication,omitempty"`
	CallbackURIs         []RedirectURI        `json:"callback_uris,omitempty"`
}
