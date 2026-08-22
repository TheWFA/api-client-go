// Package seasons gives access to the WFA Matchday API's /seasons
// endpoints.
package seasons

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Season is a summary of a season.
type Season = wfa.SeasonRef

// FullSeason is a SeasonRef with its active flag and team/competition counts.
type FullSeason struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	StartDate           wfa.Time `json:"startDate"`
	EndDate             wfa.Time `json:"endDate"`
	Active              bool     `json:"active"`
	CompetitionCount    int      `json:"competitionCount"`
	RegisteredTeamCount int      `json:"registeredTeamCount"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	// ActiveOn is an ISO date string; restricts to the season active on this
	// date.
	ActiveOn string
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	wfa.SetString(v, "activeOn", q.ActiveOn)

	return v
}

// Service gives access to the /seasons endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of seasons, newest first.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Season], error) {
	return wfa.Get[wfa.ListResponse[Season]](ctx, s.backend, "/seasons", query.encode())
}

// Get retrieves a single season by its unique identifier.
func (s *Service) Get(ctx context.Context, id int) (FullSeason, error) {
	return wfa.Get[FullSeason](ctx, s.backend, fmt.Sprintf("/seasons/%d", id), nil)
}
