package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type FoldersResource struct{ c *Client }

// --- Query params ---

type ListFoldersParams struct {
	// Generic list params
	Limit     *int    `query:"limit"`
	PageToken *string `query:"page_token"`
	Select    *string `query:"select"`

	// Folder-specific (provider-dependent)
	ParentID             *string `query:"parent_id"`              // MS/EWS
	IncludeHiddenFolders *bool   `query:"include_hidden_folders"` // MS
	SingleLevel          *bool   `query:"single_level"`           // MS
}

type FindFolderParams struct {
	Select *string `query:"select"`
}

// --- Requests ---

type CreateFolderRequest struct {
	Name            string  `json:"name"`
	ParentID        *string `json:"parent_id,omitempty"`        // MS only
	BackgroundColor *string `json:"background_color,omitempty"` // Google only
	TextColor       *string `json:"text_color,omitempty"`       // Google only
}

type UpdateFolderRequest struct {
	Name            *string `json:"name,omitempty"`
	ParentID        *string `json:"parent_id,omitempty"`        // MS only
	BackgroundColor *string `json:"background_color,omitempty"` // Google only
	TextColor       *string `json:"text_color,omitempty"`       // Google only
}

// --- Endpoints ---

// List returns folders for a grant.
func (r *FoldersResource) List(ctx context.Context, grantID string, params *ListFoldersParams) (*ListResponse[models.Folder], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/folders"
	out, headers, err := DoJSON[ListResponse[models.Folder]](r.c, ctx, http.MethodGet, path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Get fetches a single folder by ID.
func (r *FoldersResource) Get(ctx context.Context, grantID, folderID string, params *FindFolderParams) (*Response[models.Folder], error) {
	q := EncodeQuery(params)
	path := "/v3/grants/" + url.PathEscape(grantID) + "/folders/" + url.PathEscape(folderID)
	out, headers, err := DoJSON[Response[models.Folder]](r.c, ctx, http.MethodGet, path, q, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Create makes a new folder.
func (r *FoldersResource) Create(ctx context.Context, grantID string, req CreateFolderRequest) (*Response[models.Folder], error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/folders"
	out, headers, err := DoJSON[Response[models.Folder]](r.c, ctx, http.MethodPost, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Update modifies an existing folder.
func (r *FoldersResource) Update(ctx context.Context, grantID, folderID string, req UpdateFolderRequest) (*Response[models.Folder], error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/folders/" + url.PathEscape(folderID)
	out, headers, err := DoJSON[Response[models.Folder]](r.c, ctx, http.MethodPatch, path, nil, req, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	return out, nil
}

// Delete removes a folder.
func (r *FoldersResource) Delete(ctx context.Context, grantID, folderID string) error {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/folders/" + url.PathEscape(folderID)
	_, _, err := DoJSON[map[string]any](r.c, ctx, http.MethodDelete, path, nil, nil, nil)
	return err
}
