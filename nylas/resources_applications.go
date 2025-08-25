package nylas

import (
	"context"
	"net/http"

	"github.com/y3labs/nylas-go/nylas/models"
)

// ApplicationsResource mirrors the Python Applications(Resource) surface.
type ApplicationsResource struct{ c *Client }

// RedirectURIs provides access to /v3/applications/redirect-uris, matching Python's .redirect_uris property.
func (r *ApplicationsResource) RedirectURIs() *RedirectURIsResource {
	return &RedirectURIsResource{r.c}
}

// Info fetches GET /v3/applications and returns application details.
func (r *ApplicationsResource) Info(ctx context.Context) (*Response[models.Application], error) {
	const path = "/v3/applications"

	out, headers, err := DoJSON[Response[models.Application]](
		r.c, ctx, http.MethodGet, path,
		nil, // query
		nil, // body
		nil, // extra headers / per-call overrides
	)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
