package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type NotetakersResource struct{ c *Client }

// Helper to build the base path depending on whether a grant identifier is provided.
// identifier == "" → global notetakers; otherwise grant-scoped.
func ntBasePath(identifier string) string {
	if identifier == "" {
		return "/v3/notetakers"
	}
	return "/v3/grants/" + url.PathEscape(identifier) + "/notetakers"
}

// List -> GET /v3/notetakers OR /v3/grants/{identifier}/notetakers
func (r *NotetakersResource) List(
	ctx context.Context,
	identifier string, // optional: pass "" for global scope
	q *models.ListNotetakerQueryParams,
) (*ListResponse[models.Notetaker], error) {
	path := ntBasePath(identifier)
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	out, headers, err := DoJSON[ListResponse[models.Notetaker]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get -> GET /v3/notetakers/{id} OR /v3/grants/{identifier}/notetakers/{id}
func (r *NotetakersResource) Get(
	ctx context.Context,
	notetakerID string,
	identifier string, // optional
	q *models.FindNotetakerQueryParams,
) (*Response[models.Notetaker], error) {
	path := ntBasePath(identifier) + "/" + url.PathEscape(notetakerID)
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	out, headers, err := DoJSON[Response[models.Notetaker]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Invite -> POST /v3/notetakers OR /v3/grants/{identifier}/notetakers
func (r *NotetakersResource) Invite(
	ctx context.Context,
	body models.InviteNotetakerRequest,
	identifier string, // optional
) (*Response[models.Notetaker], error) {
	path := ntBasePath(identifier)
	out, headers, err := DoJSON[Response[models.Notetaker]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PATCH /v3/notetakers/{id} OR /v3/grants/{identifier}/notetakers/{id}
func (r *NotetakersResource) Update(
	ctx context.Context,
	notetakerID string,
	body models.UpdateNotetakerRequest,
	identifier string, // optional
) (*Response[models.Notetaker], error) {
	path := ntBasePath(identifier) + "/" + url.PathEscape(notetakerID)
	out, headers, err := DoJSON[Response[models.Notetaker]](r.c, ctx, http.MethodPatch, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Leave -> POST /v3/notetakers/{id}/leave OR /v3/grants/{identifier}/notetakers/{id}/leave
func (r *NotetakersResource) Leave(
	ctx context.Context,
	notetakerID string,
	identifier string, // optional
) (*Response[models.NotetakerLeaveResponse], error) {
	path := ntBasePath(identifier) + "/" + url.PathEscape(notetakerID) + "/leave"
	out, headers, err := DoJSON[Response[models.NotetakerLeaveResponse]](r.c, ctx, http.MethodPost, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// GetMedia -> GET /v3/notetakers/{id}/media OR /v3/grants/{identifier}/notetakers/{id}/media
func (r *NotetakersResource) GetMedia(
	ctx context.Context,
	notetakerID string,
	identifier string, // optional
) (*Response[models.NotetakerMedia], error) {
	path := ntBasePath(identifier) + "/" + url.PathEscape(notetakerID) + "/media"
	out, headers, err := DoJSON[Response[models.NotetakerMedia]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Cancel -> DELETE /v3/notetakers/{id}/cancel OR /v3/grants/{identifier}/notetakers/{id}/cancel
func (r *NotetakersResource) Cancel(
	ctx context.Context,
	notetakerID string,
	identifier string, // optional
) (*DeleteResponse, error) {
	path := ntBasePath(identifier) + "/" + url.PathEscape(notetakerID) + "/cancel"
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
