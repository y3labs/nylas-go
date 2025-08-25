package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type AvailabilityResource struct{ c *Client }

// Check (global) -> POST /v3/availability
func (r *AvailabilityResource) Check(
	ctx context.Context,
	body models.GetAvailabilityRequest,
) (*Response[models.GetAvailabilityResponse], error) {
	const path = "/v3/availability"
	out, headers, err := DoJSON[Response[models.GetAvailabilityResponse]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// CheckForGrant (grant-scoped) -> POST /v3/grants/{identifier}/availability
// Some deployments expose availability under a grant scope; this keeps both options available.
func (r *AvailabilityResource) CheckForGrant(
	ctx context.Context,
	grantIdentifier string,
	body models.GetAvailabilityRequest,
) (*Response[models.GetAvailabilityResponse], error) {
	path := "/v3/grants/" + url.PathEscape(grantIdentifier) + "/availability"
	out, headers, err := DoJSON[Response[models.GetAvailabilityResponse]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
