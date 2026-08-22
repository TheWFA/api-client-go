// Package persons gives access to the WFA Matchday API's /persons
// endpoints.
package persons

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/matches"
	"github.com/TheWFA/api-client-go/suspensions"
)

// Type categorizes the roles a person can hold.
type Type string

const (
	TypePlayer   Type = "player"
	TypeStaff    Type = "staff"
	TypeCoach    Type = "coach"
	TypeOfficial Type = "official"
)

// Person is a display-name-only reference to a person.
type Person = wfa.PersonRef

// FullPerson is a PersonRef with its role flags.
type FullPerson struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	FirstName  *string  `json:"firstName,omitempty"`
	LastName   *string  `json:"lastName,omitempty"`
	CreatedAt  wfa.Time `json:"createdAt"`
	IsPlayer   bool     `json:"isPlayer"`
	IsStaff    bool     `json:"isStaff"`
	IsCoach    bool     `json:"isCoach"`
	IsOfficial bool     `json:"isOfficial"`
}

// RegistrationType categorizes the kind of registration a person holds.
type RegistrationType string

const (
	RegistrationTypePlayer         RegistrationType = "player"
	RegistrationTypeHeadCoach      RegistrationType = "head-coach"
	RegistrationTypeAssistantCoach RegistrationType = "assistant-coach"
	RegistrationTypeAssistant      RegistrationType = "assistant"
	RegistrationTypeMechanic       RegistrationType = "mechanic"
)

// Registration is a single playing or staff registration, returned by
// Service.Registrations.
type Registration struct {
	Type               string             `json:"type"`
	Number             *int               `json:"number"`
	RegisteredAt       wfa.Time           `json:"registeredAt"`
	DeregisteredAt     *wfa.Time          `json:"deregisteredAt"`
	DeregisteredReason *string            `json:"deregisteredReason"`
	Team               wfa.TeamRef        `json:"team"`
	Competition        wfa.CompetitionRef `json:"competition"`
	Season             wfa.SeasonRef      `json:"season"`
}

// AppearanceMatchRef is a lightweight match reference, as embedded in
// Appearance.
type AppearanceMatchRef struct {
	ID           int         `json:"id"`
	Status       string      `json:"status"`
	ScheduledFor *wfa.Time   `json:"scheduledFor"`
	HomeTeam     wfa.TeamRef `json:"homeTeam"`
	AwayTeam     wfa.TeamRef `json:"awayTeam"`
}

// AppearanceMatchGroupRef is a lightweight match group reference, as
// embedded in Appearance.
type AppearanceMatchGroupRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Appearance is a single match appearance, returned by Service.Appearances.
type Appearance struct {
	Match         AppearanceMatchRef       `json:"match"`
	Team          wfa.TeamRef              `json:"team"`
	Competition   wfa.CompetitionRef       `json:"competition"`
	Season        wfa.SeasonRef            `json:"season"`
	MatchGroup    *AppearanceMatchGroupRef `json:"matchGroup"`
	SquadPosition string                   `json:"squadPosition"`
	Captain       bool                     `json:"captain"`
	Number        *int                     `json:"number"`
}

// StatsSummary is a summary of a person's career statistics.
type StatsSummary struct {
	Appearances     int `json:"appearances"`
	Starts          int `json:"starts"`
	Goals           int `json:"goals"`
	OwnGoals        int `json:"ownGoals"`
	Assists         int `json:"assists"`
	PenaltiesScored int `json:"penaltiesScored"`
	YellowCards     int `json:"yellowCards"`
	RedCards        int `json:"redCards"`
	Contributions   int `json:"contributions"`
}

// StatsMatchRef is a lightweight match reference, as embedded in
// GoalContribution and ReceivedCard.
type StatsMatchRef struct {
	ID           int         `json:"id"`
	ScheduledFor *wfa.Time   `json:"scheduledFor"`
	HomeTeam     wfa.TeamRef `json:"homeTeam"`
	AwayTeam     wfa.TeamRef `json:"awayTeam"`
}

// GoalContribution is a single goal or assist contribution, returned by both
// StatsService.Goals and StatsService.Assists.
type GoalContribution struct {
	ID          int                `json:"id"`
	MatchTime   *int               `json:"matchTime"`
	MatchPeriod *string            `json:"matchPeriod"`
	CreatedAt   wfa.Time           `json:"createdAt"`
	GoalType    string             `json:"goalType"`
	IsPenalty   bool               `json:"isPenalty"`
	Counterpart *wfa.PersonRef     `json:"counterpart"`
	Match       StatsMatchRef      `json:"match"`
	Team        wfa.TeamRef        `json:"team"`
	Competition wfa.CompetitionRef `json:"competition"`
	Season      wfa.SeasonRef      `json:"season"`
}

// Card is a card colour: yellow or red.
type Card string

const (
	CardYellow Card = "yellow"
	CardRed    Card = "red"
)

// CardOffenceRef is a lightweight reference to the offence a card was issued
// for.
type CardOffenceRef struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	Code             *string `json:"code"`
	SuspensionLength int     `json:"suspensionLength"`
}

