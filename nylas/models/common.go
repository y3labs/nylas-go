package models

// EmailName represents a name + email address pair.
type EmailName struct {
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

type Region string

const (
	RegionUS Region = "us"
	RegionEU Region = "eu"
)

type ListCore struct {
	Limit     *int    `query:"limit"`
	PageToken *string `query:"page_token"`
	Select    *string `query:"select"`
}
