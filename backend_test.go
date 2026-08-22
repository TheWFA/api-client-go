package wfa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestBackend(t *testing.T, handler http.HandlerFunc) (*Backend, *httptest.Server) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	backend, err := NewBackend("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}

	return backend, server
}

func TestNewBackendRequiresAPIKey(t *testing.T) {
	if _, err := NewBackend(""); err == nil {
		t.Fatal("expected an error for an empty API key")
	}
}

func TestBackendSendsAPIKeyHeader(t *testing.T) {
	var gotKey string

	backend, _ := newTestBackend(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	})

	if err := backend.Get(context.Background(), "/anything", nil, &struct{}{}); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotKey != "test-key" {
		t.Fatalf("x-api-key = %q, want %q", gotKey, "test-key")
	}
}

func TestBackendCustomHeaders(t *testing.T) {
	var gotUA, gotCustom string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		gotCustom = r.Header.Get("X-Custom")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	backend, err := NewBackend(
		"test-key",
		WithBaseURL(server.URL),
		WithUserAgent("wfa-test/1.0"),
		WithHeader("X-Custom", "abc"),
	)
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}

	if err := backend.Get(context.Background(), "/anything", nil, &struct{}{}); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotUA != "wfa-test/1.0" {
		t.Fatalf("User-Agent = %q, want %q", gotUA, "wfa-test/1.0")
	}

	if gotCustom != "abc" {
		t.Fatalf("X-Custom = %q, want %q", gotCustom, "abc")
	}
}

func TestBackendErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantCode   string
		wantMsg    string
		checkIs    func(error) bool
	}{
		{
			name:       "structured error body",
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"not_found","message":"Match not found"}}`,
			wantCode:   "not_found",
			wantMsg:    "Match not found",
			checkIs:    IsNotFound,
		},
		{
			name:       "gateway rejection shape",
			statusCode: http.StatusForbidden,
			body:       `{"message":"Forbidden"}`,
			wantMsg:    "Forbidden",
			checkIs:    IsForbidden,
		},
		{
			name:       "rate limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"error":{"code":"rate_limited","message":"Too many requests"}}`,
			wantCode:   "rate_limited",
			wantMsg:    "Too many requests",
			checkIs:    IsRateLimited,
		},
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       `{"error":{"code":"bad_request","message":"Invalid season"}}`,
			wantCode:   "bad_request",
			wantMsg:    "Invalid season",
			checkIs:    IsBadRequest,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"message":"Unauthorized"}`,
			wantMsg:    "Unauthorized",
			checkIs:    IsUnauthorized,
		},
		{
			name:       "unparseable body",
			statusCode: http.StatusInternalServerError,
			body:       `not json`,
			wantMsg:    "request failed with status 500: not json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend, _ := newTestBackend(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			})

			err := backend.Get(context.Background(), "/anything", nil, &struct{}{})
			if err == nil {
				t.Fatal("expected an error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("error is %T, want *APIError", err)
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.statusCode)
			}

			if apiErr.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", apiErr.Code, tt.wantCode)
			}

			if apiErr.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMsg)
			}

			if tt.checkIs != nil && !tt.checkIs(err) {
				t.Error("predicate check failed")
			}
		})
	}
}

func TestBackendNoContent(t *testing.T) {
	backend, _ := newTestBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var out struct{ Foo string }
	if err := backend.Get(context.Background(), "/anything", nil, &out); err != nil {
		t.Fatalf("Get: %v", err)
	}
}
