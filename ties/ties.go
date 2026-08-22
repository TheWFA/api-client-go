// Package ties gives access to the WFA Matchday API's /ties endpoints.
//
// A match's MatchGroup points here when the match belongs to a two-legged
// tie.
package ties

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Leg is a single leg of a two-legged tie.
type Leg struct {
	MatchID      int         `json:"matchId"`
	LegNumber    *int        `json:"legNumber"`
	ScheduledFor *wfa.Time   `json:"scheduledFor"`
	Status       string      `json:"status"`
	HomeTeam     wfa.TeamRef `json:"homeTeam"`
	AwayTeam     wfa.TeamRef `json:"awayTeam"`
	HomeScore    int         `json:"homeScore"`
	AwayScore    int         `json:"awayScore"`
}

// Aggregate is the aggregate score across a tie's legs.
type Aggregate struct {
	HomeScore int  `json:"homeScore"`
	AwayScore int  `json:"awayScore"`
	Complete  bool `json:"complete"`
}

// Tie is a two-legged tie, with its legs and aggregate score.
type Tie struct {
	ID          int                    `json:"id"`
	Competition wfa.CompetitionMiniRef `json:"competition"`
	Season      wfa.SeasonMiniRef      `json:"season"`
	HomeTeam    wfa.TeamRef            `json:"homeTeam"`
	AwayTeam    wfa.TeamRef            `json:"awayTeam"`
	CreatedAt   wfa.Time               `json:"createdAt"`
	Legs        []Leg                  `json:"legs"`
	Aggregate   Aggregate              `json:"aggregate"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	CompetitionID *int
	SeasonID      *int
	TeamID        *int
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)
	wfa.SetInt(v, "teamId", q.TeamID)

	return v
}

// Service gives access to the /ties endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of two-legged ties, with their legs and
// aggregate score.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Tie], error) {
	return wfa.Get[wfa.ListResponse[Tie]](ctx, s.backend, "/ties", query.encode())
}

// Get retrieves a single tie by ID.
func (s *Service) Get(ctx context.Context, id int) (Tie, error) {
	return wfa.Get[Tie](ctx, s.backend, fmt.Sprintf("/ties/%d", id), nil)
}
