package media

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWaitUntilAvailableImmediateSuccess(t *testing.T) {
	requests := 0
	method := ""
	rangeHeader := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		method = r.Method
		rangeHeader = r.Header.Get("Range")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	err := waitUntilAvailable(context.Background(), server.URL, 3, time.Millisecond, time.Millisecond)
	if err != nil {
		t.Fatalf("waitUntilAvailable returned error: %v", err)
	}
	if requests != 1 {
		t.Fatalf("expected one request, got %d", requests)
	}
	if method != http.MethodGet || rangeHeader != "bytes=0-0" {
		t.Fatalf("expected ranged GET probe, got method=%q range=%q", method, rangeHeader)
	}
}

func TestWaitUntilAvailableRetriesTransientNotFound(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		if requests < 3 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	err := waitUntilAvailable(context.Background(), server.URL, 4, time.Millisecond, time.Millisecond)
	if err != nil {
		t.Fatalf("waitUntilAvailable returned error: %v", err)
	}
	if requests != 3 {
		t.Fatalf("expected three requests, got %d", requests)
	}
}

func TestWaitUntilAvailableStopsAfterBoundedRetries(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	err := waitUntilAvailable(context.Background(), server.URL, 3, time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected availability error")
	}
	if !strings.Contains(err.Error(), "after 3 attempts") {
		t.Fatalf("expected bounded retry error, got %v", err)
	}
	if requests != 3 {
		t.Fatalf("expected three requests, got %d", requests)
	}
}

func TestWaitUntilAvailableDoesNotRetryPermanentClientError(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)

	err := waitUntilAvailable(context.Background(), server.URL, 3, time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected availability error")
	}
	if requests != 1 {
		t.Fatalf("expected no retry for permanent client error, got %d requests", requests)
	}
}

func TestWaitUntilAvailableHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitUntilAvailable(ctx, "https://example.invalid/file.pdf", 3, time.Millisecond, time.Millisecond)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}
