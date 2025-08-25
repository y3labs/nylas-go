package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) buildURL(path string, q url.Values) (string, error) {
	base := c.serverURL
	if base == "" {
		return "", fmt.Errorf("serverURL not set")
	}
	if len(q) > 0 {
		return fmt.Sprintf("%s%s?%s", base, path, q.Encode()), nil
	}
	return base + path, nil
}

func DoJSON[T any](
	c *Client,
	ctx context.Context,
	method, path string,
	query url.Values,
	body any,
	headers http.Header,
) (*T, http.Header, error) {
	if headers == nil {
		headers = http.Header{}
	}
	var rdr *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		rdr = bytes.NewReader(b)
		headers.Set("Content-Type", "application/json")
	} else {
		rdr = bytes.NewReader(nil)
	}

	u, err := c.buildURL(path, query)
	if err != nil {
		return nil, nil, err
	}

	policy := DefaultRetryPolicy()
	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		// rebuild request each attempt to reset body
		var bodyReader io.Reader
		if rdr != nil {
			rdr.Seek(0, io.SeekStart)
			bodyReader = rdr
		}
		req, err := http.NewRequest(method, u, bodyReader)
		if err != nil {
			return nil, nil, err
		}
		for k, vs := range headers {
			for _, v := range vs {
				req.Header.Set(k, v)
			}
		}

		resp, err := c.do(ctx, req)
		if err != nil {
			// Promote timeouts to SDKTimeoutError
			err = wrapTransportError(err, req, c.http.Timeout)
			lastErr = err
			if attempt < policy.MaxRetries {
				c.backoffSleep(attempt, http.Header{}, policy)
				continue
			}
			return nil, nil, err
		}

		// Handle non-2xx first; parseAPIError will close the body
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if c.ShouldRetry(resp.StatusCode) && attempt < policy.MaxRetries {
				c.backoffSleep(attempt, resp.Header, policy)
				resp.Body.Close()
				continue
			}
			return nil, nil, parseAPIError(resp)
		}

		// Success: decode body and close afterwards
		defer resp.Body.Close()

		var out T
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&out); err != nil {
			if errors.Is(err, io.EOF) {
				return &out, resp.Header, nil
			}
			return nil, nil, err
		}
		return &out, resp.Header, nil
	}
	return nil, nil, lastErr
}
