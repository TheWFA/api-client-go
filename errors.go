package wfa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// maxBodySnippetLength bounds how much of an unparseable error body is
// included in an APIError's message.
const maxBodySnippetLength = 500

// APIError is returned for any non-2xx HTTP response from the API.
type APIError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Code is the API's machine-readable error code, when available.
	Code string
	// Message is a human-readable description of the error.
	Message string
	// RawBody is the raw response body, truncated to a few hundred bytes.
	RawBody string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("wfa: %s (status %d)", e.Message, e.StatusCode)
}

// errorBody mirrors the API's own error envelope: {"error": {"code", "message"}}.
type errorBody struct {
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	// Message is the generic shape returned by the API gateway itself when a
	// request is rejected before reaching the API, e.g. a missing or invalid
	// API key: {"message": "Forbidden"}.
	Message string `json:"message"`
}

func newAPIError(statusCode int, body []byte) *APIError {
	snippet := string(body)
	if len(snippet) > maxBodySnippetLength {
		snippet = snippet[:maxBodySnippetLength]
	}

	var parsed errorBody
	if err := json.Unmarshal(body, &parsed); err == nil {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return &APIError{
				StatusCode: statusCode,
				Code:       parsed.Error.Code,
				Message:    parsed.Error.Message,
				RawBody:    snippet,
			}
		}

		if parsed.Message != "" {
			return &APIError{
				StatusCode: statusCode,
				Message:    parsed.Message,
				RawBody:    snippet,
			}
		}
	}

	message := fmt.Sprintf("request failed with status %d", statusCode)
	if snippet != "" {
		message = fmt.Sprintf("%s: %s", message, snippet)
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		RawBody:    snippet,
	}
}

func statusCodeIs(err error, statusCode int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == statusCode
}

// IsBadRequest reports whether err is an *APIError with a 400 status code.
func IsBadRequest(err error) bool { return statusCodeIs(err, http.StatusBadRequest) }

// IsUnauthorized reports whether err is an *APIError with a 401 status code.
func IsUnauthorized(err error) bool { return statusCodeIs(err, http.StatusUnauthorized) }

// IsForbidden reports whether err is an *APIError with a 403 status code.
func IsForbidden(err error) bool { return statusCodeIs(err, http.StatusForbidden) }

// IsNotFound reports whether err is an *APIError with a 404 status code.
func IsNotFound(err error) bool { return statusCodeIs(err, http.StatusNotFound) }

// IsRateLimited reports whether err is an *APIError with a 429 status code.
func IsRateLimited(err error) bool { return statusCodeIs(err, http.StatusTooManyRequests) }
