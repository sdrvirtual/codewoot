package chatwoot

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
	baseURL    *url.URL
	token      string
	accountID  int
	httpClient *http.Client

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

func WithLogger(logf func(format string, args ...any)) Option {
	return func(c *Client) {
		c.logf = logf
	}
}

func New(baseURL, token string, accountID int, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = strings.TrimSuffix(u.Path, "/")

	c := &Client{
		baseURL:    u,
		token:      token,
		accountID:  accountID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
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
	if c.token != "" {
		req.Header.Set("api_access_token", c.token)
	}
	return req, nil
}

func (c *Client) do(req *http.Request) (json.RawMessage, error) {
	if c.logf != nil {
		c.logf("%s %s", req.Method, req.URL.String())
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("chatwoot API error: status=%d body=%s", res.StatusCode, string(b))
	}
	return json.RawMessage(b), nil
}
