package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

// CredentialsResource mirrors nylas.resources.credentials.Credentials
type CredentialsResource struct{ c *Client }

// Optional top-level accessor if you want one (you already expose it under Connectors().Credentials()).
// func (c *Client) Credentials() *CredentialsResource { return &CredentialsResource{c} }

// List -> GET /v3/connectors/{provider}/creds
func (r *CredentialsResource) List(
	ctx context.Context,
	provider models.Provider,
	q *models.ListCredentialQueryParams,
) (*ListResponse[models.Credential], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider)) + "/creds"

	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}

	out, headers, err := DoJSON[ListResponse[models.Credential]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get -> GET /v3/connectors/{provider}/creds/{credential_id}
func (r *CredentialsResource) Get(
	ctx context.Context,
	provider models.Provider,
	credentialID string,
) (*Response[models.Credential], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider)) + "/creds/" + url.PathEscape(credentialID)

	out, headers, err := DoJSON[Response[models.Credential]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/connectors/{provider}/creds
func (r *CredentialsResource) Create(
	ctx context.Context,
	provider models.Provider,
	body models.CredentialRequest,
) (*Response[models.Credential], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider)) + "/creds"

	out, headers, err := DoJSON[Response[models.Credential]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PATCH /v3/connectors/{provider}/creds/{credential_id}
func (r *CredentialsResource) Update(
	ctx context.Context,
	provider models.Provider,
	credentialID string,
	body models.UpdateCredentialRequest,
) (*Response[models.Credential], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider)) + "/creds/" + url.PathEscape(credentialID)

	out, headers, err := DoJSON[Response[models.Credential]](r.c, ctx, http.MethodPatch, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Delete -> DELETE /v3/connectors/{provider}/creds/{credential_id}
func (r *CredentialsResource) Delete(
	ctx context.Context,
	provider models.Provider,
	credentialID string,
) (*DeleteResponse, error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider)) + "/creds/" + url.PathEscape(credentialID)

	out, headers, err := DoJSON[DeleteResponse](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
