package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type ContactsResource struct{ c *Client }

// List -> GET /v3/grants/{identifier}/contacts
func (r *ContactsResource) List(
	ctx context.Context,
	identifier string,
	q *models.ListContactsQueryParams,
) (*ListResponse[models.Contact], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts"
	out, headers, err := DoJSON[ListResponse[models.Contact]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get -> GET /v3/grants/{identifier}/contacts/{contact_id}
func (r *ContactsResource) Get(
	ctx context.Context,
	identifier, contactID string,
	q *models.FindContactQueryParams,
) (*Response[models.Contact], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts/" + url.PathEscape(contactID)
	out, headers, err := DoJSON[Response[models.Contact]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/grants/{identifier}/contacts
func (r *ContactsResource) Create(
	ctx context.Context,
	identifier string,
	body models.CreateContactRequest,
) (*Response[models.Contact], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts"
	out, headers, err := DoJSON[Response[models.Contact]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/contacts/{contact_id}
func (r *ContactsResource) Update(
	ctx context.Context,
	identifier, contactID string,
	body models.UpdateContactRequest,
) (*Response[models.Contact], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts/" + url.PathEscape(contactID)
	out, headers, err := DoJSON[Response[models.Contact]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Delete -> DELETE /v3/grants/{identifier}/contacts/{contact_id}
func (r *ContactsResource) Delete(
	ctx context.Context,
	identifier, contactID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts/" + url.PathEscape(contactID)
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

func (r *ContactsResource) ListGroups(
	ctx context.Context,
	identifier string,
	q *models.ListContactGroupsQueryParams,
) (*ListResponse[models.ContactGroup], error) {
	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/contacts/groups"
	out, headers, err := DoJSON[ListResponse[models.ContactGroup]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}
