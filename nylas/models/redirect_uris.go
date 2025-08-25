package models

type RedirectURISettings struct {
	Origin                     string `json:"origin,omitempty"`
	BundleID                   string `json:"bundle_id,omitempty"`
	AppStoreID                 string `json:"app_store_id,omitempty"`
	TeamID                     string `json:"team_id,omitempty"`
	PackageName                string `json:"package_name,omitempty"`
	SHA1CertificateFingerprint string `json:"sha1_certificate_fingerprint,omitempty"`
}

type RedirectURI struct {
	ID       string              `json:"id"`
	URL      string              `json:"url"`
	Platform string              `json:"platform"`
	Settings RedirectURISettings `json:"settings,omitempty"`
}

type WritableRedirectURISettings struct {
	Origin                     string `json:"origin,omitempty"`
	BundleID                   string `json:"bundle_id,omitempty"`
	AppStoreID                 string `json:"app_store_id,omitempty"`
	TeamID                     string `json:"team_id,omitempty"`
	PackageName                string `json:"package_name,omitempty"`
	SHA1CertificateFingerprint string `json:"sha1_certificate_fingerprint,omitempty"`
}

type CreateRedirectURIRequest struct {
	URL      string                      `json:"url"`
	Platform string                      `json:"platform"`
	Settings WritableRedirectURISettings `json:"settings,omitempty"`
}

type UpdateRedirectURIRequest struct {
	URL      string                      `json:"url,omitempty"`
	Platform string                      `json:"platform,omitempty"`
	Settings WritableRedirectURISettings `json:"settings,omitempty"`
}
