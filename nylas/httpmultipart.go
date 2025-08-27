package nylas

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
)

// doMultipart sends a multipart/form-data request with fields and one or more files.
func (c *Client) doMultipart(ctx context.Context, method, path string, query url.Values, fields map[string]string, files map[string]io.Reader) (*http.Response, error) {
	u, err := c.buildURL(path, query)
	if err != nil {
		return nil, err
	}

	// Build multipart body once, reuse bytes for retries
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for k, v := range fields {
		_ = mw.WriteField(k, v)
	}
	for name, r := range files {
		fw, err := mw.CreateFormFile(name, name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fw, r); err != nil {
			return nil, err
		}
	}
	_ = mw.Close()
	contentType := mw.FormDataContentType()

	policy := DefaultRetryPolicy()
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		rdr := bytes.NewReader(body.Bytes())
		req, err := http.NewRequest(method, u, rdr)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)

		resp, err := c.do(ctx, req)
		if err != nil {
			err = wrapTransportError(err, req, c.http.Timeout)
			if attempt < policy.MaxRetries {
				c.backoffSleep(attempt, http.Header{}, policy)
				continue
			}
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if c.ShouldRetry(resp.StatusCode) && attempt < policy.MaxRetries {
				c.backoffSleep(attempt, resp.Header, policy)
				// close body before retry to avoid leaks
				resp.Body.Close()
				continue
			}
			// parseAPIError will close the body
			return resp, parseAPIError(resp)
		}

		// 2xx: caller owns closing (e.g., for streaming download of response content)
		return resp, nil
	}
	return nil, fmt.Errorf("multipart: exhausted retries")
}

// doMultipartParts sends multipart/form-data with custom per-part headers.
// This is analogous to doMultipart, but lets callers set Content-Type for each part.
func (c *Client) doMultipartParts(
	ctx context.Context,
	method, path string,
	query url.Values,
	fields []FormField,
	files []FormFile,
) (*http.Response, error) {
	u, err := c.buildURL(path, query)
	if err != nil {
		return nil, err
	}

	// Build multipart once; we’ll reuse the bytes on retries.
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	// Fields
	for _, f := range fields {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="`+f.Name+`"`)
		if f.ContentType != "" {
			h.Set("Content-Type", f.ContentType)
		}
		pw, err := mw.CreatePart(h)
		if err != nil {
			return nil, err
		}
		if _, err := pw.Write([]byte(f.Value)); err != nil {
			return nil, err
		}
	}

	// Files
	for _, ff := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="`+ff.Field+`"; filename="`+ff.Filename+`"`)
		if ff.ContentType != "" {
			h.Set("Content-Type", ff.ContentType)
		}
		pw, err := mw.CreatePart(h)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(pw, ff.Reader); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}
	contentType := mw.FormDataContentType()

	policy := DefaultRetryPolicy()
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		rdr := bytes.NewReader(body.Bytes())
		req, err := http.NewRequest(method, u, rdr)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)

		resp, err := c.do(ctx, req)
		if err != nil {
			err = wrapTransportError(err, req, c.http.Timeout)
			if attempt < policy.MaxRetries {
				c.backoffSleep(attempt, http.Header{}, policy)
				continue
			}
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if c.ShouldRetry(resp.StatusCode) && attempt < policy.MaxRetries {
				c.backoffSleep(attempt, resp.Header, policy)
				resp.Body.Close()
				continue
			}
			// parseAPIError closes the body
			return resp, parseAPIError(resp)
		}

		// 2xx: caller is responsible for closing.
		return resp, nil
	}
	return nil, fmt.Errorf("multipart: exhausted retries")
}

// FileReader opens a path as a reader for multipart uploads.
func FileReader(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
