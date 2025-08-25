package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type CalendarsResource struct{ c *Client }

// Optional accessor if not present yet:
// func (c *Client) Calendars() *CalendarsResource { return &CalendarsResource{c} }

// List -> GET /v3/grants/{identifier}/calendars
func (r *CalendarsResource) List(
	ctx context.Context,
	identifier string,
	q *models.ListCalendarsQueryParams,
) (*ListResponse[models.Calendar], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars"
	out, headers, err := DoJSON[ListResponse[models.Calendar]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get -> GET /v3/grants/{identifier}/calendars/{calendar_id}
// Pass calendar_id="primary" to target primary calendar.
func (r *CalendarsResource) Get(
	ctx context.Context,
	identifier, calendarID string,
	q *models.FindCalendarQueryParams,
) (*Response[models.Calendar], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars/" + url.PathEscape(calendarID)
	out, headers, err := DoJSON[Response[models.Calendar]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/grants/{identifier}/calendars
func (r *CalendarsResource) Create(
	ctx context.Context,
	identifier string,
	body models.CreateCalendarRequest,
) (*Response[models.Calendar], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars"
	out, headers, err := DoJSON[Response[models.Calendar]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/calendars/{calendar_id}
func (r *CalendarsResource) Update(
	ctx context.Context,
	identifier, calendarID string,
	body models.UpdateCalendarRequest,
) (*Response[models.Calendar], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars/" + url.PathEscape(calendarID)
	out, headers, err := DoJSON[Response[models.Calendar]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Delete -> DELETE /v3/grants/{identifier}/calendars/{calendar_id}
func (r *CalendarsResource) Delete(
	ctx context.Context,
	identifier, calendarID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars/" + url.PathEscape(calendarID)
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

// GetAvailability (global) -> POST /v3/calendars/availability
// Mirrors Calendars.get_availability in Python.
func (r *CalendarsResource) GetAvailability(
	ctx context.Context,
	body models.GetAvailabilityRequest,
) (*Response[models.GetAvailabilityResponse], error) {
	const path = "/v3/calendars/availability"
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

// GetFreeBusy (grant-scoped) -> POST /v3/grants/{identifier}/calendars/free-busy
// Decodes Python's Union[FreeBusy, FreeBusyError] into a slice of items where Error may be set.
func (r *CalendarsResource) GetFreeBusy(
	ctx context.Context,
	identifier string,
	body models.GetFreeBusyRequest,
) (*Response[[]models.GetFreeBusyResponseItem], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/calendars/free-busy"
	out, headers, err := DoJSON[Response[[]models.GetFreeBusyResponseItem]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
