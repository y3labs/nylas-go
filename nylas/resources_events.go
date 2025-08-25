package nylas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type EventsResource struct{ c *Client }

// ----- Query param shapes (mirror Python SDK) -----

type ListEventsParams struct {
	CalendarID      string             `query:"calendar_id"` // required
	ShowCancelled   *bool              `query:"show_cancelled"`
	Title           *string            `query:"title"`
	Description     *string            `query:"description"`
	Location        *string            `query:"location"`
	Start           *int64             `query:"start"` // unix seconds
	End             *int64             `query:"end"`   // unix seconds
	ExpandRecurring *bool              `query:"expand_recurring"`
	Busy            *bool              `query:"busy"`
	OrderBy         *string            `query:"order_by"` // currently "start"
	EventType       []models.EventType `query:"event_type"`
	MasterEventID   *string            `query:"master_event_id"`
	Select          *string            `query:"select"`
	TentativeAsBusy *bool              `query:"tentative_as_busy"`
	Limit           *int               `query:"limit"`
	PageToken       *string            `query:"page_token"`
	// metadata_pair is a k/v filter; we add it manually to the query to avoid weird reflection rules.
	MetadataPair map[string]string `query:"-"`
}

func (p *ListEventsParams) toValues() url.Values {
	q := EncodeQuery(p)
	// inject metadata_pair[key]=value
	for k, v := range p.MetadataPair {
		q.Add(fmt.Sprintf("metadata_pair[%s]", k), v)
	}
	return q
}

type ListImportEventsParams struct {
	CalendarID string  `query:"calendar_id"` // required
	Start      *int64  `query:"start"`
	End        *int64  `query:"end"`
	Select     *string `query:"select"`
	PageToken  *string `query:"page_token"`
	Limit      *int    `query:"limit"`
}

type FindEventParams struct {
	CalendarID      string `query:"calendar_id"` // required
	TentativeAsBusy *bool  `query:"tentative_as_busy"`
}

type CreateEventParams struct {
	CalendarID         string `query:"calendar_id"` // required
	NotifyParticipants *bool  `query:"notify_participants"`
	TentativeAsBusy    *bool  `query:"tentative_as_busy"`
}

type UpdateEventParams = CreateEventParams

type DestroyEventParams struct {
	CalendarID         string `query:"calendar_id"`
	NotifyParticipants *bool  `query:"notify_participants"`
}

type SendRSVPParams struct {
	CalendarID string `query:"calendar_id"` // required
}

// ----- Requests (mirror Python) -----

// CreateEventRequest mirrors nylas.models.events.CreateEventRequest
type CreateEventRequest struct {
	When             CreateWhen            `json:"when"` // helper type below
	Title            *string               `json:"title,omitempty"`
	Busy             *bool                 `json:"busy,omitempty"`
	Description      *string               `json:"description,omitempty"`
	Location         *string               `json:"location,omitempty"`
	Conferencing     *models.Conferencing  `json:"conferencing,omitempty"`
	Reminders        *models.Reminders     `json:"reminders,omitempty"`
	Metadata         map[string]any        `json:"metadata,omitempty"`
	Participants     []models.Participant  `json:"participants,omitempty"`
	Recurrence       []string              `json:"recurrence,omitempty"`
	Visibility       *models.Visibility    `json:"visibility,omitempty"`
	Capacity         *int                  `json:"capacity,omitempty"`
	HideParticipants *bool                 `json:"hide_participants,omitempty"`
	Notetaker        *CreateEventNotetaker `json:"notetaker,omitempty"`
}

type CreateEventNotetaker struct {
	Name            *string                          `json:"name,omitempty"`
	MeetingSettings *models.NotetakerMeetingSettings `json:"meeting_settings,omitempty"`
}

// UpdateEventRequest mirrors nylas.models.events.UpdateEventRequest
type UpdateEventRequest struct {
	When             *UpdateWhen            `json:"when,omitempty"`
	Title            *string                `json:"title,omitempty"`
	Busy             *bool                  `json:"busy,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Location         *string                `json:"location,omitempty"`
	Conferencing     *models.Conferencing   `json:"conferencing,omitempty"`
	Reminders        *models.Reminders      `json:"reminders,omitempty"`
	Metadata         map[string]any         `json:"metadata,omitempty"`
	Participants     []models.Participant   `json:"participants,omitempty"`
	Recurrence       []string               `json:"recurrence,omitempty"`
	Visibility       *models.Visibility     `json:"visibility,omitempty"`
	Capacity         *int                   `json:"capacity,omitempty"`
	HideParticipants *bool                  `json:"hide_participants,omitempty"`
	Notetaker        *models.EventNotetaker `json:"notetaker,omitempty"`
}

// Send RSVP
type SendRSVPRequest struct {
	Status models.SendRSVPStatus `json:"status"` // "yes" | "no" | "maybe"
}

type RequestIDResponse struct {
	RequestID string `json:"request_id"`
}

// ----- helper “CreateWhen/UpdateWhen” wrappers for requests -----

// CreateWhen is a convenience alias to reuse Event.When's marshaling behavior
// while letting request structs control optionality more easily.
type createWhenAlias = models.When
type updateWhenAlias = models.When

// exported wrappers to make tags look like Python's type names
type CreateWhen = createWhenAlias
type UpdateWhen = updateWhenAlias

// ----- Resource methods -----

func (r *EventsResource) List(ctx context.Context, grantID string, params ListEventsParams) (*ListResponse[models.Event], error) {
	q := params.toValues()
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events"
	out, headers, err := DoJSON[ListResponse[models.Event]](r.c, ctx, http.MethodGet, path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *EventsResource) ListImport(ctx context.Context, grantID string, params ListImportEventsParams) (*ListResponse[models.Event], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events/import"
	out, headers, err := DoJSON[ListResponse[models.Event]](r.c, ctx, http.MethodGet, path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *EventsResource) Get(ctx context.Context, grantID, eventID string, params FindEventParams) (*Response[models.Event], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events/" + url.PathEscape(eventID)
	out, headers, err := DoJSON[Response[models.Event]](r.c, ctx, http.MethodGet, path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *EventsResource) Create(ctx context.Context, grantID string, req CreateEventRequest, params CreateEventParams) (*Response[models.Event], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events"
	out, headers, err := DoJSON[Response[models.Event]](r.c, ctx, http.MethodPost, path, q, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *EventsResource) Update(ctx context.Context, grantID, eventID string, req UpdateEventRequest, params UpdateEventParams) (*Response[models.Event], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events/" + url.PathEscape(eventID)
	out, headers, err := DoJSON[Response[models.Event]](r.c, ctx, http.MethodPatch, path, q, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

func (r *EventsResource) Delete(ctx context.Context, grantID, eventID string, params DestroyEventParams) error {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events/" + url.PathEscape(eventID)
	_, _, err := DoJSON[map[string]any](r.c, ctx, http.MethodDelete, path, q, nil, nil)
	return err
}

func (r *EventsResource) SendRSVP(ctx context.Context, grantID, eventID string, req SendRSVPRequest, params SendRSVPParams) (*RequestIDResponse, error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/events/" + url.PathEscape(eventID) + "/send-rsvp"
	out, headers, err := DoJSON[RequestIDResponse](r.c, ctx, http.MethodPost, path, q, req, nil)
	if err != nil {
		return nil, err
	}
	// prefer header request id if present
	if rid := headers.Get("x-request-id"); rid != "" {
		out.RequestID = rid
	}
	return out, nil
}
