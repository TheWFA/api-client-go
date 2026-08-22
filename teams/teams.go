// Package teams gives access to the WFA Matchday API's /teams endpoints.
package teams

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Team is a summary of a team.
type Team struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Abbreviated    string       `json:"abbreviated"`
	Nickname       string       `json:"nickname"`
	BadgeURL       string       `json:"badgeUrl"`
	GradientURL    *string      `json:"gradientUrl"`
	ThumbnailImage *string      `json:"thumbnailImage"`
	Primary        string       `json:"primary"`
	Secondary      string       `json:"secondary"`
	Club           *wfa.ClubRef `json:"club"`
}

// FullTeam is a Team with its history, when requested.
type FullTeam struct {
	ID             int                `json:"id"`
	Name           string             `json:"name"`
	Abbreviated    string             `json:"abbreviated"`
	Nickname       string             `json:"nickname"`
	BadgeURL       string             `json:"badgeUrl"`
	GradientURL    *string            `json:"gradientUrl"`
	ThumbnailImage *string            `json:"thumbnailImage"`
	Primary        string             `json:"primary"`
	Secondary      string             `json:"secondary"`
	Club           *wfa.ClubRef       `json:"club"`
	History        []wfa.HistoryEntry `json:"history,omitempty"`
}

// StaffRole is a role a person can hold on a team's staff.
type StaffRole string

const (
	StaffRoleHeadCoach      StaffRole = "head-coach"
	StaffRoleAssistantCoach StaffRole = "assistant-coach"
	StaffRoleMechanic       StaffRole = "mechanic"
	StaffRoleAssistant      StaffRole = "assistant"
)

// PlayerRegistration is a single player registration entry, returned by
// Service.Players.
type PlayerRegistration struct {
	Person             wfa.PersonRef      `json:"person"`
	Number             *int               `json:"number"`
	RegisteredAt       wfa.Time           `json:"registeredAt"`
	DeregisteredAt     *wfa.Time          `json:"deregisteredAt"`
	DeregisteredReason *string            `json:"deregisteredReason"`
	Competition        wfa.CompetitionRef `json:"competition"`
	Season             wfa.SeasonRef      `json:"season"`
}

// StaffRegistration is a single staff registration entry, returned by
// Service.Staff.
type StaffRegistration struct {
	Person             wfa.PersonRef      `json:"person"`
	Role               *StaffRole         `json:"role"`
	RegisteredAt       wfa.Time           `json:"registeredAt"`
	DeregisteredAt     *wfa.Time          `json:"deregisteredAt"`
	DeregisteredReason *string            `json:"deregisteredReason"`
	Competition        wfa.CompetitionRef `json:"competition"`
	Season             wfa.SeasonRef      `json:"season"`
}

// RegistrationCompetitionRef is a CompetitionRef with its competition type.
type RegistrationCompetitionRef struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	BadgeURL *string `json:"badgeUrl"`
	Type     string  `json:"type"`
}

// Registration is a single competition/season entry a team is registered
// into, returned by Service.Registrations.
type Registration struct {
	RegisteredAt wfa.Time                   `json:"registeredAt"`
	Competition  RegistrationCompetitionRef `json:"competition"`
	Season       wfa.SeasonRef              `json:"season"`
}

