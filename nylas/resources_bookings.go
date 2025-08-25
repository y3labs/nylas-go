package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type BookingsResource struct{ c *Client }

// Find -> GET /v3/scheduling/bookings/{booking_id}
func (r *BookingsResource) Find(
	ctx context.Context,
	bookingID string,
	q *models.FindBookingQueryParams,
) (*Response[models.Booking], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/scheduling/bookings/" + url.PathEscape(bookingID)
	out, headers, err := DoJSON[Response[models.Booking]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/scheduling/bookings
func (r *BookingsResource) Create(
	ctx context.Context,
	body models.CreateBookingRequest,
	q *models.CreateBookingQueryParams,
) (*Response[models.Booking], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/scheduling/bookings"
	out, headers, err := DoJSON[Response[models.Booking]](r.c, ctx, http.MethodPost, path, query, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Confirm -> PUT /v3/scheduling/bookings/{booking_id}
func (r *BookingsResource) Confirm(
	ctx context.Context,
	bookingID string,
	body models.ConfirmBookingRequest,
	q *models.ConfirmBookingQueryParams,
) (*Response[models.Booking], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/scheduling/bookings/" + url.PathEscape(bookingID)
	out, headers, err := DoJSON[Response[models.Booking]](r.c, ctx, http.MethodPut, path, query, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Reschedule -> PATCH /v3/scheduling/bookings/{booking_id}
func (r *BookingsResource) Reschedule(
	ctx context.Context,
	bookingID string,
	body models.RescheduleBookingRequest,
	q *models.RescheduleBookingQueryParams,
) (*Response[models.Booking], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/scheduling/bookings/" + url.PathEscape(bookingID)
	out, headers, err := DoJSON[Response[models.Booking]](r.c, ctx, "PATCH", path, query, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/scheduling/bookings/{booking_id} (with body)
func (r *BookingsResource) Destroy(
	ctx context.Context,
	bookingID string,
	body models.DeleteBookingRequest,
	q *models.DestroyBookingQueryParams,
) (*DeleteResponse, error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/scheduling/bookings/" + url.PathEscape(bookingID)
	out, headers, err := DoJSON[DeleteResponse](r.c, ctx, http.MethodDelete, path, query, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}
