package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

// ConnectorsResource mirrors nylas.resources.connectors.Connectors
type ConnectorsResource struct{ c *Client }

// Credentials sub-resource (Python exposes .credentials property). Requires you to have a CredentialsResource.
func (r *ConnectorsResource) Credentials() *CredentialsResource { return &CredentialsResource{r.c} }

// List  -> GET /v3/connectors
// Accepts optional query params: limit, page_token
func (r *ConnectorsResource) List(ctx context.Context, q *models.ListConnectorQueryParams) (*ListResponse[models.Connector], error) {
	path := "/v3/connectors"
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	out, headers, err := DoJSON[ListResponse[models.Connector]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get   -> GET /v3/connectors/{provider}
func (r *ConnectorsResource) Get(ctx context.Context, provider models.Provider) (*Response[models.Connector], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider))
	out, headers, err := DoJSON[Response[models.Connector]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/connectors
// The request body is one of:
//   - models.GoogleCreateConnectorRequest
//   - models.MicrosoftCreateConnectorRequest
//   - models.ImapCreateConnectorRequest
//   - models.VirtualCalendarsCreateConnectorRequest
func (r *ConnectorsResource) Create(ctx context.Context, body any) (*Response[models.Connector], error) {
	path := "/v3/connectors"
	out, headers, err := DoJSON[Response[models.Connector]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PATCH /v3/connectors/{provider}
func (r *ConnectorsResource) Update(ctx context.Context, provider models.Provider, body models.UpdateConnectorRequest) (*Response[models.Connector], error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider))
	out, headers, err := DoJSON[Response[models.Connector]](r.c, ctx, http.MethodPatch, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Delete -> DELETE /v3/connectors/{provider}
func (r *ConnectorsResource) Delete(ctx context.Context, provider models.Provider) (*DeleteResponse, error) {
	path := "/v3/connectors/" + url.PathEscape(string(provider))
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
