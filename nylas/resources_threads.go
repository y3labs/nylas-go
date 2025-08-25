package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type ThreadsResource struct{ c *Client }

// List -> GET /v3/grants/{identifier}/threads
func (r *ThreadsResource) List(
	ctx context.Context,
	identifier string,
	q *models.ListThreadsQueryParams,
) (*ListResponse[models.Thread], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/threads"
	out, headers, err := DoJSON[ListResponse[models.Thread]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Find -> GET /v3/grants/{identifier}/threads/{thread_id}
func (r *ThreadsResource) Find(
	ctx context.Context,
	identifier, threadID string,
) (*Response[models.Thread], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/threads/" + url.PathEscape(threadID)
	out, headers, err := DoJSON[Response[models.Thread]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/threads/{thread_id}
func (r *ThreadsResource) Update(
	ctx context.Context,
	identifier, threadID string,
	body models.UpdateThreadRequest,
) (*Response[models.Thread], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/threads/" + url.PathEscape(threadID)
	out, headers, err := DoJSON[Response[models.Thread]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/grants/{identifier}/threads/{thread_id}
func (r *ThreadsResource) Destroy(
	ctx context.Context,
	identifier, threadID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/threads/" + url.PathEscape(threadID)
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
