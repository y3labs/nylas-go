package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type SmartComposeResource struct{ c *Client }

// Optional top-level accessor (so users can do c.SmartCompose())
/*
func (c *Client) SmartCompose() *SmartComposeResource {
	return &SmartComposeResource{c: c}
}
*/

// ComposeMessage generates a suggestion for a *new* message.
// POST /v3/grants/{grantID}/messages/smart-compose
func (r *SmartComposeResource) ComposeMessage(
	ctx context.Context,
	grantID string,
	req models.ComposeMessageRequest,
) (*Response[models.ComposeMessageResponse], error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/messages/smart-compose"
	out, headers, err := DoJSON[Response[models.ComposeMessageResponse]](r.c, ctx, http.MethodPost, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// ComposeMessageReply generates a suggestion for a *reply* to an existing message.
// POST /v3/grants/{grantID}/messages/{messageID}/smart-compose
func (r *SmartComposeResource) ComposeMessageReply(
	ctx context.Context,
	grantID, messageID string,
	req models.ComposeMessageRequest,
) (*Response[models.ComposeMessageResponse], error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/messages/" + url.PathEscape(messageID) + "/smart-compose"
	out, headers, err := DoJSON[Response[models.ComposeMessageResponse]](r.c, ctx, http.MethodPost, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}
