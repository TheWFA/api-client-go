// Package organisations gives access to the WFA Matchday API's
// /organisations endpoints.
package organisations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/competitions"
)

// Organisation is a summary of an organisation.
type Organisation struct {
	ID        wfa.Snowflake `json:"id"`
	Name      string        `json:"name"`
	ShortName *string       `json:"shortName"`
	BadgeURL  *string       `json:"badgeUrl"`
	SortOrder int           `json:"sortOrder"`
	CreatedAt wfa.Time      `json:"createdAt"`
}

// FullOrganisation is an Organisation with its competitions and, when
// requested, its history.
type FullOrganisation struct {
	ID           wfa.Snowflake              `json:"id"`
	Name         string                     `json:"name"`
	ShortName    *string                    `json:"shortName"`
	BadgeURL     *string                    `json:"badgeUrl"`
	SortOrder    int                        `json:"sortOrder"`
	CreatedAt    wfa.Time                   `json:"createdAt"`
	Competitions []competitions.Competition `json:"competitions"`
	History      []wfa.HistoryEntry         `json:"history,omitempty"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)

	return v
}

// Service gives access to the /organisations endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of organisations.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Organisation], error) {
	return wfa.Get[wfa.ListResponse[Organisation]](ctx, s.backend, "/organisations", query.encode())
}

// Get retrieves detailed information about a specific organisation,
// including its competitions.
func (s *Service) Get(ctx context.Context, id wfa.Snowflake) (FullOrganisation, error) {
	return wfa.Get[FullOrganisation](ctx, s.backend, fmt.Sprintf("/organisations/%d", id), nil)
}
