// Package competitions gives access to the WFA Matchday API's /competitions
// endpoints.
package competitions

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/teams"
)

// Type categorizes a competition's format.
type Type string

const (
	TypeLeague   Type = "league"
	TypeCup      Type = "cup"
	TypeFriendly Type = "friendly"
)

// Competition is a summary of a competition.
type Competition struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	Type         Type                 `json:"type"`
	BadgeURL     *string              `json:"badgeUrl"`
	SortOrder    int                  `json:"sortOrder"`
	Hidden       bool                 `json:"hidden"`
	Organisation *wfa.OrganisationRef `json:"organisation"`
}

// MatchGroupType categorizes a competition match group's format.
type MatchGroupType string

const (
	MatchGroupTypeGameWeek  MatchGroupType = "game-week"
	MatchGroupTypePool      MatchGroupType = "pool"
	MatchGroupTypeKnockout  MatchGroupType = "knockout"
	MatchGroupTypeTwoLegged MatchGroupType = "two-legged"
)

// MatchGroupProgression describes where teams progress to from a match
// group.
type MatchGroupProgression struct {
	ToGroupID            int  `json:"toGroupId"`
	ProgressingTeamCount *int `json:"progressingTeamCount"`
}

// MatchGroup is a single stage of a competition's structure, e.g. a game
// week, pool or knockout round.
type MatchGroup struct {
	ID             int                     `json:"id"`
	GroupName      string                  `json:"groupName"`
	GroupType      *MatchGroupType         `json:"groupType"`
	RoundNumber    *int                    `json:"roundNumber"`
	AdvancingSpots *int                    `json:"advancingSpots"`
	SeasonID       int                     `json:"seasonId"`
	Progressions   []MatchGroupProgression `json:"progressions"`
}

// FullCompetition is a Competition with its current season, match groups,
// and, when requested, its history.
type FullCompetition struct {
	ID           int                   `json:"id"`
	Name         string                `json:"name"`
	Type         Type                  `json:"type"`
	BadgeURL     *string               `json:"badgeUrl"`
	SortOrder    int                   `json:"sortOrder"`
	Hidden       bool                  `json:"hidden"`
	Organisation *wfa.OrganisationRef  `json:"organisation"`
	Season       *wfa.SeasonWithActive `json:"season"`
	MatchGroups  []MatchGroup          `json:"matchGroups"`
	History      []wfa.HistoryEntry    `json:"history,omitempty"`
}

// MatchGroupTeam is a team seeded into a competition match group.
type MatchGroupTeam struct {
	MatchGroupID int         `json:"matchGroupId"`
	Team         wfa.TeamRef `json:"team"`
	SeedNumber   *int        `json:"seedNumber"`
}

// StatsSummary is a competition's aggregate stats.
type StatsSummary struct {
	Matches       int     `json:"matches"`
	Teams         int     `json:"teams"`
	Goals         int     `json:"goals"`
	OwnGoals      int     `json:"ownGoals"`
	GoalsPerMatch float64 `json:"goalsPerMatch"`
	YellowCards   int     `json:"yellowCards"`
	RedCards      int     `json:"redCards"`
	CleanSheets   int     `json:"cleanSheets"`
}

// TeamsStatsOrderBy is a sort key accepted by StatsService.Teams.
type TeamsStatsOrderBy string

const (
	TeamsStatsOrderByName           TeamsStatsOrderBy = "name"
	TeamsStatsOrderByPlayed         TeamsStatsOrderBy = "played"
	TeamsStatsOrderByWins           TeamsStatsOrderBy = "wins"
	TeamsStatsOrderByGoalsFor       TeamsStatsOrderBy = "goalsFor"
	TeamsStatsOrderByGoalsAgainst   TeamsStatsOrderBy = "goalsAgainst"
	TeamsStatsOrderByGoalDifference TeamsStatsOrderBy = "goalDifference"
	TeamsStatsOrderByCleanSheets    TeamsStatsOrderBy = "cleanSheets"
	TeamsStatsOrderByYellowCards    TeamsStatsOrderBy = "yellowCards"
	TeamsStatsOrderByRedCards       TeamsStatsOrderBy = "redCards"
	TeamsStatsOrderByPoints         TeamsStatsOrderBy = "points"
)

// TeamsStatsQuery holds the filters accepted by StatsService.Teams.
type TeamsStatsQuery struct {
	wfa.ListParams
	wfa.StatsFilterQuery
	OrderBy TeamsStatsOrderBy
}

func (q TeamsStatsQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	q.StatsFilterQuery.Apply(v)
	wfa.SetString(v, "orderBy", string(q.OrderBy))

	return v
}

// TeamStatsRow is a per-team stats aggregate row for a competition.
type TeamStatsRow struct {
	Team           wfa.TeamRef `json:"team"`
	Played         int         `json:"played"`
	Wins           int         `json:"wins"`
	Draws          int         `json:"draws"`
	Losses         int         `json:"losses"`
	GoalsFor       int         `json:"goalsFor"`
	GoalsAgainst   int         `json:"goalsAgainst"`
	GoalDifference int         `json:"goalDifference"`
	GoalsPerMatch  float64     `json:"goalsPerMatch"`
	CleanSheets    int         `json:"cleanSheets"`
	YellowCards    int         `json:"yellowCards"`
	RedCards       int         `json:"redCards"`
	Points         *int        `json:"points"`
}

