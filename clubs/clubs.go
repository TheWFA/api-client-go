// Package clubs gives access to the WFA Matchday API's /clubs endpoints.
package clubs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Club is a summary of a club.
type Club struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ClubLogo     *string  `json:"clubLogo"`
	ContactEmail *string  `json:"contactEmail"`
	CreatedAt    wfa.Time `json:"createdAt"`
}

// TeamRef is a lightweight reference to a team, as embedded in FullClub.
type TeamRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbreviated string `json:"abbreviated"`
	Nickname    string `json:"nickname"`
	BadgeURL    string `json:"badgeUrl"`
	Primary     string `json:"primary"`
	Secondary   string `json:"secondary"`
}

// FullClub is a Club with its teams and, when requested, its history.
type FullClub struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	ClubLogo     *string            `json:"clubLogo"`
	ContactEmail *string            `json:"contactEmail"`
	CreatedAt    wfa.Time           `json:"createdAt"`
	Teams        []TeamRef          `json:"teams"`
	History      []wfa.HistoryEntry `json:"history,omitempty"`
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

// Service gives access to the /clubs endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of clubs.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Club], error) {
	return wfa.Get[wfa.ListResponse[Club]](ctx, s.backend, "/clubs", query.encode())
}

// Get retrieves detailed information about a specific club, including its
// teams.
func (s *Service) Get(ctx context.Context, id int) (FullClub, error) {
	return wfa.Get[FullClub](ctx, s.backend, fmt.Sprintf("/clubs/%d", id), nil)
}
