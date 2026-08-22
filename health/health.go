// Package health gives access to the WFA Matchday API's /health endpoint.
package health

import (
	"context"

	"github.com/TheWFA/api-client-go"
)

// Status is the overall status reported by the health endpoint.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// Scope indicates whether a health check result is safe to expose publicly
// or is only meant for internal consumers.
type Scope string

const (
	ScopePublic   Scope = "public"
	ScopeInternal Scope = "internal"
)

// Health is the API's health status.
type Health struct {
	Status    Status  `json:"status"`
	LatencyMs *int    `json:"latencyMs,omitempty"`
	Error     *string `json:"error,omitempty"`
	Scope     Scope   `json:"scope"`
}

// Service retrieves the API's health status.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// Get retrieves the API's health status.
func (s *Service) Get(ctx context.Context) (Health, error) {
	return wfa.Get[Health](ctx, s.backend, "/health", nil)
}
