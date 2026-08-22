package wfa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DefaultBaseURL is the default base URL used by NewBackend, without the API
// version path segment.
const DefaultBaseURL = "https://api.thewfa.org.uk"

// apiVersionPath is appended to the configured base URL to reach the v1 API.
const apiVersionPath = "/v1"

// Backend is the shared HTTP transport used by every resource package's
// Service. Construct one with NewBackend and pass it to each resource
// package's New function — or use the github.com/TheWFA/api-client-go/client
// package, which does this for every resource at once.
//
// A Backend is safe for concurrent use by multiple goroutines.
type Backend struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	headers    map[string]string
	userAgent  string
}

// Option configures a Backend. See WithBaseURL, WithHTTPClient, WithHeader
// and WithUserAgent.
type Option func(*Backend)

// WithBaseURL overrides the API base URL, including any version path segment.
// It's primarily useful for pointing the client at a mock server in tests.
func WithBaseURL(baseURL string) Option {
	return func(b *Backend) { b.baseURL = baseURL }
}

// WithHTTPClient overrides the *http.Client used to make requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(b *Backend) { b.httpClient = httpClient }
}

// WithHeader sets an additional header sent on every request.
func WithHeader(key, value string) Option {
	return func(b *Backend) { b.headers[key] = value }
}

// WithUserAgent overrides the User-Agent header sent on every request.
func WithUserAgent(userAgent string) Option {
	return func(b *Backend) { b.userAgent = userAgent }
}

// NewBackend constructs a Backend authenticated with apiKey, sent on every
// request as the x-api-key header.
func NewBackend(apiKey string, opts ...Option) (*Backend, error) {
	if apiKey == "" {
		return nil, errors.New("wfa: apiKey is required")
	}

	b := &Backend{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL + apiVersionPath,
		httpClient: http.DefaultClient,
		headers:    map[string]string{},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b, nil
}

func (b *Backend) newRequest(ctx context.Context, method, path string, query url.Values) (*http.Request, error) {
	full := b.baseURL + path

	if len(query) > 0 {
		full += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, full, nil)
	if err != nil {
		return nil, fmt.Errorf("wfa: building request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", b.apiKey)

	if b.userAgent != "" {
		req.Header.Set("User-Agent", b.userAgent)
	}

	for k, v := range b.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (b *Backend) do(req *http.Request, out any) error {
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("wfa: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("wfa: reading response body: %w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, body)
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("wfa: decoding response: %w", err)
	}

	return nil
}

// Get performs a GET request against path with the given query parameters
// and decodes a successful JSON response into out. If out is nil, the
// response body is discarded. Non-2xx responses are converted into an
// *APIError.
//
// This is the low-level primitive resource packages build on; most callers
// should use the generic Get function instead.
func (b *Backend) Get(ctx context.Context, path string, query url.Values, out any) error {
	req, err := b.newRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return err
	}

	return b.do(req, out)
}

// Get performs a GET request against path with the given query parameters
// and decodes the JSON response body into a value of type T. It's the
// primitive every resource package's Service methods are built on.
func Get[T any](ctx context.Context, b *Backend, path string, query url.Values) (T, error) {
	var out T

	if err := b.Get(ctx, path, query, &out); err != nil {
		return out, err
	}

	return out, nil
}