// ReceivedCard is a single card received by a person, returned by
// StatsService.Cards.
type ReceivedCard struct {
	ID          int                `json:"id"`
	Card        Card               `json:"card"`
	MatchTime   *int               `json:"matchTime"`
	MatchPeriod *string            `json:"matchPeriod"`
	CreatedAt   wfa.Time           `json:"createdAt"`
	Offence     *CardOffenceRef    `json:"offence"`
	Match       StatsMatchRef      `json:"match"`
	Team        wfa.TeamRef        `json:"team"`
	Competition wfa.CompetitionRef `json:"competition"`
	Season      wfa.SeasonRef      `json:"season"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	Type          []Type
	TeamID        *int
	CompetitionID *int
	SeasonID      *int
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetEnums(v, "type", q.Type)
	wfa.SetInt(v, "teamId", q.TeamID)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)

	return v
}

// RegistrationsQuery holds the filters accepted by Service.Registrations.
type RegistrationsQuery struct {
	wfa.ListParams
	Type          []RegistrationType
	TeamID        *int
	CompetitionID *int
	SeasonID      *int
	ActiveOnly    *bool
}

func (q RegistrationsQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetEnums(v, "type", q.Type)
	wfa.SetInt(v, "teamId", q.TeamID)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)
	wfa.SetBool(v, "activeOnly", q.ActiveOnly)

	return v
}

// AppearancesQuery holds the filters accepted by Service.Appearances.
type AppearancesQuery struct {
	wfa.ListParams
	TeamID        *int
	CompetitionID *int
	SeasonID      *int
	MatchGroupID  *int
	Position      []matches.PlayerPosition
	// ScheduledFrom is an ISO date or datetime string.
	ScheduledFrom string
	// ScheduledTo is an ISO date or datetime string.
	ScheduledTo string
}

func (q AppearancesQuery) encode() url.Values {
	v := url.Values{}
	q.Apply(v)
	wfa.SetInt(v, "teamId", q.TeamID)
	wfa.SetInt(v, "competitionId", q.CompetitionID)
	wfa.SetInt(v, "seasonId", q.SeasonID)
	wfa.SetInt(v, "matchGroupId", q.MatchGroupID)
	wfa.SetEnums(v, "position", q.Position)
	wfa.SetString(v, "scheduledFrom", q.ScheduledFrom)
	wfa.SetString(v, "scheduledTo", q.ScheduledTo)

	return v
}

// StatsQuery holds the filters accepted by StatsService.Goals and
// StatsService.Assists.
type StatsQuery struct {
	wfa.ListParams
	wfa.StatsFilterQuery
}

func (q StatsQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	q.StatsFilterQuery.Apply(v)

	return v
}

// CardsQuery holds the filters accepted by StatsService.Cards.
type CardsQuery struct {
	wfa.ListParams
	wfa.StatsFilterQuery
	Card []Card
}

func (q CardsQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	q.StatsFilterQuery.Apply(v)
	wfa.SetEnums(v, "card", q.Card)

	return v
}

// Service gives access to the /persons endpoints.
type Service struct {
	backend *wfa.Backend
	Stats   *StatsService
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend, Stats: &StatsService{backend: backend}}
}

// List retrieves a paginated list of people with display names only.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Person], error) {
	return wfa.Get[wfa.ListResponse[Person]](ctx, s.backend, "/persons", query.encode())
}

// Get retrieves detailed information about a specific person.
func (s *Service) Get(ctx context.Context, id int) (FullPerson, error) {
	return wfa.Get[FullPerson](ctx, s.backend, fmt.Sprintf("/persons/%d", id), nil)
}

// Registrations retrieves a person's playing and staff registrations,
// newest first.
func (s *Service) Registrations(ctx context.Context, id int, query RegistrationsQuery) (wfa.ListResponse[Registration], error) {
	return wfa.Get[wfa.ListResponse[Registration]](ctx, s.backend, fmt.Sprintf("/persons/%d/registrations", id), query.encode())
}

// Appearances retrieves a person's match appearances.
func (s *Service) Appearances(ctx context.Context, id int, query AppearancesQuery) (wfa.ListResponse[Appearance], error) {
	return wfa.Get[wfa.ListResponse[Appearance]](ctx, s.backend, fmt.Sprintf("/persons/%d/appearances", id), query.encode())
}

// Suspensions retrieves a person's suspensions.
func (s *Service) Suspensions(ctx context.Context, id int, query suspensions.PersonQuery) (wfa.ListResponse[suspensions.Suspension], error) {
	return wfa.Get[wfa.ListResponse[suspensions.Suspension]](ctx, s.backend, fmt.Sprintf("/persons/%d/suspensions", id), query.Encode())
}

// StatsService gives access to the /persons/{id}/stats endpoints, reached
// via Service.Stats.
type StatsService struct {
	backend *wfa.Backend
}

// Summary retrieves a summary of a person's career statistics.
func (s *StatsService) Summary(ctx context.Context, id int, query wfa.StatsFilterQuery) (StatsSummary, error) {
	v := url.Values{}
	query.Apply(v)

	return wfa.Get[StatsSummary](ctx, s.backend, fmt.Sprintf("/persons/%d/stats/summary", id), v)
}

// Goals retrieves the goals scored by a person.
func (s *StatsService) Goals(ctx context.Context, id int, query StatsQuery) (wfa.ListResponse[GoalContribution], error) {
	return wfa.Get[wfa.ListResponse[GoalContribution]](ctx, s.backend, fmt.Sprintf("/persons/%d/stats/goals", id), query.encode())
}

// Assists retrieves the assists made by a person.
func (s *StatsService) Assists(ctx context.Context, id int, query StatsQuery) (wfa.ListResponse[GoalContribution], error) {
	return wfa.Get[wfa.ListResponse[GoalContribution]](ctx, s.backend, fmt.Sprintf("/persons/%d/stats/assists", id), query.encode())
}

// Cards retrieves the cards received by a person.
func (s *StatsService) Cards(ctx context.Context, id int, query CardsQuery) (wfa.ListResponse[ReceivedCard], error) {
	return wfa.Get[wfa.ListResponse[ReceivedCard]](ctx, s.backend, fmt.Sprintf("/persons/%d/stats/cards", id), query.encode())
}
