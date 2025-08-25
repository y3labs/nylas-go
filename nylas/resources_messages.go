package nylas

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/y3labs/nylas-go/nylas/models"
)

type MessagesResource struct{ c *Client }

// Expose Smart Compose resource on MessagesResource
func (r *MessagesResource) SmartCompose() *SmartComposeResource {
	return &SmartComposeResource{c: r.c}
}

// List -> GET /v3/grants/{identifier}/messages
func (r *MessagesResource) List(
	ctx context.Context,
	identifier string,
	q *models.ListMessagesQueryParams,
) (*ListResponse[models.Message], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages"
	out, headers, err := DoJSON[ListResponse[models.Message]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Find -> GET /v3/grants/{identifier}/messages/{message_id}
func (r *MessagesResource) Find(
	ctx context.Context,
	identifier, messageID string,
	q *models.FindMessageQueryParams,
) (*Response[models.Message], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/" + url.PathEscape(messageID)
	out, headers, err := DoJSON[Response[models.Message]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/messages/{message_id}
func (r *MessagesResource) Update(
	ctx context.Context,
	identifier, messageID string,
	body models.UpdateMessageRequest,
) (*Response[models.Message], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/" + url.PathEscape(messageID)
	out, headers, err := DoJSON[Response[models.Message]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/grants/{identifier}/messages/{message_id}
func (r *MessagesResource) Destroy(
	ctx context.Context,
	identifier, messageID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/" + url.PathEscape(messageID)
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

// Send -> POST /v3/grants/{identifier}/messages/send
func (r *MessagesResource) Send(
	ctx context.Context,
	identifier string,
	body models.SendMessageRequest,
) (*Response[models.Message], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/send"

	// NOTE: Python re-maps body["from_"] -> "from". Our struct already uses "from".

	out, headers, err := DoJSON[Response[models.Message]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// ListScheduledMessages -> GET /v3/grants/{identifier}/messages/schedules
func (r *MessagesResource) ListScheduledMessages(
	ctx context.Context,
	identifier string,
) (*Response[[]models.ScheduledMessage], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/schedules"
	out, headers, err := DoJSON[Response[[]models.ScheduledMessage]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// FindScheduledMessage -> GET /v3/grants/{identifier}/messages/schedules/{schedule_id}
func (r *MessagesResource) FindScheduledMessage(
	ctx context.Context,
	identifier, scheduleID string,
) (*Response[models.ScheduledMessage], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/schedules/" + url.PathEscape(scheduleID)
	out, headers, err := DoJSON[Response[models.ScheduledMessage]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// StopScheduledMessage -> DELETE /v3/grants/{identifier}/messages/schedules/{schedule_id}
func (r *MessagesResource) StopScheduledMessage(
	ctx context.Context,
	identifier, scheduleID string,
) (*Response[models.StopScheduledMessageResponse], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/schedules/" + url.PathEscape(scheduleID)
	out, headers, err := DoJSON[Response[models.StopScheduledMessageResponse]](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// CleanMessages -> PUT /v3/grants/{identifier}/messages/clean
func (r *MessagesResource) CleanMessages(
	ctx context.Context,
	identifier string,
	body models.CleanMessagesRequest,
) (*ListResponse[models.CleanMessagesResponse], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/clean"
	out, headers, err := DoJSON[ListResponse[models.CleanMessagesResponse]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// SendMultipart mirrors Python's "large attachment" path:
// - a "message" JSON part with Content-Type: application/json
// - file0, file1, ... parts for the attachments
func (r *MessagesResource) SendMultipart(
	ctx context.Context,
	identifier string,
	body any, // e.g., models.SendMessageRequest or map[string]any
	files []FormFile, // use nylas.FormFile defined in httpmultipart.go
) (*Response[models.Message], error) {
	// Marshal message payload to JSON
	msgBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	fields := []FormField{
		{Name: "message", Value: string(msgBytes), ContentType: "application/json"},
	}

	// Ensure field names are file0, file1, ... (Nylas convention)
	normalized := make([]FormFile, len(files))
	for i, f := range files {
		ff := f
		if ff.Field == "" {
			ff.Field = "file" + strconv.Itoa(i)
		}
		normalized[i] = ff
	}

	path := "/v3/grants/" + url.PathEscape(identifier) + "/messages/send"
	resp, err := r.c.doMultipartParts(ctx, http.MethodPost, path, nil, fields, normalized)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out Response[models.Message]
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil && err != io.EOF {
		return nil, err
	}
	out.Headers = resp.Header
	if out.RequestID == "" {
		out.RequestID = resp.Header.Get("x-request-id")
	}
	return &out, nil
}
