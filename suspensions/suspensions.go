// Package suspensions gives access to the WFA Matchday API's /suspensions
// endpoints.
package suspensions

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Status is the current status of a suspension.
type Status string

const (
	StatusActive     Status = "active"
	StatusServed     Status = "served"
	StatusAppealed   Status = "appealed"
	StatusOverturned Status = "overturned"
)

// OffenceRef is a lightweight reference to the offence that originated a
// suspension.
type OffenceRef struct {
	ID               wfa.Snowflake `json:"id"`
	Name             string        `json:"name"`
	SuspensionLength int           `json:"suspensionLength"`
}

// Origin describes the discipline record a suspension arose from.
type Origin struct {
	DisciplineID wfa.Snowflake `json:"disciplineId"`
	Card         string        `json:"card"`
	Offence      *OffenceRef   `json:"offence"`
}

// ServedInMatch is a lightweight match reference, as embedded in
// Suspension.ServedIn.
type ServedInMatch struct {
	ID           wfa.Snowflake   `json:"id"`
	ScheduledFor *wfa.Time       `json:"scheduledFor"`
	Status       string          `json:"status"`
	HomeTeam     wfa.TeamMiniRef `json:"homeTeam"`
	AwayTeam     wfa.TeamMiniRef `json:"awayTeam"`
}

// Suspension is a single suspension.
type Suspension struct {
	ID               wfa.Snowflake          `json:"id"`
	Person           wfa.PersonRef          `json:"person"`
	Competition      wfa.CompetitionMiniRef `json:"competition"`
	Season           wfa.SeasonMiniRef      `json:"season"`
	Origin           *Origin                `json:"origin"`
	MatchesTotal     int                    `json:"matchesTotal"`
	MatchesServed    int                    `json:"matchesServed"`
	MatchesRemaining int                    `json:"matchesRemaining"`
	Status           Status                 `json:"status"`
	CreatedAt        wfa.Time               `json:"createdAt"`
	ServedIn         []ServedInMatch        `json:"servedIn"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	PersonID        *wfa.Snowflake
	CompetitionID   *wfa.Snowflake
	SeasonID        *wfa.Snowflake
	ServedInMatchID *wfa.Snowflake
	OriginMatchID   *wfa.Snowflake
	Status          []Status
	ActiveOnly      *bool
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetSnowflake(v, "personId", q.PersonID)
	wfa.SetSnowflake(v, "competitionId", q.CompetitionID)
	wfa.SetSnowflake(v, "seasonId", q.SeasonID)
	wfa.SetSnowflake(v, "servedInMatchId", q.ServedInMatchID)
	wfa.SetSnowflake(v, "originMatchId", q.OriginMatchID)
	wfa.SetEnums(v, "status", q.Status)
	wfa.SetBool(v, "activeOnly", q.ActiveOnly)

	return v
}

// PersonQuery holds the filters accepted by persons.Service.Suspensions,
// i.e. ListQuery without PersonID.
type PersonQuery struct {
	wfa.ListParams
	CompetitionID   *wfa.Snowflake
	SeasonID        *wfa.Snowflake
	ServedInMatchID *wfa.Snowflake
	OriginMatchID   *wfa.Snowflake
	Status          []Status
	ActiveOnly      *bool
}

// Encode renders q as URL query parameters.
func (q PersonQuery) Encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetSnowflake(v, "competitionId", q.CompetitionID)
	wfa.SetSnowflake(v, "seasonId", q.SeasonID)
	wfa.SetSnowflake(v, "servedInMatchId", q.ServedInMatchID)
	wfa.SetSnowflake(v, "originMatchId", q.OriginMatchID)
	wfa.SetEnums(v, "status", q.Status)
	wfa.SetBool(v, "activeOnly", q.ActiveOnly)

	return v
}

// Service gives access to the /suspensions endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of suspensions.
//
// Filtering by fixture is split in two: ServedInMatchID returns bans sat out
// in that fixture, OriginMatchID returns bans that arose from a card shown in
// it.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Suspension], error) {
	return wfa.Get[wfa.ListResponse[Suspension]](ctx, s.backend, "/suspensions", query.encode())
}

// Get retrieves a single suspension by ID.
func (s *Service) Get(ctx context.Context, id wfa.Snowflake) (Suspension, error) {
	return wfa.Get[Suspension](ctx, s.backend, fmt.Sprintf("/suspensions/%d", id), nil)
}