// StatsSummary is a team's aggregate results and discipline stats.
type StatsSummary struct {
	Played         int `json:"played"`
	Wins           int `json:"wins"`
	Draws          int `json:"draws"`
	Losses         int `json:"losses"`
	GoalsFor       int `json:"goalsFor"`
	GoalsAgainst   int `json:"goalsAgainst"`
	GoalDifference int `json:"goalDifference"`
	CleanSheets    int `json:"cleanSheets"`
	YellowCards    int `json:"yellowCards"`
	RedCards       int `json:"redCards"`
	PlayersUsed    int `json:"playersUsed"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	ClubID *int
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "clubId", q.ClubID)

	return v
}

// PlayersQuery holds the filters accepted by Service.Players.
type PlayersQuery struct {
	wfa.ListParams
	CompetitionID *int
	SeasonID      *int
	ActiveOnly    *bool
}

func (q PlayersQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)
	wfa.SetBool(v, "activeOnly", q.ActiveOnly)

	return v
}

// StaffQuery holds the filters accepted by Service.Staff.
type StaffQuery struct {
	wfa.ListParams
	CompetitionID *int
	SeasonID      *int
	ActiveOnly    *bool
	Role          []StaffRole
}

func (q StaffQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)
	wfa.SetBool(v, "activeOnly", q.ActiveOnly)
	wfa.SetEnums(v, "role", q.Role)

	return v
}

// RegistrationsQuery holds the filters accepted by Service.Registrations.
type RegistrationsQuery struct {
	wfa.ListParams
	CompetitionID *int
	SeasonID      *int
}

func (q RegistrationsQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)

	return v
}

// Service gives access to the /teams endpoints.
type Service struct {
	backend *wfa.Backend
	Stats   *StatsService
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend, Stats: &StatsService{backend: backend}}
}

// List retrieves a paginated list of teams.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Team], error) {
	return wfa.Get[wfa.ListResponse[Team]](ctx, s.backend, "/teams", query.encode())
}

// Get retrieves detailed information about a specific team.
func (s *Service) Get(ctx context.Context, id int) (FullTeam, error) {
	return wfa.Get[FullTeam](ctx, s.backend, fmt.Sprintf("/teams/%d", id), nil)
}

// Players retrieves the playing roster for a team.
func (s *Service) Players(ctx context.Context, id int, query PlayersQuery) (wfa.ListResponse[PlayerRegistration], error) {
	return wfa.Get[wfa.ListResponse[PlayerRegistration]](ctx, s.backend, fmt.Sprintf("/teams/%d/players", id), query.encode())
}

// Staff retrieves the staff roster for a team.
func (s *Service) Staff(ctx context.Context, id int, query StaffQuery) (wfa.ListResponse[StaffRegistration], error) {
	return wfa.Get[wfa.ListResponse[StaffRegistration]](ctx, s.backend, fmt.Sprintf("/teams/%d/staff", id), query.encode())
}

// Registrations retrieves the competitions and seasons a team is entered
// into.
func (s *Service) Registrations(ctx context.Context, id int, query RegistrationsQuery) (wfa.ListResponse[Registration], error) {
	return wfa.Get[wfa.ListResponse[Registration]](ctx, s.backend, fmt.Sprintf("/teams/%d/registrations", id), query.encode())
}

// Seasons retrieves every season a team has ever been registered into,
// newest first.
func (s *Service) Seasons(ctx context.Context, id int) (wfa.UnpaginatedListResponse[wfa.SeasonRef], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[wfa.SeasonRef]](ctx, s.backend, fmt.Sprintf("/teams/%d/seasons", id), nil)
}

// StatsService gives access to the /teams/{id}/stats endpoints, reached via
// Service.Stats.
type StatsService struct {
	backend *wfa.Backend
}

// Summary retrieves aggregate results and discipline stats for a team.
func (s *StatsService) Summary(ctx context.Context, id int, query wfa.StatsFilterQuery) (StatsSummary, error) {
	v := url.Values{}
	query.Apply(v)

	return wfa.Get[StatsSummary](ctx, s.backend, fmt.Sprintf("/teams/%d/stats/summary", id), v)
}

// Players retrieves per-player stats aggregates for a team.
func (s *StatsService) Players(ctx context.Context, id int, query wfa.PlayerStatsQuery) (wfa.ListResponse[wfa.PlayerStatsRow], error) {
	return wfa.Get[wfa.ListResponse[wfa.PlayerStatsRow]](ctx, s.backend, fmt.Sprintf("/teams/%d/stats/players", id), query.Encode())
}
