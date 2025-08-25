package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type DraftsResource struct{ c *Client }

// List -> GET /v3/grants/{identifier}/drafts
func (r *DraftsResource) List(
	ctx context.Context,
	identifier string,
	q *models.ListDraftsQueryParams,
) (*ListResponse[models.Draft], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts"
	out, headers, err := DoJSON[ListResponse[models.Draft]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Find -> GET /v3/grants/{identifier}/drafts/{draft_id}
func (r *DraftsResource) Find(
	ctx context.Context,
	identifier, draftID string,
	q *models.FindDraftQueryParams,
) (*Response[models.Draft], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts/" + url.PathEscape(draftID)
	out, headers, err := DoJSON[Response[models.Draft]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/grants/{identifier}/drafts
func (r *DraftsResource) Create(
	ctx context.Context,
	identifier string,
	body models.CreateDraftRequest,
) (*Response[models.Draft], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts"

	// NOTE: Python SDK switches to multipart for >3MB. We currently do JSON only.
	// TODO: Add multipart builder for large attachments if needed.

	out, headers, err := DoJSON[Response[models.Draft]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/drafts/{draft_id}
func (r *DraftsResource) Update(
	ctx context.Context,
	identifier, draftID string,
	body models.UpdateDraftRequest,
) (*Response[models.Draft], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts/" + url.PathEscape(draftID)

	// Python uses UpdatableApiResource default (PUT).
	out, headers, err := DoJSON[Response[models.Draft]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/grants/{identifier}/drafts/{draft_id}
func (r *DraftsResource) Destroy(
	ctx context.Context,
	identifier, draftID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts/" + url.PathEscape(draftID)
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

// Send -> POST /v3/grants/{identifier}/drafts/{draft_id}
// (Python sends to the draft endpoint without "/send")
func (r *DraftsResource) Send(
	ctx context.Context,
	identifier, draftID string,
) (*Response[models.Message], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/drafts/" + url.PathEscape(draftID)
	out, headers, err := DoJSON[Response[models.Message]](r.c, ctx, http.MethodPost, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
