// Package matches gives access to the WFA Matchday API's /matches endpoints.
package matches

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Status is the current state of a match.
type Status string

const (
	StatusScheduled           Status = "scheduled"
	StatusFirstHalf           Status = "first-half"
	StatusHalfTime            Status = "half-time"
	StatusSecondHalf          Status = "second-half"
	StatusExtraTimeFirstHalf  Status = "extra-time-first-half"
	StatusHalfTimeExtraTime   Status = "half-time-extra-time"
	StatusExtraTimeSecondHalf Status = "extra-time-second-half"
	StatusPenalties           Status = "penalties"
	StatusFullTime            Status = "full-time"
	StatusPostponed           Status = "postponed"
	StatusAbandoned           Status = "abandoned"
	StatusCancelled           Status = "cancelled"
	StatusForfeit             Status = "forfeit"
	StatusAwarded             Status = "awarded"
)

// PlayerPosition is the position a player took in a match lineup.
type PlayerPosition string

const (
	PlayerPositionLeft       PlayerPosition = "left"
	PlayerPositionRight      PlayerPosition = "right"
	PlayerPositionCentre     PlayerPosition = "centre"
	PlayerPositionGoalkeeper PlayerPosition = "goalkeeper"
	PlayerPositionSubstitute PlayerPosition = "sub"
)

// Times records when each period of a match started.
type Times struct {
	FirstHalfStartedAt           *wfa.Time `json:"firstHalfStartedAt"`
	SecondHalfStartedAt          *wfa.Time `json:"secondHalfStartedAt"`
	FirstHalfExtraTimeStartedAt  *wfa.Time `json:"firstHalfExtraTimeStartedAt"`
	SecondHalfExtraTimeStartedAt *wfa.Time `json:"secondHalfExtraTimeStartedAt"`
}

// Officials lists the officials assigned to a match.
type Officials struct {
	Referee        *wfa.PersonRef `json:"referee,omitempty"`
	Assistant1     *wfa.PersonRef `json:"assistant1,omitempty"`
	Assistant2     *wfa.PersonRef `json:"assistant2,omitempty"`
	FourthOfficial *wfa.PersonRef `json:"fourthOfficial,omitempty"`
}

// GroupRef is a lightweight reference to a competition match group (e.g. a
// game week, pool or knockout round).
type GroupRef struct {
	ID            wfa.Snowflake `json:"id"`
	CompetitionID wfa.Snowflake `json:"competitionId"`
	Name          string        `json:"name"`
}

// Court is the court a match was played on.
type Court struct {
	ID       wfa.Snowflake   `json:"id"`
	Name     string          `json:"name"`
	Location wfa.LocationRef `json:"location"`
}

// Stream is a broadcast stream for a match.
type Stream struct {
	ID           wfa.Snowflake `json:"id"`
	StreamURL    *string       `json:"streamUrl"`
	Commentators *string       `json:"commentators"`
}

// Match is a summary of a single match, as returned by list endpoints.
type Match struct {
	ID               wfa.Snowflake                      `json:"id"`
	Status           Status                             `json:"status"`
	ScheduledFor     *wfa.Time                          `json:"scheduledFor"`
	Times            Times                              `json:"times"`
	Hidden           bool                               `json:"hidden"`
	HomeTeam         wfa.TeamRef                        `json:"homeTeam"`
	AwayTeam         wfa.TeamRef                        `json:"awayTeam"`
	HomeScore        int                                `json:"homeScore"`
	AwayScore        int                                `json:"awayScore"`
	HomeScorePenalty int                                `json:"homeScorePenalty"`
	AwayScorePenalty int                                `json:"awayScorePenalty"`
	Competition      wfa.CompetitionRefWithOrganisation `json:"competition"`
	Season           wfa.SeasonRef                      `json:"season"`
	Court            *Court                             `json:"court"`
	MatchGroup       *GroupRef                          `json:"matchGroup"`
	Officials        Officials                          `json:"officials"`
	Streams          []Stream                           `json:"streams"`
}

