package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type ConfigurationsResource struct{ c *Client }

// Optional list params: limit/page_token/select (parity with ListQueryParams)
type ListConfigurationsParams struct {
	Limit     *int    `url:"limit,omitempty"`
	PageToken *string `url:"page_token,omitempty"`
	Select    *string `url:"select,omitempty"`
}

// List -> GET /v3/grants/{identifier}/scheduling/configurations
func (r *ConfigurationsResource) List(
	ctx context.Context,
	identifier string,
	params *ListConfigurationsParams,
) (*ListResponse[models.Configuration], error) {
	var query url.Values
	if params != nil {
		query = EncodeQuery(*params)
	}
	path := "/v3/grants/" + url.PathEscape(identifier) + "/scheduling/configurations"
	out, headers, err := DoJSON[ListResponse[models.Configuration]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Find -> GET /v3/grants/{identifier}/scheduling/configurations/{config_id}
func (r *ConfigurationsResource) Find(
	ctx context.Context,
	identifier, configID string,
) (*Response[models.Configuration], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/scheduling/configurations/" + url.PathEscape(configID)
	out, headers, err := DoJSON[Response[models.Configuration]](r.c, ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Create -> POST /v3/grants/{identifier}/scheduling/configurations
func (r *ConfigurationsResource) Create(
	ctx context.Context,
	identifier string,
	body models.CreateConfigurationRequest,
) (*Response[models.Configuration], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/scheduling/configurations"
	out, headers, err := DoJSON[Response[models.Configuration]](r.c, ctx, http.MethodPost, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Update -> PUT /v3/grants/{identifier}/scheduling/configurations/{config_id}
func (r *ConfigurationsResource) Update(
	ctx context.Context,
	identifier, configID string,
	body models.UpdateConfigurationRequest,
) (*Response[models.Configuration], error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/scheduling/configurations/" + url.PathEscape(configID)
	out, headers, err := DoJSON[Response[models.Configuration]](r.c, ctx, http.MethodPut, path, nil, body, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Destroy -> DELETE /v3/grants/{identifier}/scheduling/configurations/{config_id}
func (r *ConfigurationsResource) Destroy(
	ctx context.Context,
	identifier, configID string,
) (*DeleteResponse, error) {
	path := "/v3/grants/" + url.PathEscape(identifier) + "/scheduling/configurations/" + url.PathEscape(configID)
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
