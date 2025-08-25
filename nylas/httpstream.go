package nylas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) doStream(ctx context.Context, method, path string, query url.Values, headers http.Header) (*http.Response, error) {
	u, err := c.buildURL(path, query)
	if err != nil {
		return nil, err
	}
	policy := DefaultRetryPolicy()
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		req, err := http.NewRequest(method, u, nil)
		if err != nil {
			return nil, err
		}
		for k, vs := range headers {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
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
			// parseAPIError will close body
			return nil, parseAPIError(resp)
		}
		// 2xx: return open stream; caller must Close()
		return resp, nil
	}
	return nil, fmt.Errorf("stream: exhausted retries")
}
