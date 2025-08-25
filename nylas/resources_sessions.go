package nylas

import (
	"context"
	"net/http"

	"github.com/y3labs/nylas-go/nylas/models"
)

type SessionsResource struct{ c *Client }

// Create -> POST /v3/scheduling/sessions
func (r *SessionsResource) Create(
	ctx context.Context,
	body models.CreateSessionRequest,
) (*Response[models.Session], error) {
	path := "/v3/scheduling/sessions"
	out, headers, err := DoJSON[Response[models.Session]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/scheduling/sessions/{session_id}
func (r *SessionsResource) Destroy(
	ctx context.Context,
	sessionID string,
) (*DeleteResponse, error) {
	path := "/v3/scheduling/sessions/" + sessionID // session_id is opaque; do not escape `/`
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