// TableRow is a single row of a competition's league table.
type TableRow struct {
	Team           wfa.TeamRef `json:"team"`
	Played         int         `json:"played"`
	Wins           int         `json:"wins"`
	Draws          int         `json:"draws"`
	Losses         int         `json:"losses"`
	GoalsFor       int         `json:"goalsFor"`
	GoalsAgainst   int         `json:"goalsAgainst"`
	GoalDifference int         `json:"goalDifference"`
	GoalsPerMatch  float64     `json:"goalsPerMatch"`
	CleanSheets    int         `json:"cleanSheets"`
	YellowCards    int         `json:"yellowCards"`
	RedCards       int         `json:"redCards"`
	Position       int         `json:"position"`
	Points         int         `json:"points"`
}

// Table is a competition's league table for a season.
type Table struct {
	Season     wfa.SeasonFull `json:"season"`
	Items      []TableRow     `json:"items"`
	TotalItems int            `json:"totalItems"`
}

// Team is a team registered for a competition, shaped identically to
// teams.Team.
type Team = teams.Team

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	OrganisationID *int
	Type           Type
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	wfa.SetInt(v, "organisationId", q.OrganisationID)
	wfa.SetString(v, "type", string(q.Type))

	return v
}

// Service gives access to the /competitions endpoints.
type Service struct {
	backend *wfa.Backend
	Stats   *StatsService
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend, Stats: &StatsService{backend: backend}}
}

// List retrieves a paginated list of competitions.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Competition], error) {
	return wfa.Get[wfa.ListResponse[Competition]](ctx, s.backend, "/competitions", query.encode())
}

// Get retrieves detailed information about a specific competition. If
// seasonID is nil, the competition's active or most recent season is used.
func (s *Service) Get(ctx context.Context, id int, seasonID *int) (FullCompetition, error) {
	v := url.Values{}
	wfa.SetInt(v, "seasonId", seasonID)

	return wfa.Get[FullCompetition](ctx, s.backend, fmt.Sprintf("/competitions/%d", id), v)
}

// Teams retrieves the teams registered for a competition and season. If
// seasonID is nil, the competition's active or most recent season is used.
func (s *Service) Teams(ctx context.Context, id int, seasonID *int) (wfa.UnpaginatedListResponse[Team], error) {
	v := url.Values{}
	wfa.SetInt(v, "seasonId", seasonID)

	return wfa.Get[wfa.UnpaginatedListResponse[Team]](ctx, s.backend, fmt.Sprintf("/competitions/%d/teams", id), v)
}

// Seasons retrieves every season a competition has run.
func (s *Service) Seasons(ctx context.Context, id int) (wfa.UnpaginatedListResponse[wfa.SeasonFull], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[wfa.SeasonFull]](ctx, s.backend, fmt.Sprintf("/competitions/%d/seasons", id), nil)
}

// Table retrieves the league table for a competition and season. If
// seasonID is nil, the competition's active or most recent season is used.
//
// Returns an *wfa.APIError with StatusCode 400 for cup and friendly
// competitions, and for a competition with no registered season.
func (s *Service) Table(ctx context.Context, id int, seasonID *int) (Table, error) {
	v := url.Values{}
	wfa.SetInt(v, "seasonId", seasonID)

	return wfa.Get[Table](ctx, s.backend, fmt.Sprintf("/competitions/%d/table", id), v)
}

// MatchGroupTeams retrieves the teams in a match group stage, with their
// bracket seeds.
func (s *Service) MatchGroupTeams(ctx context.Context, id, groupID int) (wfa.UnpaginatedListResponse[MatchGroupTeam], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[MatchGroupTeam]](ctx, s.backend, fmt.Sprintf("/competitions/%d/match-groups/%d/teams", id, groupID), nil)
}

// StatsService gives access to the /competitions/{id}/stats endpoints,
// reached via Service.Stats.
type StatsService struct {
	backend *wfa.Backend
}

// Summary retrieves aggregate statistics for a competition.
func (s *StatsService) Summary(ctx context.Context, id int, query wfa.StatsFilterQuery) (StatsSummary, error) {
	v := url.Values{}
	query.Apply(v)

	return wfa.Get[StatsSummary](ctx, s.backend, fmt.Sprintf("/competitions/%d/stats/summary", id), v)
}

// Teams retrieves per-team aggregate statistics for a competition.
func (s *StatsService) Teams(ctx context.Context, id int, query TeamsStatsQuery) (wfa.ListResponse[TeamStatsRow], error) {
	return wfa.Get[wfa.ListResponse[TeamStatsRow]](ctx, s.backend, fmt.Sprintf("/competitions/%d/stats/teams", id), query.encode())
}

// Players retrieves per-player aggregate statistics for a competition.
func (s *StatsService) Players(ctx context.Context, id int, query wfa.PlayerStatsQuery) (wfa.ListResponse[wfa.PlayerStatsRow], error) {
	return wfa.Get[wfa.ListResponse[wfa.PlayerStatsRow]](ctx, s.backend, fmt.Sprintf("/competitions/%d/stats/players", id), query.Encode())
}