// Player is a single lineup entry in a FullMatch.
type Player struct {
	Person   wfa.PersonRef  `json:"person"`
	Number   *int           `json:"number"`
	Position PlayerPosition `json:"position"`
	Captain  bool           `json:"captain"`
}

// Penalty is a single penalty shootout attempt in a FullMatch.
type Penalty struct {
	Sequence int           `json:"sequence"`
	TeamID   wfa.Snowflake `json:"teamId"`
	Scored   *bool         `json:"scored"`
	Player   wfa.PersonRef `json:"player"`
}

// FullMatch is the detailed representation of a match, including lineups,
// events and (if applicable) the penalty shootout.
type FullMatch struct {
	ID               wfa.Snowflake                      `json:"id"`
	Status           Status                             `json:"status"`
	ScheduledFor     *wfa.Time                          `json:"scheduledFor"`
	Times            Times                              `json:"times"`
	Hidden           bool                               `json:"hidden"`
	HomeTeam         wfa.TeamRef                        `json:"homeTeam"`
	AwayTeam         wfa.TeamRef                        `json:"awayTeam"`
	HomeScore        int                                `json:"homeScore"`
	AwayScore        int                                `json:"awayScore"`
	HomeScorePenalty int                                `json:"homeScorePenalty"`
	AwayScorePenalty int                                `json:"awayScorePenalty"`
	Competition      wfa.CompetitionRefWithOrganisation `json:"competition"`
	Season           wfa.SeasonRef                      `json:"season"`
	Court            *Court                             `json:"court"`
	MatchGroup       *GroupRef                          `json:"matchGroup"`
	Officials        Officials                          `json:"officials"`
	Streams          []Stream                           `json:"streams"`
	HomeLineups      []Player                           `json:"homeLineups"`
	AwayLineups      []Player                           `json:"awayLineups"`
	Events           Events                             `json:"events"`
	Penalties        []Penalty                          `json:"penalties"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams

	ID             []wfa.Snowflake
	TeamID         []wfa.Snowflake
	CompetitionID  []wfa.Snowflake
	OrganisationID []wfa.Snowflake
	SeasonID       []wfa.Snowflake
	MatchGroupID   []wfa.Snowflake
	CourtID        []wfa.Snowflake
	Status         []Status
	// ScheduledFrom is an ISO date or datetime string.
	ScheduledFrom string
	// ScheduledTo is an ISO date or datetime string.
	ScheduledTo     string
	OrderByDateDesc *bool
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetSnowflakes(v, "id", q.ID)
	wfa.SetSnowflakes(v, "teamId", q.TeamID)
	wfa.SetSnowflakes(v, "competitionId", q.CompetitionID)
	wfa.SetSnowflakes(v, "organisationId", q.OrganisationID)
	wfa.SetSnowflakes(v, "seasonId", q.SeasonID)
	wfa.SetSnowflakes(v, "matchGroupId", q.MatchGroupID)
	wfa.SetSnowflakes(v, "courtId", q.CourtID)
	wfa.SetEnums(v, "status", q.Status)
	wfa.SetString(v, "scheduledFrom", q.ScheduledFrom)
	wfa.SetString(v, "scheduledTo", q.ScheduledTo)
	wfa.SetBool(v, "orderByDateDesc", q.OrderByDateDesc)

	return v
}

// Service gives access to the /matches endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of matches matching the given query.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Match], error) {
	return wfa.Get[wfa.ListResponse[Match]](ctx, s.backend, "/matches", query.encode())
}

// Get retrieves detailed information about a specific match, including
// lineups, events and (if applicable) the penalty shootout.
func (s *Service) Get(ctx context.Context, id wfa.Snowflake) (FullMatch, error) {
	return wfa.Get[FullMatch](ctx, s.backend, fmt.Sprintf("/matches/%d", id), nil)
}
