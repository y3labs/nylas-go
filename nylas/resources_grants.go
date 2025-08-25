package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type GrantsResource struct{ c *Client }

// Matches Python's ListGrantsQueryParams
type ListGrantsParams struct {
	Limit       *int    `query:"limit"`       // default 10, max 200 (server-enforced)
	Offset      *int    `query:"offset"`      // offset-based pagination
	SortBy      *string `query:"sortBy"`      // field to sort by
	OrderBy     *string `query:"orderBy"`     // asc|desc
	Since       *int64  `query:"since"`       // unix timestamp
	Before      *int64  `query:"before"`      // unix timestamp
	Email       *string `query:"email"`       // filter by email
	GrantStatus *string `query:"grantStatus"` // filter by status
	IP          *string `query:"ip"`          // filter by IP
	Provider    *string `query:"provider"`    // provider (string enum upstream)
}

func (r *GrantsResource) List(ctx context.Context, params *ListGrantsParams) (*ListResponse[models.Grant], error) {
	q := EncodeQuery(params)
	out, headers, err := DoJSON[ListResponse[models.Grant]](r.c, ctx, http.MethodGet, "/v3/grants", q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *GrantsResource) Get(ctx context.Context, grantID string) (*Response[models.Grant], error) {
	path := "/v3/grants/" + url.PathEscape(grantID)
	out, headers, err := DoJSON[Response[models.Grant]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Matches Python's UpdateGrantRequest
type UpdateGrantRequest struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Scope    []string               `json:"scope,omitempty"`
}

func (r *GrantsResource) Update(ctx context.Context, grantID string, req UpdateGrantRequest) (*Response[models.Grant], error) {
	path := "/v3/grants/" + url.PathEscape(grantID)
	out, headers, err := DoJSON[Response[models.Grant]](r.c, ctx, http.MethodPatch, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *GrantsResource) Delete(ctx context.Context, grantID string) error {
	path := "/v3/grants/" + url.PathEscape(grantID)
	_, _, err := DoJSON[map[string]any](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	return err
}
