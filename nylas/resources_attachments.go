package nylas

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/y3labs/nylas-go/nylas/models"
)

type AttachmentsResource struct{ c *Client }

// Upload a single attachment via multipart. Returns server-side attachment id and request_id.
/*
func (r *AttachmentsResource) Upload(ctx context.Context, grantID, filename string, reader io.Reader) (*models.UploadAttachmentResponse, error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/attachments"
	resp, err := r.c.doMultipart(ctx, http.MethodPost, path, nil, nil, map[string]io.Reader{filename: reader})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}
	var out models.UploadAttachmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	// Nylas also returns x-request-id header; keep the JSON value but prefer header if JSON omits it.
	if out.RequestID == "" {
		out.RequestID = resp.Header.Get("x-request-id")
	}
	return &out, nil
}
*/

// Get returns attachment metadata (Python: find), requires message_id as a query param.
func (r *AttachmentsResource) Get(
	ctx context.Context,
	grantID, attachmentID string,
	q *models.FindAttachmentQueryParams, // required in Python; pass &models.FindAttachmentQueryParams{MessageID: "..."}
) (*Response[models.Attachment], error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/attachments/" + url.PathEscape(attachmentID)

	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}

	out, headers, err := DoJSON[Response[models.Attachment]](r.c, ctx, http.MethodGet, path, query, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Headers = headers
	if out.RequestID == "" {
		out.RequestID = headers.Get("x-request-id")
	}
	return out, nil
}

// Download returns the raw *http.Response for streaming file bytes.
// Caller MUST Close() the body when finished.
func (r *AttachmentsResource) Download(
	ctx context.Context,
	grantID, attachmentID string,
	q *models.FindAttachmentQueryParams, // ?message_id=...
) (*http.Response, error) {
	path := "/v3/grants/" + url.PathEscape(grantID) + "/attachments/" + url.PathEscape(attachmentID) + "/download"

	var query url.Values
	if q != nil {
		query = EncodeQuery(*q)
	}

	return r.c.doStream(ctx, http.MethodGet, path, query, nil)
}

// DownloadBytes reads the entire attachment body into memory and returns it.
func (r *AttachmentsResource) DownloadBytes(
	ctx context.Context,
	grantID, attachmentID string,
	q *models.FindAttachmentQueryParams,
) ([]byte, error) {
	resp, err := r.Download(ctx, grantID, attachmentID, q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}
	return io.ReadAll(resp.Body)
}
