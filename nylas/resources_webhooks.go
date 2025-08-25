package nylas

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type WebhooksResource struct{ c *Client }

// List -> GET /v3/webhooks
func (r *WebhooksResource) List(ctx context.Context) (*ListResponse[models.Webhook], error) {
	out, headers, err := DoJSON[ListResponse[models.Webhook]](r.c, ctx, http.MethodGet, "/v3/webhooks", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Find -> GET /v3/webhooks/{webhook_id}
func (r *WebhooksResource) Find(ctx context.Context, webhookID string) (*Response[models.Webhook], error) {
	path := "/v3/webhooks/" + url.PathEscape(webhookID)
	out, headers, err := DoJSON[Response[models.Webhook]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/webhooks
func (r *WebhooksResource) Create(ctx context.Context, req models.CreateWebhookRequest) (*Response[models.WebhookWithSecret], error) {
	out, headers, err := DoJSON[Response[models.WebhookWithSecret]](r.c, ctx, http.MethodPost, "/v3/webhooks", nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/webhooks/{webhook_id}
func (r *WebhooksResource) Update(ctx context.Context, webhookID string, req models.UpdateWebhookRequest) (*Response[models.Webhook], error) {
	path := "/v3/webhooks/" + url.PathEscape(webhookID)
	out, headers, err := DoJSON[Response[models.Webhook]](r.c, ctx, http.MethodPut, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/webhooks/{webhook_id}
// Returns a specialized envelope (not the common DeleteResponse).
func (r *WebhooksResource) Destroy(ctx context.Context, webhookID string) (*models.WebhookDeleteResponse, error) {
	path := "/v3/webhooks/" + url.PathEscape(webhookID)
	out, headers, err := DoJSON[models.WebhookDeleteResponse](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	// ensure request id is populated even if body didn't include it
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// RotateSecret -> PUT /v3/webhooks/{webhook_id}/rotate-secret
func (r *WebhooksResource) RotateSecret(ctx context.Context, webhookID string) (*Response[models.WebhookWithSecret], error) {
	path := "/v3/webhooks/" + url.PathEscape(webhookID) + "/rotate-secret"
	empty := map[string]any{} // server expects a JSON object
	out, headers, err := DoJSON[Response[models.WebhookWithSecret]](r.c, ctx, http.MethodPut, path, nil, empty, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// IPAddresses -> GET /v3/webhooks/ip-addresses
func (r *WebhooksResource) IPAddresses(ctx context.Context) (*Response[models.WebhookIpAddressesResponse], error) {
	out, headers, err := DoJSON[Response[models.WebhookIpAddressesResponse]](r.c, ctx, http.MethodGet, "/v3/webhooks/ip-addresses", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

func ExtractChallengeParameter(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	v := q.Get("challenge")
	if v == "" {
		return "", errors.New("invalid URL or no challenge parameter found")
	}
	return v, nil
}
