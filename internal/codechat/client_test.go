package codechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew_InvalidBase(t *testing.T) {
	_, err := New("://bad", "tok")
	if err == nil {
		t.Fatalf("expected error for invalid base URL, got nil")
	}
}

func TestNew_DefaultHTTPClientTimeout(t *testing.T) {
	c, err := New("http://example.com", "tok")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if c.httpClient == nil {
		t.Fatalf("expected default httpClient, got nil")
	}
	if c.httpClient.Timeout != 30*time.Second {
		t.Fatalf("expected timeout 30s, got %v", c.httpClient.Timeout)
	}
}

func TestWithHTTPClient_OverridesClient(t *testing.T) {
	custom := &http.Client{Timeout: 10 * time.Second}
	c, err := New("http://example.com", "tok", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if c.httpClient != custom {
		t.Fatalf("expected custom http client to be set")
	}
}

func TestWithInstanceToken_SetsAuthorizationAndInstance(t *testing.T) {
	c, err := New("http://example.com", "tok", WithInstanceToken("instTok", "instName"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if c.instanceToken != "Bearer instTok" {
		t.Fatalf("expected instanceToken 'Bearer instTok', got %q", c.instanceToken)
	}
	if c.instance != "instName" {
		t.Fatalf("expected instance 'instName', got %q", c.instance)
	}
}

func TestWithLogger_SetsLogger(t *testing.T) {
	var got string
	logger := func(format string, args ...any) {
		got = format // just confirm it was set; do() tests exercise logging
	}
	c, err := New("http://example.com", "tok", WithLogger(logger))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if c.logf == nil {
		t.Fatalf("expected logger to be set")
	}
	if got != "" {
		// not called yet here, just ensure assignment didn't log
		t.Fatalf("logger should not be called during New")
	}
}

func TestNewRequest_JSONBodyHeadersAndEncoding(t *testing.T) {
	c, err := New("http://example.com/", "gtoken", WithInstanceToken("itoken", "inst"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	body := map[string]string{"x": "<script>"}
	req, err := c.newRequest(context.Background(), http.MethodPost, "/foo", body)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}

	if req.URL.Path != "/foo" {
		t.Fatalf("expected path /foo, got %s", req.URL.Path)
	}

	if ct := req.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	if apikey := req.Header.Get("apikey"); apikey != "gtoken" {
		t.Fatalf("expected apikey header 'gtoken', got %q", apikey)
	}
	if auth := req.Header.Get("Authorization"); auth != "Bearer itoken" {
		t.Fatalf("expected Authorization 'Bearer itoken', got %q", auth)
	}

	// Ensure SetEscapeHTML(false): angle brackets should not be escaped.
	b, _ := io.ReadAll(req.Body)
	if !strings.Contains(string(b), "<script>") {
		t.Fatalf("expected JSON body to contain '<script>', got %s", string(b))
	}
}

func TestNewRequest_ReaderBody_NoJSONHeader(t *testing.T) {
	c, err := New("http://example.com", "tok")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	r := strings.NewReader("raw-data")
	req, err := c.newRequest(context.Background(), http.MethodPut, "/bar", r)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}
	if req.Header.Get("Content-Type") != "" {
		t.Fatalf("did not expect Content-Type to be set for io.Reader body")
	}
	b, _ := io.ReadAll(req.Body)
	if string(b) != "raw-data" {
		t.Fatalf("expected body 'raw-data', got %q", string(b))
	}
}

func TestDo_JSONResponse(t *testing.T) {
	var logged string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c, err := New(srv.URL, "tok", WithLogger(func(format string, args ...any) {
		logged = format // we only need to know it was called with something
	}))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	req, err := c.newRequest(context.Background(), http.MethodGet, "/json", nil)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}

	rm, res, err := c.do(req)
	if err != nil {
		t.Fatalf("do error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected response to be nil for JSON content")
	}
	if rm == nil {
		t.Fatalf("expected json body, got nil")
	}
	var v map[string]any
	if err := json.Unmarshal(rm, &v); err != nil || v["ok"] != true {
		t.Fatalf("unexpected JSON: %s err=%v", string(rm), err)
	}
	if logged == "" {
		t.Fatalf("expected do() to call logger")
	}
}

func TestDo_NonJSONResponse(t *testing.T) {
	var logged string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plain" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	c, err := New(srv.URL, "tok", WithLogger(func(format string, args ...any) {
		logged = format
	}))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	req, err := c.newRequest(context.Background(), http.MethodGet, "/plain", nil)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}

	rm, res, err := c.do(req)
	if err != nil {
		t.Fatalf("do error: %v", err)
	}
	if rm != nil {
		t.Fatalf("expected nil json body for non-JSON content")
	}
	if res == nil {
		t.Fatalf("expected non-nil response for non-JSON content")
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if string(b) != "hello" {
		t.Fatalf("unexpected body: %q", string(b))
	}
	if logged == "" {
		t.Fatalf("expected do() to call logger")
	}
}

func TestDo_APIError(t *testing.T) {
	var logged string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	c, err := New(srv.URL, "tok", WithLogger(func(format string, args ...any) {
		logged = format
	}))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	req, err := c.newRequest(context.Background(), http.MethodGet, "/err", nil)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}

	rm, res, err := c.do(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if rm != nil || res != nil {
		t.Fatalf("expected nil bodies on error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError || apiErr.Body != "boom" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
	if !strings.Contains(apiErr.Error(), "status=500") || !strings.Contains(apiErr.Error(), "boom") {
		t.Fatalf("unexpected APIError.Error(): %q", apiErr.Error())
	}
	if logged == "" {
		t.Fatalf("expected do() to call logger")
	}
}

func TestAPIError_ErrorString(t *testing.T) {
	e := &APIError{StatusCode: 400, Body: "bad"}
	got := e.Error()
	if !strings.Contains(got, "status=400") || !strings.Contains(got, "bad") {
		t.Fatalf("unexpected error string: %q", got)
	}
}

func TestNewRequest_PathJoinWithLeadingSlash(t *testing.T) {
	c, err := New("http://example.com/base/", "tok")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	req, err := c.newRequest(context.Background(), http.MethodGet, "/joined", nil)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}
	if req.URL.Path != "/base/joined" {
		t.Fatalf("expected path '/joined', got %q", req.URL.Path)
	}
}

func TestLoggerMessageFormatContainsMethodAndURL(t *testing.T) {
	var buf bytes.Buffer
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c, err := New(srv.URL, "tok", WithLogger(func(format string, args ...any) {
		buf.WriteString(fmt.Sprintf(format, args...))
	}))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	req, err := c.newRequest(context.Background(), http.MethodGet, "/x", nil)
	if err != nil {
		t.Fatalf("newRequest error: %v", err)
	}
	_, _, _ = c.do(req)

	log := buf.String()
	if !strings.Contains(log, "GET") || !strings.Contains(log, "/x") {
		t.Fatalf("expected log to contain method and URL, got %q", log)
	}
}
