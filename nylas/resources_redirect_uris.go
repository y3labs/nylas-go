package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type RedirectURIsResource struct{ c *Client }

// Convenience accessor (add to client.go as shown below)
// func (c *Client) RedirectURIs() *RedirectURIsResource { return &RedirectURIsResource{c} }

func (r *RedirectURIsResource) List(ctx context.Context) (*ListResponse[models.RedirectURI], error) {
	path := "/v3/applications/redirect-uris"
	out, headers, err := DoJSON[ListResponse[models.RedirectURI]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *RedirectURIsResource) Get(ctx context.Context, redirectURIID string) (*Response[models.RedirectURI], error) {
	path := "/v3/applications/redirect-uris/" + url.PathEscape(redirectURIID)
	out, headers, err := DoJSON[Response[models.RedirectURI]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *RedirectURIsResource) Create(ctx context.Context, req models.CreateRedirectURIRequest) (*Response[models.RedirectURI], error) {
	path := "/v3/applications/redirect-uris"
	out, headers, err := DoJSON[Response[models.RedirectURI]](r.c, ctx, http.MethodPost, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Python uses UpdatableApiResource with default method=PUT for this endpoint.
// We mirror that with http.MethodPut here for parity.
func (r *RedirectURIsResource) Update(ctx context.Context, redirectURIID string, req models.UpdateRedirectURIRequest) (*Response[models.RedirectURI], error) {
	path := "/v3/applications/redirect-uris/" + url.PathEscape(redirectURIID)
	out, headers, err := DoJSON[Response[models.RedirectURI]](r.c, ctx, http.MethodPut, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *RedirectURIsResource) Delete(ctx context.Context, redirectURIID string) (*DeleteResponse, error) {
	path := "/v3/applications/redirect-uris/" + url.PathEscape(redirectURIID)
	out, headers, err := DoJSON[DeleteResponse](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	// The body includes request_id; this just ensures it’s also available from headers if needed.
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
