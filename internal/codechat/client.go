package codechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Client struct {
	baseURL       *url.URL
	globalToken   string
	instanceToken string
	instance      string
	httpClient    *http.Client

	logf func(format string, args ...any)
}

type Option func(*Client)

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

func WithInstanceToken(instanceToken string, instance string) Option {
	return func(c *Client) {
		c.instanceToken = "Bearer " + instanceToken
		c.instance = instance
	}
}

func WithLogger(logf func(format string, args ...any)) Option {
	return func(c *Client) {
		c.logf = logf
	}
}

func New(base, token string, opts ...Option) (*Client, error) {
	if base == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = strings.TrimSuffix(u.Path, "/")

	c := &Client{
		baseURL:     u,
		globalToken: token,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("codechat API error: status=%d body=%s", e.StatusCode, e.Body)
}

func (c *Client) newRequest(ctx context.Context, method, p string, body any) (*http.Request, error) {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, p)

	var rdr io.Reader
	setJSON := false
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			rdr = r
		} else {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(body); err != nil {
				return nil, fmt.Errorf("encode body: %w", err)
			}
			rdr = buf
			setJSON = true
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), rdr)
	if err != nil {
		return nil, err
	}
	if setJSON {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("apikey", c.globalToken)
	if c.instanceToken != "" {
		req.Header.Set("Authorization", c.instanceToken)
	}
	return req, nil
}

func (c *Client) do(req *http.Request) (json.RawMessage, *http.Response, error) {
	if c.logf != nil {
		c.logf("%s %s", req.Method, req.URL.String())
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, nil, &APIError{StatusCode: res.StatusCode, Body: string(b)}
	}

	ct := res.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") || strings.Contains(ct, "json") {
		b, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		return json.RawMessage(b), nil, nil
	}

	// Non-JSON: return the response for the caller to consume/stream and close.
	return nil, res, nil
}
